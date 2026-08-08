package stream_handlers

import (
	"context"
	"log/slog"

	"github.com/twirapp/twir/libs/baseapp/lifecycle"
	bus_core "github.com/twirapp/twir/libs/bus-core"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	"github.com/twirapp/twir/libs/repositories/greetings"
	greetingsmodel "github.com/twirapp/twir/libs/repositories/greetings/model"
	channelservice "github.com/twirapp/twir/libs/services/channels"
)

type PubSubHandlers struct {
	logger              *slog.Logger
	bus                 *bus_core.Bus
	channelService      *channelservice.ChannelService
	greetingsRepository greetings.Repository
	greetingsCacher     *generic_cacher.GenericCacher[[]greetingsmodel.Greeting]
}

func New(
	lc *lifecycle.Lifecycle,
	bus *bus_core.Bus,
	channelService *channelservice.ChannelService,
	logger *slog.Logger,
	greetingsRepository greetings.Repository,
	greetingsCacher *generic_cacher.GenericCacher[[]greetingsmodel.Greeting],
) {
	service := &PubSubHandlers{
		logger:              logger,
		bus:                 bus,
		channelService:      channelService,
		greetingsRepository: greetingsRepository,
		greetingsCacher:     greetingsCacher,
	}

	lc.Append(
		lifecycle.Hook{
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
