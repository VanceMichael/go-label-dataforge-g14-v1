package service

import (
	"context"
	"database/sql"
	"github.com/VanceMichael/go-label-dataforge-g14-v1/internal/apperrors"
	"github.com/VanceMichael/go-label-dataforge-g14-v1/internal/clock"
	"github.com/VanceMichael/go-label-dataforge-g14-v1/internal/domain"
	"github.com/VanceMichael/go-label-dataforge-g14-v1/internal/repository"
)

type Sandbox struct {
	Store repository.Store
	Clock clock.Clock
}

func (s Sandbox) Reserve(ctx context.Context, l domain.Lease) error {
	if l.ExpiresAt.Before(s.Clock.Now()) {
		return apperrors.ErrExpired
	}
	return s.Store.WithTx(ctx, func(tx *sql.Tx) error { return s.Store.CreateLeaseTx(ctx, tx, l) })
}
func (s Sandbox) Run(ctx context.Context, leaseID, input string) (domain.SandboxRun, error) {
	if err := ctx.Err(); err != nil {
		return domain.SandboxRun{}, err
	}
	l, e := s.Store.GetLease(ctx, leaseID)
	if e != nil {
		return domain.SandboxRun{}, e
	}
	if l.Status != domain.LeaseReserved && l.Status != domain.LeaseRunning {
		return domain.SandboxRun{}, apperrors.ErrConflict
	}
	now := s.Clock.Now()
	run := domain.SandboxRun{ID: now.Format("20060102150405.000000000"), LeaseID: leaseID, Status: "completed", Input: input, Output: "accepted", StartedAt: &now, FinishedAt: &now}
	if e := s.Store.WithTx(ctx, func(tx *sql.Tx) error { return s.Store.CreateRunTx(ctx, tx, run) }); e != nil {
		return domain.SandboxRun{}, e
	}
	return run, nil
}
