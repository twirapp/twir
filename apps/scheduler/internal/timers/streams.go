package timers

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/lib/pq"
	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/scheduler/internal/types"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/twitch"
	"go.uber.org/zap"
)

func NewStreams(ctx context.Context, services *types.Services) {
	timeTick := lo.If(services.Config.AppEnv != "production", 15*time.Second).Else(5 * time.Minute)
	ticker := time.NewTicker(timeTick)

	go func() {
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			case <-ticker.C:
				processStreams(services)
			}
		}
	}()
}

func processStreams(services *types.Services) {
	var channels []model.Channels
	err := services.Gorm.
		Where(`"isEnabled" = ? and "isBanned" = ?`, true, false).
		Select("id", `"isEnabled"`, `"isBanned"`).
		Find(&channels).Error
	if err != nil {
		zap.S().Error(err)
		return
	}

	var existedStreams []model.ChannelsStreams
	err = services.Gorm.Select("id", `"userId"`, `"parsedMesages"`).Find(&existedStreams).Error
	if err != nil {
		zap.S().Error(err)
		return
	}

	twitchClient, err := twitch.NewAppClient(*services.Config, services.Grpc.Tokens)
	if err != nil {
		zap.S().Error(err)
		return
	}

	chunks := lo.Chunk(channels, 100)
	wg := &sync.WaitGroup{}

	wg.Add(len(chunks))

	for _, chunk := range chunks {
		go func(chunk []model.Channels) {
			defer wg.Done()
			usersIds := lo.Map(
				chunk, func(channel model.Channels, _ int) string {
					return channel.ID
				},
			)

			streams, err := twitchClient.GetStreams(
				&helix.StreamsParams{
					UserIDs: usersIds,
				},
			)

			if err != nil || streams.ErrorMessage != "" {
				zap.S().Error(err)
				return
			}

			for _, channel := range chunk {
				twitchStream, twitchStreamExists := lo.Find(
					streams.Data.Streams, func(stream helix.Stream) bool {
						return stream.UserID == channel.ID
					},
				)
				dbStream, dbStreamExists := lo.Find(
					existedStreams, func(stream model.ChannelsStreams) bool {
						return stream.UserId == channel.ID
					},
				)

				tags := &pq.StringArray{}
				for _, tag := range twitchStream.Tags {
					*tags = append(*tags, tag)
				}

				channelStream := &model.ChannelsStreams{
					ID:             twitchStream.ID,
					UserId:         twitchStream.UserID,
					UserLogin:      twitchStream.UserLogin,
					UserName:       twitchStream.UserName,
					GameId:         twitchStream.GameID,
					GameName:       twitchStream.GameName,
					CommunityIds:   nil,
					Type:           twitchStream.Type,
					Title:          twitchStream.Title,
					ViewerCount:    twitchStream.ViewerCount,
					StartedAt:      twitchStream.StartedAt,
					Language:       twitchStream.Language,
					ThumbnailUrl:   twitchStream.ThumbnailURL,
					TagIds:         nil,
					Tags:           tags,
					IsMature:       twitchStream.IsMature,
					ParsedMessages: dbStream.ParsedMessages,
				}

				if twitchStreamExists && dbStreamExists {
					err = services.Gorm.Where(`"userId" = ?`, channel.ID).Save(channelStream).Error
					if err != nil {
						zap.S().Error(err)
						return
					}
				}

				if twitchStreamExists && !dbStreamExists {
					bytes, err := json.Marshal(
						&streamOnlineMessage{
							StreamID:  channelStream.ID,
							ChannelID: channelStream.UserId,
						},
					)
					if err != nil {
						zap.S().Error(err)
						return
					}

					services.PubSub.Publish("stream.online", bytes)
				}

				if !twitchStreamExists && dbStreamExists {
					err = services.Gorm.Where(`"userId" = ?`, channel.ID).Delete(&model.ChannelsStreams{}).Error
					if err != nil {
						zap.S().Error(err)
						return
					}

					bytes, err := json.Marshal(
						&streamOfflineMessage{
							ChannelID: channelStream.UserId,
						},
					)
					if err != nil {
						zap.S().Error(err)
						return
					}

					services.PubSub.Publish("stream.offline", bytes)
				}
			}
		}(chunk)
	}

	wg.Wait()
}

// { streamId: stream.id, channelId: channel }
type streamOnlineMessage struct {
	StreamID  string `json:"streamId"`
	ChannelID string `json:"channelId"`
}

type streamOfflineMessage struct {
	ChannelID string `json:"channelId"`
}
