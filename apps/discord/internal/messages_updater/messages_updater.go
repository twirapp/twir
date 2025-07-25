package messages_updater

import (
	"context"
	"log/slog"

	"github.com/nicklaw5/helix/v2"
	"github.com/twirapp/twir/apps/discord/internal/discord_go"
	"github.com/twirapp/twir/apps/discord/internal/sended_messages_store"
	cfg "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/logger"
	"github.com/twirapp/twir/libs/twitch"
	buscore "github.com/twirapp/twir/libs/bus-core"
	bustwitch "github.com/twirapp/twir/libs/bus-core/twitch"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In

	Store   *sended_messages_store.SendedMessagesStore
	Logger  logger.Logger
	LC      fx.Lifecycle
	Config  cfg.Config
	DB      *gorm.DB
	Discord *discord_go.Discord
	Bus     *buscore.Bus
}

func New(opts Opts) (*MessagesUpdater, error) {
	twitchClient, err := twitch.NewAppClient(opts.Config, opts.Bus)
	if err != nil {
		return nil, err
	}

	updater := &MessagesUpdater{
		store:        opts.Store,
		logger:       opts.Logger.WithComponent("messages_updater"),
		config:       opts.Config,
		db:           opts.DB,
		discord:      opts.Discord,
		twitchClient: twitchClient,
		twirBus:      opts.Bus,
	}

	opts.LC.Append(
		fx.Hook{
			OnStart: func(_ context.Context) error {
				opts.Bus.Channel.StreamOnline.SubscribeGroup(
					"discord",
					func(ctx context.Context, data bustwitch.StreamOnlineMessage) (struct{}, error) {
						if err := updater.processOnline(ctx, data.ChannelID); err != nil {
							opts.Logger.Error("Failed to process online", slog.Any("err", err))
							return struct{}{}, err
						}

						return struct{}{}, nil
					},
				)

				opts.Bus.Channel.StreamOffline.SubscribeGroup(
					"discord",
					func(ctx context.Context, data bustwitch.StreamOfflineMessage) (struct{}, error) {
						if err := updater.processOffline(ctx, data.ChannelID); err != nil {
							opts.Logger.Error("Failed to process offline", slog.Any("err", err))
							return struct{}{}, err
						}

						return struct{}{}, nil
					},
				)

				return nil
			},
			OnStop: func(_ context.Context) error {
				opts.Bus.Channel.StreamOnline.Unsubscribe()
				opts.Bus.Channel.StreamOffline.Unsubscribe()

				return nil
			},
		},
	)

	return updater, nil
}

type MessagesUpdater struct {
	store   *sended_messages_store.SendedMessagesStore
	logger  logger.Logger
	config  cfg.Config
	db      *gorm.DB
	discord *discord_go.Discord

	twirBus      *buscore.Bus
	twitchClient *helix.Client
}
