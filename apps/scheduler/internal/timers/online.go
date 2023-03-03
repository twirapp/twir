package timers

import (
	"context"
	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/samber/lo"
	"github.com/satont/go-helix/v2"
	"github.com/satont/tsuwari/apps/scheduler/internal/types"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/twitch"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"sync"
	"time"
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
				var streams []model.ChannelsStreams
				err := services.Gorm.Find(&streams).Error
				if err != nil {
					zap.S().Error(err)
				} else {
					wg := &sync.WaitGroup{}
					for _, stream := range streams {
						wg.Add(1)
						twitchClient, err := twitch.NewUserClient(stream.UserId, *services.Config, services.Grpc.Tokens)
						if err != nil {
							zap.S().Error(err)
						} else {
							go func(userId string) {
								defer wg.Done()
								var chatters []helix.ChatChatter
								cursor := ""
								for {
									req, err := twitchClient.GetChannelChatChatters(&helix.GetChatChattersParams{
										BroadcasterID: stream.UserId,
										ModeratorID:   stream.UserId,
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
									for _, chatter := range chatters {
										var user model.Users
										err := tx.Where("id = ?", chatter.UserID).First(&user).Error
										if err != nil {
											if err == gorm.ErrRecordNotFound {
												user = model.Users{
													ID: chatter.UserID,
												}
												err = tx.Create(&user).Error
												if err != nil {
													return err
												}
											} else {
												return err
											}
										} else {
											tx.Delete(&model.UsersOnline{}).Where(`"userId" = ?`, chatter.UserID)
											tx.Save(&model.UsersOnline{
												ID:        uuid.New().String(),
												ChannelId: userId,
												UserId:    null.StringFrom(chatter.UserID),
												UserName:  null.StringFrom(chatter.UserLogin),
											})
										}
									}
									return nil
								})

								if err != nil {
									zap.S().Error(err)
								}
							}(stream.UserId)
						}
					}
					wg.Wait()
				}

			}
		}
	}()
}
