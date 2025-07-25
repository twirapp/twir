package discord_go

import (
	"context"

	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	cfg "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/logger"
	"go.uber.org/fx"
	"gorm.io/gorm"

	"github.com/diamondburned/arikawa/v3/session/shard"
)

type Opts struct {
	fx.In

	LC     fx.Lifecycle
	Config cfg.Config
	Logger logger.Logger
	Db     *gorm.DB
}

type Discord struct {
	*shard.Manager

	logger logger.Logger
	db     *gorm.DB
}

func New(opts Opts) (*Discord, error) {
	if opts.Config.DiscordBotToken == "" {
		return &Discord{}, nil
	}

	log := opts.Logger.WithComponent("discord_session")
	discord := &Discord{
		logger: log,
		db:     opts.Db,
	}

	newShard := state.NewShardFunc(
		func(m *shard.Manager, s *state.State) {
			// Add the needed Gateway intents.
			s.AddIntents(gateway.IntentGuilds)
			s.AddIntents(gateway.IntentGuildMessages)
			s.AddIntents(gateway.IntentDirectMessages)
			s.AddIntents(gateway.IntentGuildMembers)

			s.AddHandler(
				func(c *gateway.ReadyEvent) {
					discord.handleShardReady(c)
				},
			)

			s.AddHandler(
				func(c *gateway.GuildDeleteEvent) {
					discord.handleGuildDelete(c)
				},
			)
		},
	)

	shardManager, err := shard.NewManager("Bot "+opts.Config.DiscordBotToken, newShard)
	if err != nil {
		return nil, err
	}

	opts.LC.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				if err := shardManager.Open(context.Background()); err != nil {
					return err
				}

				return nil
			},
			OnStop: func(ctx context.Context) error {
				return shardManager.Close()
			},
		},
	)

	discord.Manager = shardManager

	return discord, nil
}
