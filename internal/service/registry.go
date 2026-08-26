package service

import (
	"context"
	"database/sql"
	"github.com/VanceMichael/go-label-dataforge-g14-v1/internal/apperrors"
	"github.com/VanceMichael/go-label-dataforge-g14-v1/internal/audit"
	"github.com/VanceMichael/go-label-dataforge-g14-v1/internal/clock"
	"github.com/VanceMichael/go-label-dataforge-g14-v1/internal/domain"
	"github.com/VanceMichael/go-label-dataforge-g14-v1/internal/repository"
)

type Registry struct {
	Store repository.Store
	Clock clock.Clock
}

func (s Registry) Create(ctx context.Context, res domain.Resource, ver domain.ResourceVersion, rev domain.Review, req string) (domain.Resource, error) {
	if res.Status == "" {
		res.Status = domain.ResourceDraft
	}
	if ver.Number == 0 {
		ver.Number = 1
	}
	if rev.Status == "" {
		rev.Status = "pending"
	}
	if e := s.Store.WithTx(ctx, func(tx *sql.Tx) error { return s.Store.CreateResourceTx(ctx, tx, res, ver, rev) }); e != nil {
		return res, e
	}
	return res, nil
}
func (s Registry) Submit(ctx context.Context, id string, version int, actor, tenant, req string) error {
	res, e := s.Store.GetResource(ctx, id)
	if e != nil {
		return e
	}
	if res.Version != version || !res.Status.Can(domain.ResourceSubmitted) {
		return apperrors.ErrConflict
	}
	event := audit.Event("audit-submit", tenant, actor, "resource", id, "submit", "ok", req, "resource submitted", s.Clock)
	// The status transition and the audit record must commit or roll back
	// together: committing the status in its own transaction and then writing
	// the audit in a second one leaves the resource submitted with no audit
	// (and no way to retry) when the audit write fails.
	return s.Store.WithTx(ctx, func(tx *sql.Tx) error {
		if e := s.Store.UpdateResourceStatusTx(ctx, tx, id, domain.ResourceSubmitted, version); e != nil {
			return e
		}
		return s.Store.AddAuditTx(ctx, tx, event)
	})
}
