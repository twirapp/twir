package bus_listener

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/twirapp/twir/apps/timers/internal/manager"
	"github.com/twirapp/twir/libs/baseapp/lifecycle"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/generic"
	"github.com/twirapp/twir/libs/bus-core/timers"
)

type server struct {
	manager *manager.Manager
}

func New(
	lc *lifecycle.Lifecycle,
	_ *slog.Logger,
	bus *buscore.Bus,
	managerService *manager.Manager,
) error {
	s := &server{
		manager: managerService,
	}

	lc.Append(
		lifecycle.Hook{
			OnStart: func(ctx context.Context) error {
				bus.Timers.AddTimer.SubscribeGroup("timers", s.onAddTimerToQueue)
				bus.Timers.RemoveTimer.SubscribeGroup("timers", s.onRemoveTimerFromQueue)
				bus.ChatMessages.SubscribeGroup("timers", s.onChatMessage)

				return nil
			},
			OnStop: func(ctx context.Context) error {
				bus.Timers.AddTimer.Unsubscribe()
				bus.Timers.RemoveTimer.Unsubscribe()
				bus.ChatMessages.Unsubscribe()
				return nil
			},
		},
	)

	return nil
}

func (c *server) onAddTimerToQueue(ctx context.Context, t timers.AddOrRemoveTimerRequest) (
	struct{},
	error,
) {
	return struct{}{}, c.manager.AddTimerById(ctx, manager.TimerID(uuid.MustParse(t.TimerID)))
}

func (c *server) onRemoveTimerFromQueue(
	ctx context.Context,
	t timers.AddOrRemoveTimerRequest,
) (struct{}, error) {
	c.manager.RemoveTimerById(manager.TimerID(uuid.MustParse(t.TimerID)))
	return struct{}{}, nil
}

func (c *server) onChatMessage(ctx context.Context, m generic.ChatMessage) (struct{}, error) {
	channelID, err := uuid.Parse(m.ChannelID)
	if err != nil {
		return struct{}{}, fmt.Errorf("parse message channel id: %w", err)
	}

	c.manager.OnChatMessage(channelID)
	return struct{}{}, nil
}
