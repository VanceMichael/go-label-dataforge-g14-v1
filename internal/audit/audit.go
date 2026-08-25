package audit

import (
	"github.com/VanceMichael/go-label-dataforge-g14-v1/internal/clock"
	"github.com/VanceMichael/go-label-dataforge-g14-v1/internal/domain"
)

func Event(id, tenant, actor, obj, objID, action, result, req, detail string, now clock.Clock) domain.AuditEvent {
	return domain.AuditEvent{ID: id, TenantID: tenant, ActorID: actor, ObjectType: obj, ObjectID: objID, Action: action, Result: result, RequestID: req, Detail: detail, CreatedAt: now.Now()}
}
