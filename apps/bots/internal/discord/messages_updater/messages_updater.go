package discordmessagesupdater

import (
	"log/slog"

	"github.com/nicklaw5/helix/v2"
	"github.com/twirapp/twir/apps/bots/internal/discord/discord_go"
	"github.com/twirapp/twir/apps/bots/internal/discord/sended_messages_store"
	buscore "github.com/twirapp/twir/libs/bus-core"
	cfg "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/logger"
	channelsintegrationsdiscord "github.com/twirapp/twir/libs/repositories/channels_integrations_discord"
	"github.com/twirapp/twir/libs/twitch"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In

	Store       *sended_messages_store.SendedMessagesStore
	Logger      *slog.Logger
	LC          fx.Lifecycle
	Config      cfg.Config
	DB          *gorm.DB
	Discord     *discord_go.Discord
	Bus         *buscore.Bus
	DiscordRepo channelsintegrationsdiscord.Repository
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
		db:           opts.DB,
		discord:      opts.Discord,
		twitchClient: twitchClient,
		twirBus:      opts.Bus,
		discordRepo:  opts.DiscordRepo,
	}

	return updater, nil
}

type MessagesUpdater struct {
	store   *sended_messages_store.SendedMessagesStore
	logger  *slog.Logger
	config  cfg.Config
	db      *gorm.DB
	discord *discord_go.Discord

	twirBus      *buscore.Bus
	twitchClient *helix.Client
	discordRepo  channelsintegrationsdiscord.Repository
}
