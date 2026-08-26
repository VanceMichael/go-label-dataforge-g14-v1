package sqlite

import (
	"context"
	"database/sql"
	"github.com/VanceMichael/go-label-dataforge-g14-v1/internal/domain"
	"time"
)

func SeedDemo(ctx context.Context, r *Repo) error {
	now := time.Now().UTC()
	tenant := domain.Tenant{ID: "tenant-demo", Name: "甘肃数据工场", CreatedAt: now}
	user := domain.User{ID: "user-provider", TenantID: tenant.ID, Email: "provider@example.com", Role: domain.RoleProvider, Active: true, CreatedAt: now}
	reviewer := domain.User{ID: "user-reviewer", TenantID: tenant.ID, Email: "reviewer@example.com", Role: domain.RoleReviewer, Active: true, CreatedAt: now}
	return r.WithTx(ctx, func(tx *sql.Tx) error {
		if e := r.CreateTenant(ctx, tx, tenant); e != nil {
			return e
		}
		if e := r.CreateUser(ctx, tx, user); e != nil {
			return e
		}
		return r.CreateUser(ctx, tx, reviewer)
	})
}
func ClearDemo(ctx context.Context, d *DB) error {
	tables := []string{"outbox_jobs", "audit_events", "product_releases", "data_products", "sandbox_runs", "sandbox_leases", "authorization_requests", "registration_reviews", "resource_versions", "data_resources", "sessions", "users", "tenants"}
	return d.WithTx(ctx, func(tx *sql.Tx) error {
		for _, t := range tables {
			if _, e := tx.ExecContext(ctx, "DELETE FROM "+t); e != nil {
				return e
			}
		}
		return nil
	})
}
