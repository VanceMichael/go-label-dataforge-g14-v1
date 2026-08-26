package worker

import (
	"context"
	"github.com/VanceMichael/go-label-dataforge-g14-v1/internal/clock"
	"github.com/VanceMichael/go-label-dataforge-g14-v1/internal/domain"
	"github.com/VanceMichael/go-label-dataforge-g14-v1/internal/repository"
	"log/slog"
	"time"
)

type Worker struct {
	Store    repository.Store
	Clock    clock.Clock
	Interval time.Duration
	Log      *slog.Logger
}

func (w Worker) Run(ctx context.Context) error {
	ticker := time.NewTicker(w.Interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			w.tick(ctx)
		}
	}
}
func (w Worker) tick(ctx context.Context) {
	now := w.Clock.Now()
	if n, e := w.Store.ExpireLeases(ctx, now); e != nil {
		w.Log.Error("lease recovery failed", "err", e)
	} else if n > 0 {
		w.Log.Info("leases expired", "count", n)
	}
	if n, e := w.Store.ExpireReviews(ctx, now); e != nil {
		w.Log.Error("review recovery failed", "err", e)
	} else if n > 0 {
		w.Log.Info("reviews expired", "count", n)
	}
	for i := 0; i < 4; i++ {
		j, e := w.Store.ClaimOutbox(ctx, now)
		if e != nil {
			break
		}
		if e = w.dispatch(ctx, j); e != nil {
			_ = w.Store.FinishOutbox(ctx, j.ID, e)
		} else {
			_ = w.Store.FinishOutbox(ctx, j.ID, nil)
		}
	}
}
func (w Worker) dispatch(ctx context.Context, j domain.OutboxJob) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}
