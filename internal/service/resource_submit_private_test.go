package service_test

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"
	"time"

	"github.com/VanceMichael/go-label-dataforge-g14-v1/internal/clock"
	"github.com/VanceMichael/go-label-dataforge-g14-v1/internal/domain"
	"github.com/VanceMichael/go-label-dataforge-g14-v1/internal/service"
	"github.com/VanceMichael/go-label-dataforge-g14-v1/internal/storage/sqlite"
)

func TestResourceSubmitAuditFailureRollsBack(t *testing.T) {
	ctx := context.Background()
	db, err := sqlite.Open(ctx, filepath.Join(t.TempDir(), "resource-submit.db"))
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	defer db.Close()
	repo := sqlite.NewRepo(db)
	now := time.Date(2026, 8, 26, 0, 0, 0, 0, time.UTC)
	tenant := domain.Tenant{ID: "tenant-1", Name: "DataForge", CreatedAt: now}
	owner := domain.User{ID: "user-1", TenantID: tenant.ID, Email: "owner@example.com", Role: domain.RoleProvider, Active: true, CreatedAt: now}
	if err := repo.WithTx(ctx, func(tx *sql.Tx) error {
		if err := repo.CreateTenant(ctx, tx, tenant); err != nil {
			return err
		}
		return repo.CreateUser(ctx, tx, owner)
	}); err != nil {
		t.Fatalf("seed owner: %v", err)
	}
	registry := service.Registry{Store: repo, Clock: clock.Fixed{T: now}}
	resource := domain.Resource{ID: "resource-1", TenantID: tenant.ID, OwnerID: owner.ID, Code: "DF-1", Name: "原始资源", Description: "待提交", Status: domain.ResourceDraft, Version: 1, CreatedAt: now, UpdatedAt: now}
	if _, err := registry.Create(ctx, resource, domain.ResourceVersion{ID: "resource-1-v1", ResourceID: resource.ID, Number: 1, Schema: "json", ContentHash: "hash", CreatedAt: now}, domain.Review{ID: "review-1", ResourceID: resource.ID, Status: "pending", CreatedAt: now}, "create-1"); err != nil {
		t.Fatalf("create resource: %v", err)
	}
	if err := repo.WithTx(ctx, func(tx *sql.Tx) error {
		return repo.AddAuditTx(ctx, tx, domain.AuditEvent{ID: "audit-submit", TenantID: tenant.ID, ActorID: owner.ID, ObjectType: "resource", ObjectID: resource.ID, Action: "collision", Result: "ok", RequestID: "seed", Detail: "reserved event id", CreatedAt: now})
	}); err != nil {
		t.Fatalf("seed audit collision: %v", err)
	}
	if err := registry.Submit(ctx, resource.ID, 1, owner.ID, tenant.ID, "submit-1"); err == nil {
		t.Fatal("expected audit persistence failure")
	}
	got, err := repo.GetResource(ctx, resource.ID)
	if err != nil {
		t.Fatalf("read resource after failed submit: %v", err)
	}
	if got.Status != domain.ResourceDraft || got.Version != 1 {
		t.Fatalf("failed submit changed resource: status=%s version=%d", got.Status, got.Version)
	}
	var auditCount int
	if err := db.SQL.QueryRowContext(ctx, "SELECT COUNT(*) FROM audit_events WHERE object_id=?", resource.ID).Scan(&auditCount); err != nil {
		t.Fatalf("count audit events: %v", err)
	}
	if auditCount != 1 {
		t.Fatalf("failed submit changed audit history: count=%d", auditCount)
	}
	if _, err := db.SQL.ExecContext(ctx, "DELETE FROM audit_events WHERE id=?", "audit-submit"); err != nil {
		t.Fatalf("clear audit collision: %v", err)
	}
	if err := registry.Submit(ctx, resource.ID, 1, owner.ID, tenant.ID, "submit-2"); err != nil {
		t.Fatalf("retry submit after audit recovery: %v", err)
	}
	got, err = repo.GetResource(ctx, resource.ID)
	if err != nil {
		t.Fatalf("read resource after retry: %v", err)
	}
	if got.Status != domain.ResourceSubmitted || got.Version != 2 {
		t.Fatalf("successful retry did not submit resource: status=%s version=%d", got.Status, got.Version)
	}
	if err := db.SQL.QueryRowContext(ctx, "SELECT COUNT(*) FROM audit_events WHERE object_id=?", resource.ID).Scan(&auditCount); err != nil {
		t.Fatalf("count audit events after retry: %v", err)
	}
	if auditCount != 1 {
		t.Fatalf("successful retry wrote wrong audit history: count=%d", auditCount)
	}
}
