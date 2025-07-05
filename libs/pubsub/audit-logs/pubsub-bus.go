package auditlog

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
	buscore "github.com/twirapp/twir/libs/bus-core"
	busauditlog "github.com/twirapp/twir/libs/bus-core/audit-logs"
	"go.uber.org/fx"
)

type (
	BusPubSub struct {
		bus *buscore.Bus

		subs       map[string]BusSubscription
		subsLocker sync.RWMutex
	}

	BusSubscription struct {
		channel      chan AuditLog
		dashboardIDs []string
		done         chan struct{}
	}
)

var (
	_ PubSub       = (*BusPubSub)(nil)
	_ Subscription = (*BusSubscription)(nil)
)

func NewBusPubSub(bus *buscore.Bus) *BusPubSub {
	return &BusPubSub{
		bus:        bus,
		subs:       make(map[string]BusSubscription),
		subsLocker: sync.RWMutex{},
	}
}

func NewBusPubSubFx(bus *buscore.Bus, lc fx.Lifecycle) *BusPubSub {
	bps := NewBusPubSub(bus)

	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				return bps.Start()
			},
			OnStop: func(ctx context.Context) error {
				bps.Stop()
				return nil
			},
		},
	)

	return bps
}

func (b *BusPubSub) Start() error {
	err := b.bus.AuditLogs.Logs.Subscribe(
		func(ctx context.Context, msg busauditlog.NewAuditLogMessage) (struct{}, error) {
			auditLog := fromBusNewAuditLogMessage(msg)

			b.subsLocker.RLock()
			for _, sub := range b.subs {
				if filter(auditLog, sub.dashboardIDs...) {
					sub.channel <- auditLog
				}
			}
			b.subsLocker.RUnlock()

			return struct{}{}, nil
		},
	)
	if err != nil {
		return fmt.Errorf("subscribe to audit logs from bus: %w", err)
	}

	return nil
}

func (b *BusPubSub) Stop() {
	b.bus.AuditLogs.Logs.Unsubscribe()
}

func (b *BusPubSub) Publish(ctx context.Context, auditLog AuditLog) error {
	auditLogMsg := toBusNewAuditLogMessage(auditLog)

	if err := b.bus.AuditLogs.Logs.Publish(ctx, auditLogMsg); err != nil {
		return fmt.Errorf("publish audit log to bus: %w", err)
	}

	return nil
}

func (b *BusPubSub) Subscribe(_ context.Context, dashboardIDs ...string) (Subscription, error) {
	var (
		sub = BusSubscription{
			channel:      make(chan AuditLog),
			dashboardIDs: dashboardIDs,
			done:         make(chan struct{}),
		}

		subID = uuid.NewString()
	)

	b.subsLocker.Lock()
	b.subs[subID] = sub
	b.subsLocker.Unlock()

	go func() {
		for range sub.done {
			b.subsLocker.Lock()
			delete(b.subs, subID)
			b.subsLocker.Unlock()
		}
	}()

	return sub, nil
}

func (b BusSubscription) Channel() <-chan AuditLog {
	return b.channel
}

func (b BusSubscription) Close() error {
	b.done <- struct{}{}
	close(b.done)
	return nil
}
