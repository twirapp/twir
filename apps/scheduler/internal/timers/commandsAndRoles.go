package timers

import (
	"context"
	"log/slog"
	"time"

	"github.com/twirapp/twir/apps/scheduler/internal/services"
	config "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/logger"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type CommandsAndRolesOpts struct {
	fx.In
	Lc fx.Lifecycle

	Logger logger.Logger
	Config config.Config

	RolesService    *services.Roles
	CommandsService *services.Commands
	Gorm            *gorm.DB
}

func NewCommandsAndRoles(opts CommandsAndRolesOpts) {
	timeTick := 15 * time.Second
	if opts.Config.AppEnv == "production" {
		timeTick = 5 * time.Minute
	}
	ticker := time.NewTicker(timeTick)

	ctx, cancel := context.WithCancel(context.Background())

	opts.Lc.Append(
		fx.Hook{
			OnStart: func(_ context.Context) error {
				go func() {
					for {
						select {
						case <-ctx.Done():
							ticker.Stop()
							return
						case <-ticker.C:
							if err := opts.RolesService.CreateDefaultRoles(ctx); err != nil {
								opts.Logger.Error("error while creating default roles", slog.Any("err", err))
								return
							}

							if err := opts.CommandsService.CreateDefaultCommands(ctx); err != nil {
								opts.Logger.Error("error while creating default commands", slog.Any("err", err))
								return
							}
						}
					}
				}()

				return nil
			},
			OnStop: func(_ context.Context) error {
				cancel()
				return nil
			},
		},
	)
}
