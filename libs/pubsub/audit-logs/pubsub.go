package auditlog

import (
	"context"
)

type PubSub interface {
	Publish(ctx context.Context, auditLog AuditLog) error
	// Subscribe subscribes for new audit logs and filter them by the provided IDs of dashboards.
	Subscribe(ctx context.Context, dashboardIDs ...string) (Subscription, error)
}

type Subscription interface {
	Channel() <-chan AuditLog
	Close() error
}
