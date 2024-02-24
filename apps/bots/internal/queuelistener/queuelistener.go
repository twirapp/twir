package queuelistener

import (
	"context"

	"github.com/satont/twir/apps/bots/internal/messagehandler"
	"github.com/satont/twir/libs/logger"
	"github.com/satont/twir/libs/types/types/services"
	"github.com/satont/twir/libs/types/types/services/twitch"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In
	LC fx.Lifecycle

	MessageHandler *messagehandler.MessageHandler
	Logger         logger.Logger
	Bus            *services.Bus
}

type QueueListener struct {
	messageHandler *messagehandler.MessageHandler
	logger         logger.Logger
	bus            *services.Bus
}

func New(opts Opts) (*QueueListener, error) {
	listener := &QueueListener{
		messageHandler: opts.MessageHandler,
		logger:         opts.Logger,
		bus:            opts.Bus,
	}

	opts.LC.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				return listener.subscribe()
			},
			OnStop: func(ctx context.Context) error {
				listener.unsubscribe()
				return nil
			},
		},
	)

	return listener, nil
}

func (c *QueueListener) subscribe() error {
	c.bus.BotsMessages.Subscribe(
		func(ctx context.Context, data twitch.TwitchChatMessage) struct{} {
			if err := c.messageHandler.Handle(ctx, data); err != nil {
				c.logger.Error("failed to handle message", "error", err)
			}

			return struct{}{}
		},
	)

	return nil
}

func (c *QueueListener) unsubscribe() {
	c.bus.BotsMessages.Unsubscribe()
}
