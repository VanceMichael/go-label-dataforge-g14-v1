package repository

import (
	"context"
	"database/sql"
	"github.com/VanceMichael/go-label-dataforge-g14-v1/internal/domain"
	"github.com/VanceMichael/go-label-dataforge-g14-v1/internal/pagination"
	"time"
)

type ResourceQuery interface {
	ListResources(context.Context, string, pagination.Query) ([]domain.Resource, int, error)
	ListAudit(context.Context, string, string, time.Time, int) ([]domain.AuditEvent, error)
}
type SQLResourceQuery struct{ DB *sql.DB }

func (q SQLResourceQuery) ListResources(ctx context.Context, tenant string, p pagination.Query) ([]domain.Resource, int, error) {
	p = pagination.Normalize(p)
	rows, e := q.DB.QueryContext(ctx, "SELECT id,tenant_id,owner_id,code,name,description,status,version,created_at,updated_at FROM data_resources WHERE tenant_id=? ORDER BY updated_at DESC LIMIT ? OFFSET ?", tenant, p.Limit, p.Offset)
	if e != nil {
		return nil, 0, e
	}
	defer rows.Close()
	out := make([]domain.Resource, 0, p.Limit)
	for rows.Next() {
		var r domain.Resource
		var st, cr, up string
		if e = rows.Scan(&r.ID, &r.TenantID, &r.OwnerID, &r.Code, &r.Name, &r.Description, &st, &r.Version, &cr, &up); e != nil {
			return nil, 0, e
		}
		r.Status = domain.ResourceStatus(st)
		r.CreatedAt, _ = time.Parse(time.RFC3339Nano, cr)
		r.UpdatedAt, _ = time.Parse(time.RFC3339Nano, up)
		out = append(out, r)
	}
	var total int
	_ = q.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM data_resources WHERE tenant_id=?", tenant).Scan(&total)
	return out, total, rows.Err()
}
func (q SQLResourceQuery) ListAudit(ctx context.Context, tenant, obj string, after time.Time, limit int) ([]domain.AuditEvent, error) {
	rows, e := q.DB.QueryContext(ctx, "SELECT id,tenant_id,actor_id,object_type,object_id,action,result,request_id,detail,created_at FROM audit_events WHERE tenant_id=? AND object_type=? AND created_at>? ORDER BY created_at LIMIT ?", tenant, obj, after.UTC().Format(time.RFC3339Nano), limit)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := []domain.AuditEvent{}
	for rows.Next() {
		var a domain.AuditEvent
		var cr string
		if e = rows.Scan(&a.ID, &a.TenantID, &a.ActorID, &a.ObjectType, &a.ObjectID, &a.Action, &a.Result, &a.RequestID, &a.Detail, &cr); e != nil {
			return nil, e
		}
		a.CreatedAt, _ = time.Parse(time.RFC3339Nano, cr)
		out = append(out, a)
	}
	return out, rows.Err()
}
