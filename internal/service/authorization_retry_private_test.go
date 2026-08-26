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

func TestAuthorizationRetryReusesCommittedRequest(t *testing.T) {
	ctx := context.Background()
	now := time.Date(2026, 8, 26, 0, 0, 0, 0, time.UTC)
	db, err := sqlite.Open(ctx, filepath.Join(t.TempDir(), "authorization-retry.db"))
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	defer db.Close()
	repo := sqlite.NewRepo(db)
	tenant := domain.Tenant{ID: "tenant-auth", Name: "授权租户", CreatedAt: now}
	applicant := domain.User{ID: "user-auth", TenantID: tenant.ID, Email: "applicant@example.com", Role: domain.RoleProvider, Active: true, CreatedAt: now}
	resource := domain.Resource{ID: "resource-auth", TenantID: tenant.ID, OwnerID: applicant.ID, Code: "AUTH-1", Name: "授权资源", Description: "可授权资源", Status: domain.ResourceRegistered, Version: 2, CreatedAt: now, UpdatedAt: now}
	if err := repo.WithTx(ctx, func(tx *sql.Tx) error {
		if err := repo.CreateTenant(ctx, tx, tenant); err != nil {
			return err
		}
		if err := repo.CreateUser(ctx, tx, applicant); err != nil {
			return err
		}
		_, err := tx.ExecContext(ctx, "INSERT INTO data_resources(id,tenant_id,owner_id,code,name,description,status,version,created_at,updated_at) VALUES(?,?,?,?,?,?,?,?,?,?)", resource.ID, resource.TenantID, resource.OwnerID, resource.Code, resource.Name, resource.Description, resource.Status, resource.Version, now.Format(time.RFC3339Nano), now.Format(time.RFC3339Nano))
		return err
	}); err != nil {
		t.Fatalf("seed authorization context: %v", err)
	}
	a := domain.Authorization{ID: "authorization-retry-1", TenantID: tenant.ID, ResourceID: resource.ID, ApplicantID: applicant.ID, Purpose: "sandbox", Quota: 10, Version: 1, CreatedAt: now, UpdatedAt: now}
	authorizer := service.Authorizer{Store: repo, Clock: clock.Fixed{T: now}}
	if err := authorizer.Request(ctx, a); err != nil {
		t.Fatalf("first authorization request: %v", err)
	}
	if err := authorizer.Request(ctx, a); err != nil {
		t.Fatalf("retry should reuse committed authorization: %v", err)
	}
	got, err := repo.GetAuthorization(ctx, a.ID)
	if err != nil {
		t.Fatalf("read committed authorization: %v", err)
	}
	if got.Status != domain.AuthRequested || got.Version != a.Version || got.Purpose != a.Purpose {
		t.Fatalf("retry changed committed authorization: status=%s version=%d purpose=%s", got.Status, got.Version, got.Purpose)
	}
	var count int
	if err := db.SQL.QueryRowContext(ctx, "SELECT COUNT(*) FROM authorization_requests WHERE id=?", a.ID).Scan(&count); err != nil {
		t.Fatalf("count authorization rows: %v", err)
	}
	if count != 1 {
		t.Fatalf("retry created duplicate authorization rows: %d", count)
	}
}
