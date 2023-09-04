package grpc_impl

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/satont/twir/libs/grpc/generated/tokens"
	"github.com/satont/twir/libs/twitch"

	"github.com/nicklaw5/helix/v2"
	cfg "github.com/satont/twir/libs/config"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/generated/watched"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type WatchedGrpcServerOpts struct {
	Db         *gorm.DB
	Cfg        *cfg.Config
	Logger     *zap.Logger
	TokensGrpc tokens.TokensClient
}

type WatchedGrpcServer struct {
	watched.UnimplementedWatchedServer

	db         *gorm.DB
	cfg        *cfg.Config
	logger     *zap.Logger
	tokensGrpc tokens.TokensClient
}

func New(opts *WatchedGrpcServerOpts) *WatchedGrpcServer {
	return &WatchedGrpcServer{
		db:         opts.Db,
		cfg:        opts.Cfg,
		logger:     opts.Logger,
		tokensGrpc: opts.TokensGrpc,
	}
}

func (c *WatchedGrpcServer) IncrementByChannelId(
	ctx context.Context,
	data *watched.Request,
) (*emptypb.Empty, error) {
	twitchClient, err := twitch.NewBotClientWithContext(ctx, data.BotId, *c.cfg, c.tokensGrpc)
	if err != nil {
		c.logger.Sugar().Error(err)
		return nil, errors.New("cannot create api for bot")
	}

	wg := sync.WaitGroup{}

	for _, channel := range data.ChannelsId {
		wg.Add(1)

		go func(channel string) {
			defer wg.Done()
			chatters := make([]helix.ChatChatter, 0)
			cursor := ""

			for {
				reqParams := &helix.GetChatChattersParams{
					BroadcasterID: channel,
					ModeratorID:   data.BotId,
					First:         "1000",
				}

				if cursor != "" {
					reqParams.After = cursor
				}

				req, err := twitchClient.GetChannelChatChatters(reqParams)
				if err != nil {
					break
				}

				chatters = append(chatters, req.Data.Chatters...)

				if len(req.Data.Chatters) == 0 {
					break
				}

				cursor = req.Data.Pagination.Cursor

				chattersIds := lo.Map(
					chatters, func(item helix.ChatChatter, _ int) string {
						return item.UserID
					},
				)
				var existedChatters []model.Users
				if err := c.db.
					WithContext(ctx).
					Where(`"users"."id" IN ?`, chattersIds).
					Joins("Stats", c.db.Where(&model.UsersStats{ChannelID: channel}), channel).
					Find(&existedChatters).
					Error; err != nil {
					c.logger.Sugar().Error(err)
					return
				}

				usersForCreate := make([]model.Users, 0, len(chatters))
				for _, chatter := range chatters {
					if chatter.UserID == "" {
						continue
					}

					var existedChatter *model.Users
					for _, item := range existedChatters {
						if item.ID == chatter.UserID {
							existedChatter = &item
							break
						}
					}

					if existedChatter == nil {
						newUser := model.Users{
							ID:     chatter.UserID,
							ApiKey: uuid.New().String(),
							Stats: &model.UsersStats{
								ID:        uuid.New().String(),
								UserID:    chatter.UserID,
								ChannelID: channel,
								Watched:   0,
							},
						}
						usersForCreate = append(usersForCreate, newUser)
					} else {
						if existedChatter.Stats == nil {
							err := c.db.Create(
								&model.UsersStats{
									ID:        uuid.New().String(),
									UserID:    chatter.UserID,
									ChannelID: channel,
									Watched:   0,
								},
							).Error
							if err != nil {
								c.logger.Sugar().Error(err)
							}
						} else {
							incTime := 5 * time.Minute

							err := c.db.
								Model(&model.UsersStats{}).
								Where("id = ?", existedChatter.Stats.ID).
								Updates(
									map[string]any{
										"watched": existedChatter.Stats.Watched + incTime.Milliseconds(),
									},
								).Error
							if err != nil {
								c.logger.Sugar().Error(err)
							}
						}
					}
				}

				if len(usersForCreate) > 0 {
					err := c.db.Create(&usersForCreate).Error
					if err != nil {
						c.logger.Sugar().Error(err)
					}
				}

				if req.Data.Pagination.Cursor == "" {
					break
				}
			}
		}(channel)
	}

	wg.Wait()

	return &emptypb.Empty{}, nil
}
