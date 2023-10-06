package discord_go

import (
	"context"
	"errors"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	cfg "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/logger"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	LC     fx.Lifecycle
	Config cfg.Config
	Logger logger.Logger
}

type Discord struct {
	*discordgo.Session

	logger logger.Logger
}

func New(opts Opts) (*Discord, error) {
	if opts.Config.DiscordBotToken == "" {
		return nil, errors.New("discord bot token is empty")
	}

	log := opts.Logger.WithComponent("discord_session")

	dgo, err := discordgo.New("Bot " + opts.Config.DiscordBotToken)
	if err != nil {
		return nil, err
	}

	dgo.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)

	opts.LC.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				if err := dgo.Open(); err != nil {
					return err
				}

				log.Info(
					"Discord bot is running",
					slog.Group(
						"bot",
						slog.String("id", dgo.State.User.ID),
						slog.String("name", dgo.State.User.Username),
						slog.Bool("verified", dgo.State.User.Verified),
					),
				)

				return nil
			},
			OnStop: func(ctx context.Context) error {
				return dgo.Close()
			},
		},
	)

	return &Discord{
		Session: dgo,
		logger:  log,
	}, nil
}
