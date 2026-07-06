package song_requests

import (
	"context"
	"log/slog"

	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/api"
	"github.com/twirapp/twir/libs/wsrouter"
	"go.uber.org/fx"
)

type BridgeOpts struct {
	fx.In
	LC                   fx.Lifecycle
	WsRouter             wsrouter.WsRouter
	TwirBus              *buscore.Bus
	Logger               *slog.Logger
	PlaybackStateService *PlaybackStateService
}

type Bridge struct {
	wsRouter             wsrouter.WsRouter
	twirBus              *buscore.Bus
	logger               *slog.Logger
	playbackStateService *PlaybackStateService
}

func NewBridge(opts BridgeOpts) *Bridge {
	b := &Bridge{
		wsRouter:             opts.WsRouter,
		twirBus:              opts.TwirBus,
		logger:               opts.Logger,
		playbackStateService: opts.PlaybackStateService,
	}

	opts.LC.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			b.twirBus.Api.SongRequestAddToQueue.SubscribeGroup("api",
				func(ctx context.Context, data api.SongRequestAddToQueue) (struct{}, error) {
					return struct{}{}, b.wsRouter.Publish(
						"api.songRequestQueue."+data.ChannelID, data,
					)
				},
			)
			b.logger.Info("Subscribed to SongRequestAddToQueue events")

			b.twirBus.Api.SongRequestRemoveFromQueue.SubscribeGroup("api",
				func(ctx context.Context, data api.SongRequestRemoveFromQueue) (struct{}, error) {
					return struct{}{}, b.wsRouter.Publish(
						"api.songRequestQueueRemove."+data.ChannelID, data,
					)
				},
			)
			b.logger.Info("Subscribed to SongRequestRemoveFromQueue events")

			b.twirBus.Api.SongRequestPlaybackState.SubscribeGroup("api",
				func(ctx context.Context, data api.SongRequestPlaybackState) (struct{}, error) {
					return struct{}{}, b.wsRouter.Publish(
						"api.songRequestPlayback."+data.ChannelID, data,
					)
				},
			)
			b.logger.Info("Subscribed to SongRequestPlaybackState events")

			b.playbackStateService.StartTicker(ctx, b.wsRouter)
			b.logger.Info("Started playback state ticker")

			return nil
		},
		OnStop: func(ctx context.Context) error {
			b.twirBus.Api.SongRequestAddToQueue.Unsubscribe()
			b.twirBus.Api.SongRequestRemoveFromQueue.Unsubscribe()
			b.twirBus.Api.SongRequestPlaybackState.Unsubscribe()
			return nil
		},
	})

	return b
}
