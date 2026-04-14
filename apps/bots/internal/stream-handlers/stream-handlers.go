package stream_handlers

import (
	"context"
	"log/slog"

	bus_core "github.com/twirapp/twir/libs/bus-core"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	channelsrepository "github.com/twirapp/twir/libs/repositories/channels"
	"github.com/twirapp/twir/libs/repositories/greetings"
	greetingsmodel "github.com/twirapp/twir/libs/repositories/greetings/model"
	user_platform_accounts "github.com/twirapp/twir/libs/repositories/user_platform_accounts"
	"go.uber.org/fx"
)

type PubSubHandlers struct {
	logger                   *slog.Logger
	bus                      *bus_core.Bus
	channelsRepo             channelsrepository.Repository
	greetingsRepository      greetings.Repository
	greetingsCacher          *generic_cacher.GenericCacher[[]greetingsmodel.Greeting]
	userPlatformAccountsRepo user_platform_accounts.Repository
}

type Opts struct {
	fx.In

	LC fx.Lifecycle

	Bus                      *bus_core.Bus
	ChannelsRepo             channelsrepository.Repository
	Logger                   *slog.Logger
	GreetingsRepository      greetings.Repository
	GreetingsCacher          *generic_cacher.GenericCacher[[]greetingsmodel.Greeting]
	UserPlatformAccountsRepo user_platform_accounts.Repository
}

func New(opts Opts) {
	service := &PubSubHandlers{
		logger:                   opts.Logger,
		bus:                      opts.Bus,
		channelsRepo:             opts.ChannelsRepo,
		greetingsRepository:      opts.GreetingsRepository,
		greetingsCacher:          opts.GreetingsCacher,
		userPlatformAccountsRepo: opts.UserPlatformAccountsRepo,
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
