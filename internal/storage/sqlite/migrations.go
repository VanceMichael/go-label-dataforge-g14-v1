package sqlite

import (
	"context"
	"database/sql"
)

func Migrate(ctx context.Context, db *sql.DB) error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS schema_migrations(version INTEGER PRIMARY KEY, applied_at TEXT NOT NULL);`,
		`CREATE TABLE IF NOT EXISTS tenants(id TEXT PRIMARY KEY,name TEXT NOT NULL UNIQUE,created_at TEXT NOT NULL);`,
		`CREATE TABLE IF NOT EXISTS users(id TEXT PRIMARY KEY,tenant_id TEXT NOT NULL REFERENCES tenants(id),email TEXT NOT NULL,role TEXT NOT NULL,active INTEGER NOT NULL DEFAULT 1,created_at TEXT NOT NULL,UNIQUE(tenant_id,email));`,
		`CREATE TABLE IF NOT EXISTS sessions(id TEXT PRIMARY KEY,user_id TEXT NOT NULL REFERENCES users(id),token_hash TEXT NOT NULL UNIQUE,expires_at TEXT NOT NULL,created_at TEXT NOT NULL,revoked_at TEXT);`,
		`CREATE TABLE IF NOT EXISTS data_resources(id TEXT PRIMARY KEY,tenant_id TEXT NOT NULL REFERENCES tenants(id),owner_id TEXT NOT NULL REFERENCES users(id),code TEXT NOT NULL,name TEXT NOT NULL,description TEXT NOT NULL,status TEXT NOT NULL,version INTEGER NOT NULL DEFAULT 0,created_at TEXT NOT NULL,updated_at TEXT NOT NULL,UNIQUE(tenant_id,code));`,
		`CREATE TABLE IF NOT EXISTS resource_versions(id TEXT PRIMARY KEY,resource_id TEXT NOT NULL REFERENCES data_resources(id),number INTEGER NOT NULL,schema TEXT NOT NULL,content_hash TEXT NOT NULL,created_at TEXT NOT NULL,UNIQUE(resource_id,number));`,
		`CREATE TABLE IF NOT EXISTS registration_reviews(id TEXT PRIMARY KEY,resource_id TEXT NOT NULL REFERENCES data_resources(id),reviewer_id TEXT REFERENCES users(id),decision TEXT NOT NULL,notes TEXT NOT NULL,status TEXT NOT NULL,lease_until TEXT,created_at TEXT NOT NULL);`,
		`CREATE TABLE IF NOT EXISTS authorization_requests(id TEXT PRIMARY KEY,tenant_id TEXT NOT NULL REFERENCES tenants(id),resource_id TEXT NOT NULL REFERENCES data_resources(id),applicant_id TEXT NOT NULL REFERENCES users(id),status TEXT NOT NULL,purpose TEXT NOT NULL,quota INTEGER NOT NULL,used INTEGER NOT NULL DEFAULT 0,version INTEGER NOT NULL DEFAULT 0,created_at TEXT NOT NULL,updated_at TEXT NOT NULL);`,
		`CREATE TABLE IF NOT EXISTS sandbox_leases(id TEXT PRIMARY KEY,authorization_id TEXT NOT NULL REFERENCES authorization_requests(id),holder_id TEXT NOT NULL REFERENCES users(id),status TEXT NOT NULL,expires_at TEXT NOT NULL,version INTEGER NOT NULL DEFAULT 0,created_at TEXT NOT NULL);`,
		`CREATE TABLE IF NOT EXISTS sandbox_runs(id TEXT PRIMARY KEY,lease_id TEXT NOT NULL REFERENCES sandbox_leases(id),status TEXT NOT NULL,input TEXT NOT NULL,output TEXT NOT NULL,started_at TEXT,finished_at TEXT);`,
		`CREATE TABLE IF NOT EXISTS data_products(id TEXT PRIMARY KEY,tenant_id TEXT NOT NULL REFERENCES tenants(id),owner_id TEXT NOT NULL REFERENCES users(id),name TEXT NOT NULL,resource_id TEXT NOT NULL REFERENCES data_resources(id),status TEXT NOT NULL,version INTEGER NOT NULL DEFAULT 0,created_at TEXT NOT NULL,updated_at TEXT NOT NULL);`,
		`CREATE TABLE IF NOT EXISTS product_releases(id TEXT PRIMARY KEY,product_id TEXT NOT NULL REFERENCES data_products(id),version INTEGER NOT NULL,status TEXT NOT NULL,published_at TEXT,created_at TEXT NOT NULL,UNIQUE(product_id,version));`,
		`CREATE TABLE IF NOT EXISTS audit_events(id TEXT PRIMARY KEY,tenant_id TEXT NOT NULL REFERENCES tenants(id),actor_id TEXT,object_type TEXT NOT NULL,object_id TEXT NOT NULL,action TEXT NOT NULL,result TEXT NOT NULL,request_id TEXT NOT NULL,detail TEXT NOT NULL,created_at TEXT NOT NULL);`,
		`CREATE TABLE IF NOT EXISTS outbox_jobs(id TEXT PRIMARY KEY,topic TEXT NOT NULL,payload TEXT NOT NULL,attempts INTEGER NOT NULL DEFAULT 0,available_at TEXT NOT NULL,locked_until TEXT,status TEXT NOT NULL,last_error TEXT NOT NULL DEFAULT '',created_at TEXT NOT NULL);`,
		`CREATE INDEX IF NOT EXISTS idx_resources_status ON data_resources(status,updated_at);`,
		`CREATE INDEX IF NOT EXISTS idx_leases_expiry ON sandbox_leases(status,expires_at);`,
		`CREATE INDEX IF NOT EXISTS idx_outbox_ready ON outbox_jobs(status,available_at);`,
	}
	for _, s := range stmts {
		if _, e := db.ExecContext(ctx, s); e != nil {
			return e
		}
	}
	var n int
	if e := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations").Scan(&n); e != nil {
		return e
	}
	if n == 0 {
		if _, e := db.ExecContext(ctx, "INSERT INTO schema_migrations(version,applied_at) VALUES(1,datetime('now'))"); e != nil {
			return e
		}
	}
	return nil
}
