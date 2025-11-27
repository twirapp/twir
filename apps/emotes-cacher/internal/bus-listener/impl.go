package bus_listener

import (
	"context"
	"log/slog"
	"slices"

	"github.com/twirapp/twir/apps/emotes-cacher/internal/emotes_store"
	buscore "github.com/twirapp/twir/libs/bus-core"
	emotes_cacher "github.com/twirapp/twir/libs/bus-core/emotes-cacher"
	"go.uber.org/fx"
)

type BusListener struct {
	logger *slog.Logger
	bus    *buscore.Bus
	store  *emotes_store.EmotesStore
}

type Opts struct {
	fx.In
	Lc fx.Lifecycle

	Logger *slog.Logger
	Bus    *buscore.Bus
	Store  *emotes_store.EmotesStore
}

func New(opts Opts) error {
	impl := &BusListener{
		logger: opts.Logger,
		bus:    opts.Bus,
		store:  opts.Store,
	}

	opts.Lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				if err := impl.bus.EmotesCacher.GetChannelEmotes.SubscribeGroup(
					emotes_cacher.EmotesCacherGetChannelEmotesSubject,
					impl.GetChannelEmotes,
				); err != nil {
					return err
				}

				if err := impl.bus.EmotesCacher.GetGlobalEmotes.SubscribeGroup(
					emotes_cacher.EmotesCacherGetGlobalEmotesSubject,
					func(
						ctx context.Context,
						data emotes_cacher.GetGlobalEmotesRequest,
					) (emotes_cacher.Response, error) {
						return impl.GetChannelEmotes(
							ctx,
							emotes_cacher.GetChannelEmotesRequest{
								ChannelID: emotes_store.GlobalChannelID,
							},
						)
					},
				); err != nil {
					return err
				}

				return nil
			},
			OnStop: func(ctx context.Context) error {
				return nil
			},
		},
	)

	return nil
}

func (c *BusListener) GetChannelEmotes(
	_ context.Context,
	request emotes_cacher.GetChannelEmotesRequest,
) (emotes_cacher.Response, error) {
	var result []emotes_cacher.Emote

	for serviceName, emotes := range c.store.GetChannelEmotesServices(emotes_store.ChannelID(request.ChannelID)) {
		if len(request.ServiceIn) > 0 && !slices.Contains(request.ServiceIn, serviceName) {
			continue
		}

		for _, emote := range emotes {
			result = append(
				result, emotes_cacher.Emote{
					ID:      string(emote.ID),
					Name:    emote.Name,
					Service: serviceName,
				},
			)
		}
	}

	return emotes_cacher.Response{
		Emotes: result,
	}, nil
}
