package pubsub_handlers

import (
	"context"
	"github.com/satont/twir/apps/bots/internal/bots"
	"github.com/satont/twir/libs/logger"
	"github.com/satont/twir/libs/pubsub"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type handlers struct {
	db          *gorm.DB
	logger      logger.Logger
	botsService *bots.Service
}

type Opts struct {
	fx.In

	LC fx.Lifecycle

	DB          *gorm.DB
	PubSub      *pubsub.PubSub
	Logger      logger.Logger
	BotsService *bots.Service
}

func New(opts Opts) {
	service := &handlers{
		db:          opts.DB,
		logger:      opts.Logger,
		botsService: opts.BotsService,
	}

	opts.LC.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				opts.PubSub.Subscribe(
					"user.update", func(data []byte) {
						service.userUpdate(data)
					},
				)
				opts.PubSub.Subscribe(
					"stream.online", func(data []byte) {
						service.streamsOnline(data)
					},
				)
				opts.PubSub.Subscribe(
					"stream.offline", func(data []byte) {
						service.streamsOffline(data)
					},
				)
				return nil
			},
			OnStop: nil,
		},
	)
}
