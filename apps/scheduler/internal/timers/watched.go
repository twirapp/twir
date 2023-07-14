package timers

import (
	"context"
	"time"

	"github.com/samber/lo"
	"go.uber.org/zap"

	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/generated/watched"

	"github.com/satont/twir/apps/scheduler/internal/types"
)

func NewWatched(ctx context.Context, services *types.Services) {
	timeTick := lo.If(services.Config.AppEnv != "production", 15*time.Second).Else(5 * time.Minute)
	ticker := time.NewTicker(timeTick)

	go func() {
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case t := <-ticker.C:
				zap.S().Debugf("watched timer ticked at %s", t)

				var streams []model.ChannelsStreams
				err := services.Gorm.
					Preload("Channel").
					Find(&streams).
					Error
				if err != nil {
					zap.S().Error(err)
					continue
				}

				groups := lo.GroupBy(streams, func(item model.ChannelsStreams) string {
					return item.Channel.BotID
				})

				for botId, streams := range groups {
					chunks := lo.Chunk(streams, 100)

					for _, chunk := range chunks {
						channelsIds := lo.Map(chunk, func(stream model.ChannelsStreams, _ int) string {
							return stream.Channel.ID
						})

						_, err := services.Grpc.Watched.IncrementByChannelId(ctx, &watched.Request{
							BotId:      botId,
							ChannelsId: channelsIds,
						})
						if err != nil {
							zap.S().Error(err)
						}
					}
				}
			}
		}
	}()
}
