package discordmessagesupdater

import (
	"context"
	"log/slog"

	"github.com/nicklaw5/helix/v2"
	"github.com/twirapp/twir/apps/bots/internal/discord/discord_go"
	"github.com/twirapp/twir/apps/bots/internal/discord/sended_messages_store"
	"github.com/twirapp/twir/libs/baseapp/lifecycle"
	buscore "github.com/twirapp/twir/libs/bus-core"
	cfg "github.com/twirapp/twir/libs/config"
	loggerlib "github.com/twirapp/twir/libs/logger"
	channelsintegrationsdiscord "github.com/twirapp/twir/libs/repositories/channels_integrations_discord"
	streamsrepository "github.com/twirapp/twir/libs/repositories/streams"
	"github.com/twirapp/twir/libs/twitch"
)

func New(
	store *sended_messages_store.SendedMessagesStore,
	logger *slog.Logger,
	lc *lifecycle.Lifecycle,
	cfg cfg.Config,
	discord *discord_go.Discord,
	bus *buscore.Bus,
	discordRepo channelsintegrationsdiscord.Repository,
	streamsRepo streamsrepository.Repository,
) (*MessagesUpdater, error) {
	twitchClient, err := twitch.NewAppClient(cfg, bus)
	if err != nil {
		return nil, err
	}

	updater := &MessagesUpdater{
		store:        store,
		logger:       loggerlib.WithComponent(logger, "messages_updater"),
		config:       cfg,
		discord:      discord,
		twitchClient: twitchClient,
		twirBus:      bus,
		discordRepo:  discordRepo,
		streamsRepo:  streamsRepo,
	}

	closeCtx, closeFunc := context.WithCancel(context.Background())

	// Start periodic updater in background
	lc.Append(
		lifecycle.Hook{
			OnStart: func(ctx context.Context) error {
				go updater.StartPeriodicUpdater(closeCtx)
				return nil
			},
			OnStop: func(ctx context.Context) error {
				closeFunc()
				return nil
			},
		},
	)

	return updater, nil
}

type MessagesUpdater struct {
	store   *sended_messages_store.SendedMessagesStore
	logger  *slog.Logger
	config  cfg.Config
	discord *discord_go.Discord

	twirBus      *buscore.Bus
	twitchClient *helix.Client
	discordRepo  channelsintegrationsdiscord.Repository
	streamsRepo  streamsrepository.Repository
}
