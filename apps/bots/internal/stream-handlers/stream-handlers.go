package stream_handlers

import (
	"context"

	"github.com/satont/twir/libs/logger"
	bus_core "github.com/twirapp/twir/libs/bus-core"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type PubSubHandlers struct {
	db     *gorm.DB
	logger logger.Logger
	bus    *bus_core.Bus
}

type Opts struct {
	fx.In

	LC fx.Lifecycle

	DB     *gorm.DB
	Bus    *bus_core.Bus
	Logger logger.Logger
}

func New(opts Opts) {
	service := &PubSubHandlers{
		db:     opts.DB,
		logger: opts.Logger,
		bus:    opts.Bus,
	}

	opts.LC.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				service.bus.StreamOnline.SubscribeGroup("bots", service.streamsOnline)
				service.bus.StreamOffline.SubscribeGroup("bots", service.streamsOffline)
				return nil
			},
			OnStop: func(ctx context.Context) error {
				service.bus.StreamOnline.Unsubscribe()
				service.bus.StreamOffline.Unsubscribe()
				return nil
			},
		},
	)
}
