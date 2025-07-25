package timers

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	config "github.com/twirapp/twir/libs/config"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/logger"
	"github.com/twirapp/twir/libs/twitch"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type OnlineUsersOpts struct {
	fx.In
	Lc fx.Lifecycle

	Logger logger.Logger
	Config config.Config

	Gorm    *gorm.DB
	TwirBus *buscore.Bus
}

type onlineUsers struct {
	config  config.Config
	logger  logger.Logger
	db      *gorm.DB
	twirBus *buscore.Bus
}

func NewOnlineUsers(opts OnlineUsersOpts) {
	timeTick := 15 * time.Second
	if opts.Config.AppEnv == "production" {
		timeTick = 5 * time.Minute
	}
	ticker := time.NewTicker(timeTick)

	ctx, cancel := context.WithCancel(context.Background())

	s := &onlineUsers{
		config:  opts.Config,
		logger:  opts.Logger,
		db:      opts.Gorm,
		twirBus: opts.TwirBus,
	}

	opts.Lc.Append(
		fx.Hook{
			OnStart: func(_ context.Context) error {
				go func() {
					for {
						select {
						case <-ctx.Done():
							ticker.Stop()
							return
						case <-ticker.C:
							s.updateOnlineUsers(ctx)
						}
					}
				}()

				return nil
			},
			OnStop: func(_ context.Context) error {
				cancel()
				return nil
			},
		},
	)
}

func (c *onlineUsers) updateOnlineUsers(ctx context.Context) {
	streams, err := c.getStreams(ctx)
	if err != nil {
		c.logger.Error("cannot get streams", slog.Any("err", err))
		return
	}

	var wg sync.WaitGroup
	for _, stream := range streams {
		if c.shouldSkipStream(stream) {
			continue
		}

		userId := stream.UserId
		wg.Add(1)

		go func() {
			defer wg.Done()
			if updateErr := c.updateStreamUsers(ctx, userId); updateErr != nil {
				c.logger.Error("cannot update stream users", slog.Any("err", updateErr))
			}
		}()
	}

	wg.Wait()
}

func (c *onlineUsers) getStreams(
	ctx context.Context,
) ([]*model.ChannelsStreams, error) {
	var streams []*model.ChannelsStreams
	err := c.db.WithContext(ctx).Preload("Channel").Preload("Channel.User").Find(&streams).Error
	return streams, err
}

func (c *onlineUsers) shouldSkipStream(stream *model.ChannelsStreams) bool {
	return stream.Channel == nil || (!stream.Channel.IsEnabled || stream.Channel.User.IsBanned)
}

func (c *onlineUsers) updateStreamUsers(
	ctx context.Context,
	broadcasterID string,
) error {
	twitchClient, err := twitch.NewUserClientWithContext(
		ctx,
		broadcasterID,
		c.config,
		c.twirBus,
	)
	if err != nil {
		return err
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
			return fmt.Errorf("cannot get channel chat chatters: %w", err)
		}
		if req.ErrorMessage != "" {
			return fmt.Errorf("cannot get channel chat chatters: %s", req.ErrorMessage)
		}

		chatters := req.Data.Chatters
		if len(chatters) == 0 {
			return nil
		}

		usersIdsForRequest := lo.Map(
			chatters,
			func(chatter helix.ChatChatter, _ int) string {
				return chatter.UserID
			},
		)

		err = c.db.WithContext(ctx).Transaction(
			func(tx *gorm.DB) error {
				var existedUsers []model.Users
				if err := tx.
					Select("id").
					Where("id IN ?", usersIdsForRequest).
					Find(&existedUsers).Error; err != nil {
					return fmt.Errorf("cannot get existed users: %w", err)
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
					return fmt.Errorf("cannot delete online users: %w", err)
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
						return fmt.Errorf("cannot create users: %w", err)
					}
				}

				if len(onlineChattersForCreate) > 0 {
					if err := tx.Create(&onlineChattersForCreate).Error; err != nil {
						return fmt.Errorf("cannot create online users: %w", err)
					}
				}

				return nil
			},
		)

		if err != nil {
			return fmt.Errorf("cannot update stream users: %w", err)
		}

		if req.Data.Pagination.Cursor == "" {
			return nil
		}

		cursor = req.Data.Pagination.Cursor
	}
}
