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
	"github.com/twirapp/twir/libs/logger"
	channelsintegrationsdiscord "github.com/twirapp/twir/libs/repositories/channels_integrations_discord"
	streamsrepository "github.com/twirapp/twir/libs/repositories/streams"
	"github.com/twirapp/twir/libs/twitch"
)

type Opts struct {
	Store       *sended_messages_store.SendedMessagesStore
	Logger      *slog.Logger
	LC          *lifecycle.Lifecycle
	Config      cfg.Config
	Discord     *discord_go.Discord
	Bus         *buscore.Bus
	DiscordRepo channelsintegrationsdiscord.Repository
	StreamsRepo streamsrepository.Repository
}

func New(opts Opts) (*MessagesUpdater, error) {
	twitchClient, err := twitch.NewAppClient(opts.Config, opts.Bus)
	if err != nil {
		return nil, err
	}

	updater := &MessagesUpdater{
		store:        opts.Store,
		logger:       logger.WithComponent(opts.Logger, "messages_updater"),
		config:       opts.Config,
		discord:      opts.Discord,
		twitchClient: twitchClient,
		twirBus:      opts.Bus,
		discordRepo:  opts.DiscordRepo,
		streamsRepo:  opts.StreamsRepo,
	}

	closeCtx, closeFunc := context.WithCancel(context.Background())

	// Start periodic updater in background
	opts.LC.Append(
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
