package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/VanceMichael/go-label-dataforge-g14-v1/internal/apperrors"
	"github.com/VanceMichael/go-label-dataforge-g14-v1/internal/domain"
	"time"
)

type Repo struct{ DB *DB }

func NewRepo(db *DB) *Repo                                               { return &Repo{DB: db} }
func iso(t time.Time) string                                             { return t.UTC().Format(time.RFC3339Nano) }
func parse(s string) time.Time                                           { t, _ := time.Parse(time.RFC3339Nano, s); return t }
func (r *Repo) WithTx(ctx context.Context, fn func(*sql.Tx) error) error { return r.DB.WithTx(ctx, fn) }
func (r *Repo) CreateTenant(ctx context.Context, tx *sql.Tx, t domain.Tenant) error {
	_, e := tx.ExecContext(ctx, "INSERT INTO tenants(id,name,created_at) VALUES(?,?,?)", t.ID, t.Name, iso(t.CreatedAt))
	return e
}
func (r *Repo) CreateUser(ctx context.Context, tx *sql.Tx, u domain.User) error {
	_, e := tx.ExecContext(ctx, "INSERT INTO users(id,tenant_id,email,role,active,created_at) VALUES(?,?,?,?,?,?)", u.ID, u.TenantID, u.Email, u.Role, u.Active, iso(u.CreatedAt))
	return e
}
func (r *Repo) FindUserByEmail(ctx context.Context, tenant, email string) (domain.User, error) {
	var u domain.User
	var role string
	var active int
	var created string
	e := r.DB.SQL.QueryRowContext(ctx, "SELECT id,tenant_id,email,role,active,created_at FROM users WHERE tenant_id=? AND email=?", tenant, email).Scan(&u.ID, &u.TenantID, &u.Email, &role, &active, &created)
	if errors.Is(e, sql.ErrNoRows) {
		return u, apperrors.ErrNotFound
	}
	u.Role = domain.Role(role)
	u.Active = active == 1
	u.CreatedAt = parse(created)
	return u, e
}
func (r *Repo) CreateSession(ctx context.Context, s domain.Session) error {
	_, e := r.DB.SQL.ExecContext(ctx, "INSERT INTO sessions(id,user_id,token_hash,expires_at,created_at) VALUES(?,?,?,?,?)", s.ID, s.UserID, s.TokenHash, iso(*s.ExpiresAt), iso(*s.CreatedAt))
	return e
}
func (r *Repo) FindSession(ctx context.Context, id string) (domain.Session, error) {
	var s domain.Session
	var ex, cr, rv sql.NullString
	e := r.DB.SQL.QueryRowContext(ctx, "SELECT id,user_id,token_hash,expires_at,created_at,revoked_at FROM sessions WHERE id=?", id).Scan(&s.ID, &s.UserID, &s.TokenHash, &ex, &cr, &rv)
	if errors.Is(e, sql.ErrNoRows) {
		return s, apperrors.ErrNotFound
	}
	x := parse(ex.String)
	c := parse(cr.String)
	s.ExpiresAt = &x
	s.CreatedAt = &c
	if rv.Valid {
		v := parse(rv.String)
		s.RevokedAt = &v
	}
	return s, e
}
func (r *Repo) RevokeSession(ctx context.Context, id string) error {
	res, e := r.DB.SQL.ExecContext(ctx, "UPDATE sessions SET revoked_at=datetime('now') WHERE id=? AND revoked_at IS NULL", id)
	if e == nil {
		n, _ := res.RowsAffected()
		if n == 0 {
			return apperrors.ErrNotFound
		}
	}
	return e
}
func (r *Repo) CreateResourceTx(ctx context.Context, tx *sql.Tx, res domain.Resource, ver domain.ResourceVersion, rev domain.Review) error {
	if _, e := tx.ExecContext(ctx, "INSERT INTO data_resources(id,tenant_id,owner_id,code,name,description,status,version,created_at,updated_at) VALUES(?,?,?,?,?,?,?,?,?,?)", res.ID, res.TenantID, res.OwnerID, res.Code, res.Name, res.Description, res.Status, res.Version, iso(res.CreatedAt), iso(res.UpdatedAt)); e != nil {
		return e
	}
	if _, e := tx.ExecContext(ctx, "INSERT INTO resource_versions(id,resource_id,number,schema,content_hash,created_at) VALUES(?,?,?,?,?,?)", ver.ID, ver.ResourceID, ver.Number, ver.Schema, ver.ContentHash, iso(ver.CreatedAt)); e != nil {
		return e
	}
	_, e := tx.ExecContext(ctx, "INSERT INTO registration_reviews(id,resource_id,reviewer_id,decision,notes,status,created_at) VALUES(?,?,?,?,?,?,?)", rev.ID, rev.ResourceID, nil, rev.Decision, rev.Notes, rev.Status, iso(rev.CreatedAt))
	return e
}
func (r *Repo) GetResource(ctx context.Context, id string) (domain.Resource, error) {
	var x domain.Resource
	var st, cr, up string
	e := r.DB.SQL.QueryRowContext(ctx, "SELECT id,tenant_id,owner_id,code,name,description,status,version,created_at,updated_at FROM data_resources WHERE id=?", id).Scan(&x.ID, &x.TenantID, &x.OwnerID, &x.Code, &x.Name, &x.Description, &st, &x.Version, &cr, &up)
	if errors.Is(e, sql.ErrNoRows) {
		return x, apperrors.ErrNotFound
	}
	x.Status = domain.ResourceStatus(st)
	x.CreatedAt = parse(cr)
	x.UpdatedAt = parse(up)
	return x, e
}
func (r *Repo) UpdateResourceStatusTx(ctx context.Context, tx *sql.Tx, id string, status domain.ResourceStatus, version int) error {
	res, e := tx.ExecContext(ctx, "UPDATE data_resources SET status=?,version=version+1,updated_at=datetime('now') WHERE id=? AND version=?", status, id, version)
	if e != nil {
		return e
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return apperrors.ErrConflict
	}
	return nil
}
func (r *Repo) CreateAuthorizationTx(ctx context.Context, tx *sql.Tx, a domain.Authorization) error {
	_, e := tx.ExecContext(ctx, "INSERT INTO authorization_requests(id,tenant_id,resource_id,applicant_id,status,purpose,quota,used,version,created_at,updated_at) VALUES(?,?,?,?,?,?,?,?,?,?,?)", a.ID, a.TenantID, a.ResourceID, a.ApplicantID, a.Status, a.Purpose, a.Quota, a.Used, a.Version, iso(a.CreatedAt), iso(a.UpdatedAt))
	return e
}
func (r *Repo) GetAuthorization(ctx context.Context, id string) (domain.Authorization, error) {
	var a domain.Authorization
	var st, cr, up string
	e := r.DB.SQL.QueryRowContext(ctx, "SELECT id,tenant_id,resource_id,applicant_id,status,purpose,quota,used,version,created_at,updated_at FROM authorization_requests WHERE id=?", id).Scan(&a.ID, &a.TenantID, &a.ResourceID, &a.ApplicantID, &st, &a.Purpose, &a.Quota, &a.Used, &a.Version, &cr, &up)
	if errors.Is(e, sql.ErrNoRows) {
		return a, apperrors.ErrNotFound
	}
	a.Status = domain.AuthorizationStatus(st)
	a.CreatedAt = parse(cr)
	a.UpdatedAt = parse(up)
	return a, e
}
func (r *Repo) ApproveAuthorizationTx(ctx context.Context, tx *sql.Tx, id string, version int) error {
	res, e := tx.ExecContext(ctx, "UPDATE authorization_requests SET status='active',version=version+1,updated_at=datetime('now') WHERE id=? AND status='approved' AND version=?", id, version)
	if e != nil {
		return e
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return apperrors.ErrConflict
	}
	return nil
}
func (r *Repo) CreateLeaseTx(ctx context.Context, tx *sql.Tx, l domain.Lease) error {
	_, e := tx.ExecContext(ctx, "INSERT INTO sandbox_leases(id,authorization_id,holder_id,status,expires_at,version,created_at) VALUES(?,?,?,?,?,?,?)", l.ID, l.AuthorizationID, l.HolderID, l.Status, iso(l.ExpiresAt), l.Version, iso(l.CreatedAt))
	return e
}
func (r *Repo) GetLease(ctx context.Context, id string) (domain.Lease, error) {
	var l domain.Lease
	var st, ex, cr string
	e := r.DB.SQL.QueryRowContext(ctx, "SELECT id,authorization_id,holder_id,status,expires_at,version,created_at FROM sandbox_leases WHERE id=?", id).Scan(&l.ID, &l.AuthorizationID, &l.HolderID, &st, &ex, &l.Version, &cr)
	if errors.Is(e, sql.ErrNoRows) {
		return l, apperrors.ErrNotFound
	}
	l.Status = domain.LeaseStatus(st)
	l.ExpiresAt = parse(ex)
	l.CreatedAt = parse(cr)
	return l, e
}
func (r *Repo) UpdateLeaseTx(ctx context.Context, tx *sql.Tx, id string, status domain.LeaseStatus, version int) error {
	res, e := tx.ExecContext(ctx, "UPDATE sandbox_leases SET status=?,version=version+1 WHERE id=? AND version=?", status, id, version)
	if e != nil {
		return e
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return apperrors.ErrConflict
	}
	return nil
}
func (r *Repo) CreateRunTx(ctx context.Context, tx *sql.Tx, run domain.SandboxRun) error {
	_, e := tx.ExecContext(ctx, "INSERT INTO sandbox_runs(id,lease_id,status,input,output) VALUES(?,?,?,?,?)", run.ID, run.LeaseID, run.Status, run.Input, run.Output)
	return e
}
func (r *Repo) CreateProductTx(ctx context.Context, tx *sql.Tx, p domain.Product) error {
	_, e := tx.ExecContext(ctx, "INSERT INTO data_products(id,tenant_id,owner_id,name,resource_id,status,version,created_at,updated_at) VALUES(?,?,?,?,?,?,?,?,?)", p.ID, p.TenantID, p.OwnerID, p.Name, p.ResourceID, p.Status, p.Version, iso(p.CreatedAt), iso(p.UpdatedAt))
	return e
}
func (r *Repo) GetProduct(ctx context.Context, id string) (domain.Product, error) {
	var p domain.Product
	var st, cr, up string
	e := r.DB.SQL.QueryRowContext(ctx, "SELECT id,tenant_id,owner_id,name,resource_id,status,version,created_at,updated_at FROM data_products WHERE id=?", id).Scan(&p.ID, &p.TenantID, &p.OwnerID, &p.Name, &p.ResourceID, &st, &p.Version, &cr, &up)
	if errors.Is(e, sql.ErrNoRows) {
		return p, apperrors.ErrNotFound
	}
	p.Status = domain.ProductStatus(st)
	p.CreatedAt = parse(cr)
	p.UpdatedAt = parse(up)
	return p, e
}
func (r *Repo) PublishProductTx(ctx context.Context, tx *sql.Tx, id string, version int, rel domain.ProductRelease) error {
	res, e := tx.ExecContext(ctx, "UPDATE data_products SET status='published',version=version+1,updated_at=datetime('now') WHERE id=? AND status='pending_release' AND version=?", id, version)
	if e != nil {
		return e
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return apperrors.ErrConflict
	}
	_, e = tx.ExecContext(ctx, "INSERT INTO product_releases(id,product_id,version,status,published_at,created_at) VALUES(?,?,?,?,?,?)", rel.ID, rel.ProductID, rel.Version, rel.Status, iso(*rel.PublishedAt), iso(rel.CreatedAt))
	return e
}
func (r *Repo) AddAuditTx(ctx context.Context, tx *sql.Tx, a domain.AuditEvent) error {
	_, e := tx.ExecContext(ctx, "INSERT INTO audit_events(id,tenant_id,actor_id,object_type,object_id,action,result,request_id,detail,created_at) VALUES(?,?,?,?,?,?,?,?,?,?)", a.ID, a.TenantID, a.ActorID, a.ObjectType, a.ObjectID, a.Action, a.Result, a.RequestID, a.Detail, iso(a.CreatedAt))
	return e
}
func (r *Repo) EnqueueTx(ctx context.Context, tx *sql.Tx, j domain.OutboxJob) error {
	_, e := tx.ExecContext(ctx, "INSERT INTO outbox_jobs(id,topic,payload,attempts,available_at,status,last_error,created_at) VALUES(?,?,?,?,?,'pending','',?)", j.ID, j.Topic, j.Payload, j.Attempts, iso(j.AvailableAt), iso(j.CreatedAt))
	return e
}
func (r *Repo) ClaimOutbox(ctx context.Context, now time.Time) (domain.OutboxJob, error) {
	var j domain.OutboxJob
	err := r.DB.WithTx(ctx, func(tx *sql.Tx) error {
		var av, cr string
		e := tx.QueryRowContext(ctx, "SELECT id,topic,payload,attempts,available_at,created_at FROM outbox_jobs WHERE status='pending' AND available_at<=? ORDER BY created_at LIMIT 1", iso(now)).Scan(&j.ID, &j.Topic, &j.Payload, &j.Attempts, &av, &cr)
		if e != nil {
			return e
		}
		res, e := tx.ExecContext(ctx, "UPDATE outbox_jobs SET status='processing',locked_until=? WHERE id=? AND status='pending'", iso(now.Add(time.Minute)), j.ID)
		if e != nil {
			return e
		}
		n, _ := res.RowsAffected()
		if n == 0 {
			return apperrors.ErrConflict
		}
		j.AvailableAt = parse(av)
		j.CreatedAt = parse(cr)
		j.Status = "processing"
		return nil
	})
	return j, err
}
func (r *Repo) FinishOutbox(ctx context.Context, id string, jobErr error) error {
	if jobErr == nil {
		_, e := r.DB.SQL.ExecContext(ctx, "UPDATE outbox_jobs SET status='done',locked_until=NULL WHERE id=?", id)
		return e
	}
	_, e := r.DB.SQL.ExecContext(ctx, "UPDATE outbox_jobs SET status=CASE WHEN attempts+1>=5 THEN 'failed' ELSE 'pending' END,attempts=attempts+1,last_error=?,available_at=datetime('now','+5 seconds'),locked_until=NULL WHERE id=?", jobErr.Error(), id)
	return e
}
func (r *Repo) ExpireLeases(ctx context.Context, now time.Time) (int, error) {
	res, e := r.DB.SQL.ExecContext(ctx, "UPDATE sandbox_leases SET status='expired',version=version+1 WHERE status IN ('reserved','running') AND expires_at<?", iso(now))
	if e != nil {
		return 0, e
	}
	n, _ := res.RowsAffected()
	return int(n), nil
}
func (r *Repo) ExpireReviews(ctx context.Context, now time.Time) (int, error) {
	res, e := r.DB.SQL.ExecContext(ctx, "UPDATE registration_reviews SET status='expired' WHERE status='claimed' AND lease_until<?", iso(now))
	if e != nil {
		return 0, e
	}
	n, _ := res.RowsAffected()
	return int(n), nil
}
func ensure(e error) error {
	if e == sql.ErrNoRows {
		return apperrors.ErrNotFound
	}
	return e
}

var _ = fmt.Sprintf
