package timers

import (
	"context"
	"time"

	"github.com/samber/lo"
	"github.com/satont/twir/apps/scheduler/internal/types"
	model "github.com/satont/twir/libs/gomodels"
	"go.uber.org/zap"
)

func NewCommandsAndRoles(ctx context.Context, services *types.Services) {
	timeTick := lo.If(services.Config.AppEnv != "production", 15*time.Second).Else(5 * time.Minute)
	ticker := time.NewTicker(timeTick)

	go func() {
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			case <-ticker.C:
				var channels []model.Channels
				if err := services.Gorm.Select(`"id"`).Find(&channels).Error; err != nil {
					zap.S().Error(err)
					return
				}

				channelIds := lo.Map(
					channels, func(channel model.Channels, _ int) string {
						return channel.ID
					},
				)

				if err := services.Roles.CreateDefaultRoles(ctx, channelIds); err != nil {
					zap.S().Error(err)
					return
				}

				if err := services.Commands.CreateDefaultCommands(ctx, channelIds); err != nil {
					zap.S().Error(err)
					return
				}
			}
		}
	}()
}
