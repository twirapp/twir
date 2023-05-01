package timers

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/scheduler/internal/types"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/twitch"
	"go.uber.org/zap"
	"gorm.io/gorm"
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

						err = services.Gorm.Transaction(func(tx *gorm.DB) error {
							err = tx.Where(`"channelId" = ?`, broadcasterId).Delete(&model.UsersOnline{}).Error

							if err != nil {
								return err
							}

							for _, chatter := range chatters {
								user := &model.Users{}
								err = tx.Where("id = ?", chatter.UserID).Find(user).Error
								if err != nil {
									return err
								}

								if user.ID == "" {
									user = &model.Users{
										ID:     chatter.UserID,
										ApiKey: uuid.New().String(),
										Stats: &model.UsersStats{
											ID:        uuid.New().String(),
											ChannelID: broadcasterId,
											UserID:    chatter.UserID,
										},
									}
									err = tx.Create(user).Error
									if err != nil {
										return err
									}

									continue
								}

								err = tx.Save(&model.UsersOnline{
									ID:        uuid.New().String(),
									ChannelId: broadcasterId,
									UserId:    null.StringFrom(chatter.UserID),
									UserName:  null.StringFrom(chatter.UserLogin),
								}).Error

								if err != nil {
									return err
								}

								continue
							}
							return nil
						})

						if err != nil {
							zap.S().Error(err)
						}
					}(stream.UserId)
				}
				streamsWg.Wait()
			}
		}
	}()
}
