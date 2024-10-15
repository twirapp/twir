package auditlog

import (
	"context"
	"fmt"
	"slices"

	buscore "github.com/twirapp/twir/libs/bus-core"
	busauditlog "github.com/twirapp/twir/libs/bus-core/audit-logs"
)

type (
	BusPubSub struct {
		bus *buscore.Bus
	}

	BusSubscription struct {
		bus     *buscore.Bus
		channel <-chan AuditLog
	}
)

var (
	_ PubSub       = (*BusPubSub)(nil)
	_ Subscription = (*BusSubscription)(nil)
)

func NewBusPubSub(bus *buscore.Bus) BusPubSub {
	return BusPubSub{bus: bus}
}

func (b BusPubSub) Publish(_ context.Context, auditLog AuditLog) error {
	if err := b.bus.AuditLogs.Logs.Publish(toBusNewAuditLogMessage(auditLog)); err != nil {
		return fmt.Errorf("publish audit log to bus: %w", err)
	}

	return nil
}

func (b BusPubSub) Subscribe(_ context.Context, userIDs ...string) (Subscription, error) {
	channel := make(chan AuditLog)

	err := b.bus.AuditLogs.Logs.Subscribe(
		func(ctx context.Context, msg busauditlog.NewAuditLogMessage) struct{} {
			if len(userIDs) != 0 && msg.UserID == nil {
				return struct{}{}
			}

			if msg.UserID != nil {
				if !slices.Contains(userIDs, *msg.UserID) {
					return struct{}{}
				}
			}

			channel <- fromBusNewAuditLogMessage(msg)

			return struct{}{}
		},
	)
	if err != nil {
		return nil, fmt.Errorf("subscribe to audit logs from bus: %w", err)
	}

	return BusSubscription{
		bus:     b.bus,
		channel: channel,
	}, nil
}

func (b BusSubscription) Channel() <-chan AuditLog {
	return b.channel
}

func (b BusSubscription) Close() error {
	b.bus.AuditLogs.Logs.Unsubscribe()
	return nil
}
