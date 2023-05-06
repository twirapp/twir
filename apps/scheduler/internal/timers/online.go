package timers

import (
	"context"
	"gorm.io/gorm"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	loParallel "github.com/samber/lo/parallel"
	"github.com/satont/tsuwari/apps/scheduler/internal/types"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/twitch"
	"go.uber.org/zap"
)

func NewOnlineUsers(ctx context.Context, services *types.Services) {
	timeTick := lo.If(services.Config.AppEnv != "production", 15*time.Second).Else(5 * time.Minute)
	ticker := time.NewTicker(timeTick)

	go func() {
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				break
			case <-ticker.C:
				var streams []*model.ChannelsStreams
				err := services.Gorm.Find(&streams).Error
				if err != nil {
					zap.S().Error(err)
					return
				}

				streamsWg := &sync.WaitGroup{}
				for _, stream := range streams {
					streamsWg.Add(1)

					twitchClient, err := twitch.NewUserClient(stream.UserId, *services.Config, services.Grpc.Tokens)
					if err != nil {
						zap.S().Error(err)
						streamsWg.Done()
						continue
					}

					go func(broadcasterId string) {
						defer streamsWg.Done()
						var chatters []helix.ChatChatter
						cursor := ""
						for {
							req, err := twitchClient.GetChannelChatChatters(&helix.GetChatChattersParams{
								BroadcasterID: broadcasterId,
								ModeratorID:   broadcasterId,
								After:         cursor,
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

						chattersChunks := lo.Chunk(chatters, 1000)

						err = services.Gorm.Where(`"channelId" = ?`, broadcasterId).Delete(&model.UsersOnline{}).Error
						if err != nil {
							zap.S().Error(err)
							return
						}

						loParallel.ForEach(chattersChunks, func(chunk []helix.ChatChatter, _ int) {
							dbChatters := make([]*model.Users, 0, len(chunk))
							err = services.Gorm.
								Where(
									`"id" IN ?`,
									lo.Map(chunk, func(chatter helix.ChatChatter, _ int) string {
										return chatter.UserID
									}),
								).
								Find(&dbChatters).
								Error

							if err != nil {
								zap.S().Error(err)
								return
							}

							var usersForCreate []*model.Users
							var usersOnlineForCreate []*model.UsersOnline

							for _, chatter := range chunk {
								isExists := lo.SomeBy(dbChatters, func(item *model.Users) bool {
									return item.ID == chatter.UserID
								})
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

								usersOnlineForCreate = append(usersOnlineForCreate, &model.UsersOnline{
									ID:        uuid.New().String(),
									ChannelId: broadcasterId,
									UserId:    null.StringFrom(chatter.UserID),
									UserName:  null.StringFrom(chatter.UserLogin),
								})
							}

							err = services.Gorm.Transaction(func(tx *gorm.DB) error {
								err = services.Gorm.CreateInBatches(usersForCreate, 1000).Error
								if err != nil {
									return err
								}

								err = services.Gorm.CreateInBatches(usersOnlineForCreate, 1000).Error
								if err != nil {
									return err
								}

								return nil
							})
							if err != nil {
								zap.S().Error(err)
							}
						})
					}(stream.UserId)
				}
				streamsWg.Wait()
			}
		}
	}()
}
