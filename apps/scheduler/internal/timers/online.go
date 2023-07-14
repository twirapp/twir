package timers

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	loParallel "github.com/samber/lo/parallel"
	"github.com/satont/twir/apps/scheduler/internal/types"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/twitch"
	"go.uber.org/zap"
)

func NewOnlineUsers(ctx context.Context, services *types.Services) {
	timeTick := lo.If(services.Config.AppEnv != "production", 15*time.Second).Else(5 * time.Minute)
	ticker := time.NewTicker(timeTick)

	go func() {
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case t := <-ticker.C:
				zap.S().Debugf("Online users timer tick at %s", t)

				var streams []*model.ChannelsStreams
				err := services.Gorm.Preload("Channel").Find(&streams).Error
				if err != nil {
					zap.S().Error(err)
					return
				}

				loParallel.ForEach(streams, func(stream *model.ChannelsStreams, _ int) {
					if stream.Channel != nil && (!stream.Channel.IsEnabled || stream.Channel.IsBanned) {
						return
					}

					twitchClient, err := twitch.NewUserClient(ctx, stream.UserId, *services.Config, services.Grpc.Tokens)
					if err != nil {
						zap.S().Error(err)
						return
					}

					broadcasterId := stream.UserId
					var chatters []helix.ChatChatter
					cursor := ""
					for {
						req, err := twitchClient.GetChannelChatChatters(&helix.GetChatChattersParams{
							BroadcasterID: broadcasterId,
							ModeratorID:   broadcasterId,
							After:         cursor, // TODO: cursor is not changing, should it?
						})
						if err != nil {
							zap.S().Error(err)
						} else {
							chatters = append(chatters, req.Data.Chatters...)
							if req.Data.Pagination.Cursor == "" {
								break
							}
						}
					}

					err = services.Gorm.Where(`"channelId" = ?`, broadcasterId).Delete(&model.UsersOnline{}).Error
					if err != nil {
						zap.S().Error(err)
						return
					}

					chattersChunks := lo.Chunk(chatters, 1000)
					loParallel.ForEach(chattersChunks, func(chunk []helix.ChatChatter, _ int) {
						dbChatters := make([]*model.Users, 0, len(chunk))
						err = services.Gorm.
							Where(
								`"id" IN ?`,
								lo.Map(
									chunk, func(chatter helix.ChatChatter, _ int) string {
										return chatter.UserID
									},
								),
							).
							Find(&dbChatters).
							Error
						if err != nil {
							zap.S().Error(err)
							return
						}

						usersForCreate := make([]*model.Users, 0, len(chunk))
						usersOnlineForCreate := make([]*model.UsersOnline, len(chunk))
						for i, chatter := range chunk {
							isExists := lo.SomeBy(
								dbChatters, func(item *model.Users) bool {
									return item.ID == chatter.UserID
								},
							)
							if !isExists {
								usersForCreate = append(usersForCreate, &model.Users{
									ID:     chatter.UserID,
									ApiKey: uuid.New().String(),
									Stats: &model.UsersStats{
										ID:        uuid.New().String(),
										ChannelID: broadcasterId,
										UserID:    chatter.UserID,
									},
								})
							}

							usersOnlineForCreate[i] = &model.UsersOnline{
								ID:        uuid.New().String(),
								ChannelId: broadcasterId,
								UserId:    null.StringFrom(chatter.UserID),
								UserName:  null.StringFrom(chatter.UserLogin),
							}
						}

						err = services.Gorm.Transaction(func(tx *gorm.DB) error {
							if err := tx.CreateInBatches(usersForCreate, 1000).Error; err != nil {
								return err
							}

							if err := tx.CreateInBatches(usersOnlineForCreate, 1000).Error; err != nil {
								return err
							}

							return nil
						})
						if err != nil {
							zap.S().Error(err)
						}
					})
				})
			}
		}
	}()
}
