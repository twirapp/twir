package timers

import (
	"context"
	"log/slog"
	"time"

	"github.com/twirapp/twir/apps/scheduler/internal/services"
	"github.com/twirapp/twir/libs/baseapp/lifecycle"
	config "github.com/twirapp/twir/libs/config"
	twirlogger "github.com/twirapp/twir/libs/logger"
	"gorm.io/gorm"
)

func NewCommandsAndRoles(
	lc *lifecycle.Lifecycle,
	logger *slog.Logger,
	cfg config.Config,
	rolesService *services.Roles,
	commandsService *services.Commands,
	_ *gorm.DB,
) {
	timeTick := 15 * time.Second
	if cfg.AppEnv == "production" {
		timeTick = 5 * time.Minute
	}
	ticker := time.NewTicker(timeTick)

	ctx, cancel := context.WithCancel(context.Background())

	lc.Append(
		lifecycle.Hook{
			OnStart: func(_ context.Context) error {
				go func() {
					for {
						select {
						case <-ctx.Done():
							ticker.Stop()
							return
						case <-ticker.C:
							if err := rolesService.CreateDefaultRoles(ctx); err != nil {
								logger.Error("error while creating default roles", twirlogger.Error(err))
								return
							}

							if err := commandsService.CreateDefaultCommands(ctx); err != nil {
								logger.Error("error while creating default commands", twirlogger.Error(err))
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
