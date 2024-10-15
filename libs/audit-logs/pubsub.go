package auditlog

import (
	"context"
)

type PubSub interface {
	Publish(ctx context.Context, auditLog AuditLog) error
	// Subscribe subscribes for new audit logs and filter them by the provided IDs of users.
	Subscribe(ctx context.Context, userIDs ...string) (Subscription, error)
}

type Subscription interface {
	Channel() <-chan AuditLog
	Close() error
}
