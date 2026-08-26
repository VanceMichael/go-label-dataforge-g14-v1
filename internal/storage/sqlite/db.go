package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	_ "modernc.org/sqlite"
	"os"
	"path/filepath"
)

type DB struct{ SQL *sql.DB }

func Open(ctx context.Context, path string) (*DB, error) {
	if dir := filepath.Dir(path); dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, err
		}
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(8)
	if _, err = db.ExecContext(ctx, "PRAGMA foreign_keys=ON; PRAGMA busy_timeout=5000;"); err != nil {
		db.Close()
		return nil, err
	}
	d := &DB{SQL: db}
	if err := Migrate(ctx, db); err != nil {
		db.Close()
		return nil, err
	}
	return d, nil
}
func (d *DB) Close() error                   { return d.SQL.Close() }
func (d *DB) Ping(ctx context.Context) error { return d.SQL.PingContext(ctx) }
func tx(ctx context.Context, db *sql.DB, fn func(*sql.Tx) error) error {
	t, e := db.BeginTx(ctx, nil)
	if e != nil {
		return e
	}
	if e = fn(t); e != nil {
		_ = t.Rollback()
		return e
	}
	return t.Commit()
}
func (d *DB) WithTx(ctx context.Context, fn func(*sql.Tx) error) error { return tx(ctx, d.SQL, fn) }
func qstr(s string) string                                             { return fmt.Sprintf("%q", s) }
