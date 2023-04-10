package grpc_impl

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/satont/tsuwari/libs/grpc/generated/tokens"
	"github.com/satont/tsuwari/libs/twitch"

	"github.com/gofrs/uuid"
	"github.com/nicklaw5/helix/v2"
	cfg "github.com/satont/tsuwari/libs/config"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/grpc/generated/watched"
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
	twitchClient, err := twitch.NewBotClient(data.BotId, *c.cfg, c.tokensGrpc)
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

				if req.Data.Pagination.Cursor == "" || len(req.Data.Chatters) == 0 {
					break
				}

				cursor = req.Data.Pagination.Cursor
			}

			chattersWg := sync.WaitGroup{}
			for _, chatter := range chatters {
				if chatter.UserID == "" {
					continue
				}
				chattersWg.Add(1)

				go func(chatter helix.ChatChatter) {
					defer chattersWg.Done()
					user := model.Users{}
					err := c.db.
						Where(`"users"."id" = ?`, chatter.UserID).
						Preload("Stats", `"channelId" = ?`, channel).
						Find(&user).Error
					if err != nil {
						c.logger.Sugar().Error(err)
						return
					}

					if user.ID == "" {
						apiKey, _ := uuid.NewV4()
						statsId, _ := uuid.NewV4()
						newUser := &model.Users{
							ID:     chatter.UserID,
							ApiKey: apiKey.String(),
							Stats: &model.UsersStats{
								ID:        statsId.String(),
								UserID:    chatter.UserID,
								ChannelID: channel,
								Messages:  0,
								Watched:   0,
							},
						}

						if err := c.db.Create(&newUser).Error; err != nil {
							c.logger.Sugar().Error(err)
						}
					} else if user.Stats == nil {
						statsId, _ := uuid.NewV4()
						err := c.db.Create(&model.UsersStats{
							ID:        statsId.String(),
							UserID:    chatter.UserID,
							ChannelID: channel,
							Messages:  0,
							Watched:   0,
						}).Error
						if err != nil {
							c.logger.Sugar().Error(err)
						}
					} else {
						time := 5 * time.Minute

						err := c.db.
							Model(&model.UsersStats{}).
							Where("id = ?", user.Stats.ID).
							Updates(map[string]any{
								"watched": user.Stats.Watched + time.Milliseconds(),
							}).Error
						if err != nil {
							c.logger.Sugar().Error(err)
						}
					}
				}(chatter)
			}

			chattersWg.Wait()
		}(channel)
	}

	wg.Wait()

	return &emptypb.Empty{}, nil
}
