package twirbus

import (
	"context"

	"github.com/satont/twir/libs/logger"
	"github.com/twirapp/twir/apps/chat-translator/internal/services/handle_message"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In
	LC fx.Lifecycle

	BusCore              *buscore.Bus
	HandleMessageService *handle_message.Service
	Logger               logger.Logger
}

func New(opts Opts) {
	opts.LC.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				if err := opts.BusCore.ChatMessages.SubscribeGroup(
					"chat-translator",
					opts.HandleMessageService.Handle,
				); err != nil {
					return nil
				}

				opts.Logger.Info("Subscribed to messages")

				return nil
			},
			OnStop: func(ctx context.Context) error {
				opts.BusCore.ChatMessages.Unsubscribe()
				return nil
			},
		},
	)
}
