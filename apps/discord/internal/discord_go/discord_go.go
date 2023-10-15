package discord_go

import (
	"context"
	"errors"

	cfg "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/logger"
	"github.com/switchupcb/disgo"
	"github.com/switchupcb/disgo/shard"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In

	LC     fx.Lifecycle
	Config cfg.Config
	Logger logger.Logger
	Db     *gorm.DB
}

type Discord struct {
	*disgo.Client

	logger logger.Logger
	db     *gorm.DB
}

func New(opts Opts) (*Discord, error) {
	if opts.Config.DiscordBotToken == "" {
		return nil, errors.New("discord bot token is empty")
	}

	log := opts.Logger.WithComponent("discord_session")
	discord := &Discord{
		logger: log,
		db:     opts.Db,
	}

	bot := &disgo.Client{
		ApplicationID:  opts.Config.DiscordClientID,
		Authentication: disgo.BotToken(opts.Config.DiscordBotToken),
		Authorization:  &disgo.Authorization{},
		Config:         disgo.DefaultConfig(),
		Handlers: &disgo.Handlers{
			Ready:       []func(*disgo.Ready){discord.handleReady},
			GuildDelete: []func(*disgo.GuildDelete){discord.handleGuildDelete},
		},
		Sessions: disgo.NewSessionManager(),
	}

	bot.Config.Gateway.ShardManager = new(shard.InstanceShardManager)
	// auto sharding
	bot.Config.Gateway.ShardManager.SetNumShards(0)
	s := bot.Config.Gateway.ShardManager

	opts.LC.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {

				if err := s.Connect(bot); err != nil {
					return err
				}

				return nil
			},
			OnStop: func(ctx context.Context) error {
				return s.Disconnect()
			},
		},
	)

	discord.Client = bot

	return discord, nil
}
