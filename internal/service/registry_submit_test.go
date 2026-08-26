package service_test

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/VanceMichael/go-label-dataforge-g14-v1/internal/apperrors"
	"github.com/VanceMichael/go-label-dataforge-g14-v1/internal/clock"
	"github.com/VanceMichael/go-label-dataforge-g14-v1/internal/domain"
	"github.com/VanceMichael/go-label-dataforge-g14-v1/internal/service"
	"github.com/VanceMichael/go-label-dataforge-g14-v1/internal/storage/sqlite"
)

// failingAuditStore wraps the real repo and makes the audit write fail a
// configurable number of times, simulating the audit store erroring during the
// write phase of a submit.
type failingAuditStore struct {
	*sqlite.Repo
	mu        sync.Mutex
	remaining int
}

func (f *failingAuditStore) AddAuditTx(ctx context.Context, tx *sql.Tx, a domain.AuditEvent) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.remaining > 0 {
		f.remaining--
		return apperrors.ErrUnavailable
	}
	return f.Repo.AddAuditTx(ctx, tx, a)
}

func newTestStore(t *testing.T, auditFails int) *failingAuditStore {
	t.Helper()
	dir := t.TempDir()
	db, err := sqlite.Open(context.Background(), filepath.Join(dir, "test.db"))
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	if err := sqlite.SeedDemo(context.Background(), sqlite.NewRepo(db)); err != nil {
		t.Fatalf("seed: %v", err)
	}
	return &failingAuditStore{Repo: sqlite.NewRepo(db), remaining: auditFails}
}

func TestSubmitAuditWriteFailureKeepsDraft(t *testing.T) {
	store := newTestStore(t, 1)
	now := time.Now().UTC()
	reg := service.Registry{Store: store, Clock: clock.Fixed{T: now}}

	res := domain.Resource{
		ID: "res-1", TenantID: "tenant-demo", OwnerID: "user-provider",
		Code: "demo-code", Name: "demo", Description: "d",
		Status: domain.ResourceDraft, Version: 1,
		CreatedAt: now, UpdatedAt: now,
	}
	ver := domain.ResourceVersion{ID: "res-1-v1", ResourceID: "res-1", Number: 1, Schema: "json", ContentHash: "h", CreatedAt: now}
	rev := domain.Review{ID: "res-1-review", ResourceID: "res-1", Status: "pending", Notes: "", CreatedAt: now}
	if _, err := reg.Create(context.Background(), res, ver, rev, "req-1"); err != nil {
		t.Fatalf("create: %v", err)
	}

	// First submit: the audit write fails. The submit must fail and the
	// resource must remain draft with its version unchanged.
	err := reg.Submit(context.Background(), "res-1", 1, "user-provider", "tenant-demo", "req-2")
	if !errors.Is(err, apperrors.ErrUnavailable) {
		t.Fatalf("expected unavailable error, got %v", err)
	}
	got, err := store.GetResource(context.Background(), "res-1")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Status != domain.ResourceDraft {
		t.Fatalf("expected draft after audit failure, got %s", got.Status)
	}
	if got.Version != 1 {
		t.Fatalf("expected version unchanged after audit failure, got %d", got.Version)
	}

	// Dependency recovered: retrying the submit with the same version must
	// now succeed and move the resource to submitted.
	if err := reg.Submit(context.Background(), "res-1", 1, "user-provider", "tenant-demo", "req-3"); err != nil {
		t.Fatalf("retry submit: %v", err)
	}
	got, err = store.GetResource(context.Background(), "res-1")
	if err != nil {
		t.Fatalf("get after retry: %v", err)
	}
	if got.Status != domain.ResourceSubmitted {
		t.Fatalf("expected submitted after retry, got %s", got.Status)
	}
	if got.Version != 2 {
		t.Fatalf("expected version 2 after retry, got %d", got.Version)
	}
}
