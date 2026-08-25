package service

import (
	"context"
	"github.com/VanceMichael/go-label-dataforge-g14-v1/internal/apperrors"
	"github.com/VanceMichael/go-label-dataforge-g14-v1/internal/clock"
	"github.com/VanceMichael/go-label-dataforge-g14-v1/internal/repository"
	"time"
)

type Recovery struct {
	Store repository.Store
	Clock clock.Clock
}

func (r Recovery) Reclaim(ctx context.Context) (int, int, error) {
	now := r.Clock.Now()
	leases, e := r.Store.ExpireLeases(ctx, now)
	if e != nil {
		return 0, 0, e
	}
	reviews, e := r.Store.ExpireReviews(ctx, now)
	if e != nil {
		return leases, 0, e
	}
	return leases, reviews, nil
}
func (r Recovery) EnsureLease(ctx context.Context, id string) error {
	l, e := r.Store.GetLease(ctx, id)
	if e != nil {
		return e
	}
	if l.ExpiresAt.Before(r.Clock.Now()) {
		return apperrors.ErrExpired
	}
	return nil
}
func Backoff(attempt int) time.Duration {
	if attempt < 0 {
		attempt = 0
	}
	if attempt > 8 {
		attempt = 8
	}
	return time.Duration(1<<attempt) * time.Second
}
