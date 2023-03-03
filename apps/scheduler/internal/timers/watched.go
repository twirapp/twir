package timers

import (
	"context"
	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/scheduler/internal/types"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/grpc/generated/watched"
	"go.uber.org/zap"
	"time"
)

func NewWatched(ctx context.Context, services *types.Services) {
	timeTick := lo.If(services.Config.AppEnv != "production", 15*time.Second).Else(5 * time.Minute)
	ticker := time.NewTicker(timeTick)

	go func() {
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				break
			case <-ticker.C:
				var streams []model.ChannelsStreams
				err := services.Gorm.
					Preload("Channel").
					Find(&streams).
					Error
				if err != nil {
					zap.S().Error(err)
				} else {
					groups := lo.GroupBy(streams, func(item model.ChannelsStreams) string {
						return item.Channel.BotID
					})

					for botId, streams := range groups {
						chunks := lo.Chunk(streams, 100)

						for _, chunk := range chunks {
							var channelsIds []string
							for _, stream := range chunk {
								channelsIds = append(channelsIds, stream.Channel.ID)
							}

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
		}
	}()
}
