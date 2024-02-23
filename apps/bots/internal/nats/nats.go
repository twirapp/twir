package nats

import (
	"context"

	"github.com/nats-io/nats.go"
	"github.com/satont/twir/apps/bots/internal/messagehandler"
	cfg "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/logger"
	"github.com/satont/twir/libs/types/types/services"
	"github.com/satont/twir/libs/types/types/services/twitch"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In
	LC fx.Lifecycle

	Config         cfg.Config
	MessageHandler *messagehandler.MessageHandler
	Logger         logger.Logger
}

func New(opts Opts) (*nats.Conn, error) {
	nc, err := nats.Connect(opts.Config.NatsUrl)
	if err != nil {
		return nil, err
	}

	chatSub, err := nc.Subscribe(
		twitch.TOPIC_CHAT_MESSAGE,
		func(msg *nats.Msg) {
			msg.InProgress()
			data, _ := services.Decode[twitch.TwitchChatMessage](msg.Data)
			if err := opts.MessageHandler.Handle(context.Background(), data); err != nil {
				opts.Logger.Error("failed to handle message", "error", err)
			}
			msg.Ack()
		},
	)
	if err != nil {
		return nil, err
	}

	opts.LC.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				return nil
			},
			OnStop: func(ctx context.Context) error {
				chatSub.Unsubscribe()
				nc.Close()
				return nil
			},
		},
	)

	return nc, nil
}
