package service

import (
	"context"
	"database/sql"
	"github.com/VanceMichael/go-label-dataforge-g14-v1/internal/apperrors"
	"github.com/VanceMichael/go-label-dataforge-g14-v1/internal/clock"
	"github.com/VanceMichael/go-label-dataforge-g14-v1/internal/domain"
	"github.com/VanceMichael/go-label-dataforge-g14-v1/internal/repository"
)

type Authorizer struct {
	Store repository.Store
	Clock clock.Clock
}

func (s Authorizer) Request(ctx context.Context, a domain.Authorization) error {
	if a.Status == "" {
		a.Status = domain.AuthRequested
	}
	return s.Store.WithTx(ctx, func(tx *sql.Tx) error { return s.Store.CreateAuthorizationTx(ctx, tx, a) })
}
func (s Authorizer) Approve(ctx context.Context, id string, version int) error {
	a, e := s.Store.GetAuthorization(ctx, id)
	if e != nil {
		return e
	}
	if a.Version != version || a.Status != domain.AuthRequested {
		return apperrors.ErrConflict
	}
	return s.Store.WithTx(ctx, func(tx *sql.Tx) error { return s.Store.ApproveAuthorizationTx(ctx, tx, id, version) })
}
