package stream_handlers

import (
	"context"
	"log/slog"

	bus_core "github.com/twirapp/twir/libs/bus-core"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	"github.com/twirapp/twir/libs/repositories/greetings"
	greetingsmodel "github.com/twirapp/twir/libs/repositories/greetings/model"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type PubSubHandlers struct {
	db                  *gorm.DB
	logger              *slog.Logger
	bus                 *bus_core.Bus
	greetingsRepository greetings.Repository
	greetingsCacher     *generic_cacher.GenericCacher[[]greetingsmodel.Greeting]
}

type Opts struct {
	fx.In

	LC fx.Lifecycle

	DB                  *gorm.DB
	Bus                 *bus_core.Bus
	Logger              *slog.Logger
	GreetingsRepository greetings.Repository
	GreetingsCacher     *generic_cacher.GenericCacher[[]greetingsmodel.Greeting]
}

func New(opts Opts) {
	service := &PubSubHandlers{
		db:                  opts.DB,
		logger:              opts.Logger,
		bus:                 opts.Bus,
		greetingsRepository: opts.GreetingsRepository,
		greetingsCacher:     opts.GreetingsCacher,
	}

	opts.LC.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				service.bus.Channel.StreamOnline.SubscribeGroup("bots", service.streamsOnline)
				service.bus.Channel.StreamOffline.SubscribeGroup("bots", service.streamsOffline)
				return nil
			},
			OnStop: func(ctx context.Context) error {
				service.bus.Channel.StreamOnline.Unsubscribe()
				service.bus.Channel.StreamOffline.Unsubscribe()
				return nil
			},
		},
	)
}
