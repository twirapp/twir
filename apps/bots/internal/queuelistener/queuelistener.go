package queuelistener

import (
	"context"

	"github.com/satont/twir/apps/bots/internal/messagehandler"
	"github.com/satont/twir/libs/logger"
	bus_core "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/twitch"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In
	LC fx.Lifecycle

	MessageHandler *messagehandler.MessageHandler
	Logger         logger.Logger
	Bus            *bus_core.Bus
}

type QueueListener struct {
	messageHandler *messagehandler.MessageHandler
	logger         logger.Logger
	bus            *bus_core.Bus
}

func New(opts Opts) (*QueueListener, error) {
	listener := &QueueListener{
		messageHandler: opts.MessageHandler,
		logger:         opts.Logger,
		bus:            opts.Bus,
	}

	opts.LC.Append(
		fx.Hook{
			OnStart: func(_ context.Context) error {
				return listener.bus.BotsMessages.Subscribe(
					func(ctx context.Context, data twitch.TwitchChatMessage) struct{} {
						if err := listener.messageHandler.Handle(ctx, data); err != nil {
							listener.logger.Error("failed to handle message", "error", err)
						}

						return struct{}{}
					},
				)
			},
			OnStop: func(ctx context.Context) error {
				listener.bus.BotsMessages.Unsubscribe()
				return nil
			},
		},
	)

	return listener, nil
}
