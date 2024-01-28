package pubsub_handlers

import (
	"context"

	"github.com/satont/twir/libs/logger"
	"github.com/satont/twir/libs/pubsub"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type PubSubHandlers struct {
	db     *gorm.DB
	logger logger.Logger
}

type Opts struct {
	fx.In

	LC fx.Lifecycle

	DB     *gorm.DB
	PubSub *pubsub.PubSub
	Logger logger.Logger
}

func New(opts Opts) {
	service := &PubSubHandlers{
		db:     opts.DB,
		logger: opts.Logger,
	}

	opts.LC.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
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
