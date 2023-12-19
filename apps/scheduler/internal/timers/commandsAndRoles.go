package timers

import (
	"context"
	"log/slog"
	"time"

	"github.com/samber/lo"
	"github.com/satont/twir/apps/scheduler/internal/services"
	config "github.com/satont/twir/libs/config"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/logger"
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
	timeTick := lo.If(opts.Config.AppEnv != "production", 15*time.Second).Else(5 * time.Minute)
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
							var channels []model.Channels
							if err := opts.Gorm.WithContext(ctx).Select(`"id"`).Find(&channels).Error; err != nil {
								opts.Logger.Error("error while getting channels", slog.Any("err", err))
								return
							}

							channelIds := lo.Map(
								channels, func(channel model.Channels, _ int) string {
									return channel.ID
								},
							)

							if err := opts.RolesService.CreateDefaultRoles(ctx, channelIds); err != nil {
								opts.Logger.Error("error while creating default roles", slog.Any("err", err))
								return
							}

							if err := opts.CommandsService.CreateDefaultCommands(ctx, channelIds); err != nil {
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
