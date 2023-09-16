package timers

import (
	"context"
	"time"

	"github.com/samber/lo"
	"github.com/satont/twir/apps/scheduler/internal/types"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/generated/emotes_cacher"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

func NewEmotes(ctx context.Context, services *types.Services) {
	timeTick := lo.If(services.Config.AppEnv != "production", 15*time.Second).Else(5 * time.Minute)

	channelsTicker := time.NewTicker(timeTick)
	globalTicker := time.NewTicker(timeTick)

	go func() {
		for {
			select {
			case <-ctx.Done():
				channelsTicker.Stop()
				return
			case <-channelsTicker.C:
				var channels []model.Channels
				err := services.Gorm.
					Where(`"isEnabled" = ? and "isBanned" = ?`, true, false).
					Select("isEnabled", "id").
					Find(&channels).
					Error
				if err != nil {
					zap.S().Error(err)
				} else {
					for _, channel := range channels {
						_, err = services.Grpc.Emotes.CacheChannelEmotes(
							ctx,
							&emotes_cacher.Request{
								ChannelId: channel.ID,
							},
						)
						if err != nil {
							zap.S().Error(err)
						}
					}
				}
			}
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				globalTicker.Stop()
				break
			case <-globalTicker.C:
				_, err := services.Grpc.Emotes.CacheGlobalEmotes(ctx, &emptypb.Empty{})
				if err != nil {
					zap.S().Error(err)
				}
			}
		}
	}()
}
