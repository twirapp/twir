package timers

import (
	"context"
	"time"

	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/scheduler/internal/types"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/twitch"
	"github.com/satont/twir/libs/utils"
	"go.uber.org/zap"
)

type BannedChannels struct {
	services *types.Services
}

func NewBannerChannels(ctx context.Context, services *types.Services) *BannedChannels {
	timeTick := lo.If(services.Config.AppEnv != "production", 15*time.Second).Else(5 * time.Minute)
	ticker := time.NewTicker(timeTick)

	bannedChannels := &BannedChannels{
		services: services,
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			case <-ticker.C:
				bannedChannels.process(ctx)
			}
		}
	}()

	bannedChannels.process(ctx)

	return bannedChannels
}

func (c *BannedChannels) process(ctx context.Context) {
	var channels []model.Channels
	if err := c.services.Gorm.WithContext(ctx).Find(&channels).Error; err != nil {
		zap.S().Error("failed to get channels", zap.Error(err))
		return
	}

	chunks := lo.Chunk(channels, 100)

	wg := utils.NewGoroutinesGroup()

	for _, chunk := range chunks {
		chunk := chunk

		wg.Go(
			func() {
				twitchClient, err := twitch.NewAppClientWithContext(
					ctx, *c.services.Config,
					c.services.Grpc.Tokens,
				)
				if err != nil {
					zap.S().Error("failed to create twitch client", zap.Error(err))
					return
				}

				channelsReq, err := twitchClient.GetUsers(
					&helix.UsersParams{
						IDs: lo.Map(chunk, func(channel model.Channels, _ int) string { return channel.ID }),
					},
				)
				if err != nil {
					zap.S().Error("failed to get users", zap.Error(err))
					return
				}
				if channelsReq.ErrorMessage != "" {
					zap.S().Error("failed to get users", zap.String("error", channelsReq.ErrorMessage))
					return
				}

				for _, channel := range chunk {
					exists := lo.SomeBy(
						channelsReq.Data.Users, func(user helix.User) bool {
							return user.ID == channel.ID
						},
					)

					err := c.services.Gorm.WithContext(ctx).Model(&channel).Update(
						`"isTwitchBanned"`,
						!exists,
					).Error
					if err != nil {
						zap.S().Error("failed to update channel", zap.Error(err))
						continue
					}
				}
			},
		)
	}

	wg.Wait()
}
