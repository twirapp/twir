package song_requests

import (
	"context"
	"github.com/twirapp/twir/libs/baseapp/lifecycle"
	"log/slog"

	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/api"
	"github.com/twirapp/twir/libs/wsrouter"
)

type BridgeOpts struct {
	LC                   *lifecycle.Lifecycle
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

	opts.LC.Append(lifecycle.Hook{
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

			return nil
		},
		OnStop: func(ctx context.Context) error {
			b.twirBus.Api.SongRequestAddToQueue.Unsubscribe()
			b.twirBus.Api.SongRequestRemoveFromQueue.Unsubscribe()
			return nil
		},
	})

	return b
}
