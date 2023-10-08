package messages_updater

import (
	"context"
	"time"

	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/discord/internal/discord_go"
	"github.com/satont/twir/apps/discord/internal/sended_messages_store"
	cfg "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/grpc/generated/tokens"
	"github.com/satont/twir/libs/logger"
	"github.com/satont/twir/libs/twitch"
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

	TokensGrpc tokens.TokensClient
}

func New(opts Opts) (*MessagesUpdater, error) {
	twitchClient, err := twitch.NewAppClient(opts.Config, opts.TokensGrpc)
	if err != nil {
		return nil, err
	}

	updater := &MessagesUpdater{
		store:        opts.Store,
		logger:       opts.Logger.WithComponent("messages_updater"),
		config:       opts.Config,
		db:           opts.DB,
		discord:      opts.Discord,
		tokensGrpc:   opts.TokensGrpc,
		stopChan:     make(chan struct{}),
		twitchClient: twitchClient,
	}

	opts.LC.Append(
		fx.Hook{
			OnStart: func(_ context.Context) error {
				go updater.poll()
				updater.logger.Info("Messages updater is running")

				return nil
			},
			OnStop: func(_ context.Context) error {
				updater.stopChan <- struct{}{}
				close(updater.stopChan)
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

	tokensGrpc   tokens.TokensClient
	twitchClient *helix.Client

	stopChan chan struct{}
}

func (c *MessagesUpdater) poll() {
	ticker := time.NewTicker(
		lo.If(
			c.config.AppEnv != "production",
			10*time.Second,
		).Else(5 * time.Minute),
	)

	_, cancel := context.WithCancel(context.Background())

	for {
		select {
		case <-c.stopChan:
			cancel()
			break
		case <-ticker.C:
			// c.process(ctx)
		}
	}
}
