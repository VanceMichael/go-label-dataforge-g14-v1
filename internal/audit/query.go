package audit

import (
	"context"
	"github.com/VanceMichael/go-label-dataforge-g14-v1/internal/domain"
	"github.com/VanceMichael/go-label-dataforge-g14-v1/internal/repository"
	"time"
)

type Query struct{ Store repository.ResourceQuery }

func (q Query) Since(ctx context.Context, tenant, object string, after time.Time, limit int) ([]domain.AuditEvent, error) {
	return q.Store.ListAudit(ctx, tenant, object, after, limit)
}
func Compact(events []domain.AuditEvent) map[string]int {
	m := map[string]int{}
	for _, e := range events {
		m[e.Action]++
	}
	return m
}
