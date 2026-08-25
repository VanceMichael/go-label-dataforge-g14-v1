package outbox

import (
	"context"
	"github.com/VanceMichael/go-label-dataforge-g14-v1/internal/clock"
	"github.com/VanceMichael/go-label-dataforge-g14-v1/internal/domain"
	"github.com/VanceMichael/go-label-dataforge-g14-v1/internal/repository"
	"time"
)

type Dispatcher struct {
	Store repository.Store
	Clock clock.Clock
}

func (d Dispatcher) Enqueue(ctx context.Context, tx interface{}, job domain.OutboxJob) error {
	_ = tx
	_ = ctx
	_ = job
	return nil
}
func NewJob(id, topic, payload string, c clock.Clock) domain.OutboxJob {
	now := c.Now()
	return domain.OutboxJob{ID: id, Topic: topic, Payload: payload, AvailableAt: now, CreatedAt: now}
}

var _ = time.Second
