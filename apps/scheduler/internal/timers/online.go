package timers

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/scheduler/internal/types"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/twitch"
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
				return
			case <-ticker.C:
				updateOnlineUsers(ctx, services)
			}
		}
	}()
}

func updateOnlineUsers(ctx context.Context, services *types.Services) {
	streams, err := getStreams(ctx, services)
	if err != nil {
		zap.S().Error(err)
		return
	}

	var wg sync.WaitGroup
	for _, stream := range streams {
		if shouldSkipStream(stream) {
			continue
		}
		wg.Add(1)
		go func(broadcasterID string) {
			defer wg.Done()
			updateStreamUsers(ctx, broadcasterID, services)
		}(stream.UserId)
	}
	wg.Wait()
}

func getStreams(ctx context.Context, services *types.Services) ([]*model.ChannelsStreams, error) {
	var streams []*model.ChannelsStreams
	err := services.Gorm.WithContext(ctx).Preload("Channel").Find(&streams).Error
	return streams, err
}

func shouldSkipStream(stream *model.ChannelsStreams) bool {
	return stream.Channel == nil || (!stream.Channel.IsEnabled || stream.Channel.IsBanned)
}

func updateStreamUsers(ctx context.Context, broadcasterID string, services *types.Services) {
	twitchClient, err := twitch.NewUserClientWithContext(
		ctx,
		broadcasterID,
		*services.Config,
		services.Grpc.Tokens,
	)
	if err != nil {
		zap.S().Error(err)
		return
	}

	var cursor string

	for {
		params := &helix.GetChatChattersParams{
			BroadcasterID: broadcasterID,
			ModeratorID:   broadcasterID,
			After:         cursor,
			First:         "1000",
		}
		req, err := twitchClient.GetChannelChatChatters(params)
		if err != nil {
			zap.S().Error(err)
			return
		}
		if req.ErrorMessage != "" {
			zap.S().Error(req.ErrorMessage)
			return
		}

		chatters := req.Data.Chatters

		if len(chatters) == 0 {
			return
		}

		usersIdsForRequest := lo.Map(
			chatters, func(chatter helix.ChatChatter, _ int) string {
				return chatter.UserID
			},
		)

		err = services.Gorm.WithContext(ctx).Transaction(
			func(tx *gorm.DB) error {
				var existedUsers []model.Users
				if err := tx.
					Select("id").
					Where("id IN ?", usersIdsForRequest).
					Find(&existedUsers).Error; err != nil {
					return err
				}
				var usersForCreate []model.Users
				for _, chatter := range chatters {
					_, chatterExists := lo.Find(
						existedUsers, func(user model.Users) bool {
							return user.ID == chatter.UserID
						},
					)
					if !chatterExists {
						usersForCreate = append(
							usersForCreate,
							model.Users{
								ID:     chatter.UserID,
								ApiKey: uuid.New().String(),
								Stats: &model.UsersStats{
									ID:        uuid.New().String(),
									UserID:    chatter.UserID,
									ChannelID: broadcasterID,
								},
							},
						)
					}
				}

				if err := tx.Where(
					`"channelId" = ?`,
					broadcasterID,
				).Delete(&model.UsersOnline{}).Error; err != nil {
					return err
				}

				onlineChattersForCreate := make([]model.UsersOnline, 0, len(chatters))
				for _, chatter := range chatters {
					onlineChattersForCreate = append(
						onlineChattersForCreate,
						model.UsersOnline{
							ID:        uuid.New().String(),
							ChannelId: broadcasterID,
							UserId:    null.StringFrom(chatter.UserID),
							UserName:  null.StringFrom(chatter.UserLogin),
						},
					)
				}

				if len(usersForCreate) > 0 {
					if err := tx.Create(&usersForCreate).Error; err != nil {
						return err
					}
				}

				if len(onlineChattersForCreate) > 0 {
					if err := tx.Create(&onlineChattersForCreate).Error; err != nil {
						return err
					}
				}

				return nil
			},
		)

		if err != nil {
			zap.S().Error(err)
			return
		}

		if req.Data.Pagination.Cursor == "" {
			return
		}

		cursor = req.Data.Pagination.Cursor
	}
}
