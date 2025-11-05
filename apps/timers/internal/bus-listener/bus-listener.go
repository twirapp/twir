package bus_listener

import (
	"context"

	"github.com/google/uuid"
	"github.com/twirapp/twir/apps/timers/internal/manager"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/timers"
	"github.com/twirapp/twir/libs/bus-core/twitch"
	"github.com/twirapp/twir/libs/logger"
	"go.uber.org/fx"
)

type server struct {
	manager *manager.Manager
}

type Opts struct {
	fx.In

	Lc      fx.Lifecycle
	Logger  logger.Logger
	Bus     *buscore.Bus
	Manager *manager.Manager
}

func New(opts Opts) error {
	s := &server{
		manager: opts.Manager,
	}

	opts.Lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				opts.Bus.Timers.AddTimer.SubscribeGroup("timers", s.onAddTimerToQueue)
				opts.Bus.Timers.RemoveTimer.SubscribeGroup("timers", s.onRemoveTimerFromQueue)
				opts.Bus.ChatMessages.SubscribeGroup("timers", s.onChatMessage)

				return nil
			},
			OnStop: func(ctx context.Context) error {
				opts.Bus.Timers.AddTimer.Unsubscribe()
				opts.Bus.Timers.RemoveTimer.Unsubscribe()
				opts.Bus.ChatMessages.Unsubscribe()
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

func (c *server) onChatMessage(ctx context.Context, m twitch.TwitchChatMessage) (struct{}, error) {
	c.manager.OnChatMessage(m.BroadcasterUserId)
	return struct{}{}, nil
}
