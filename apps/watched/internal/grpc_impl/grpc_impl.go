package grpc_impl

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/gofrs/uuid"
	"github.com/satont/go-helix/v2"
	cfg "github.com/satont/tsuwari/libs/config"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/grpc/generated/watched"
	"github.com/satont/tsuwari/libs/twitch"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type WatchedGrpcServerOpts struct {
	Db     *gorm.DB
	Cfg    *cfg.Config
	Logger *zap.Logger
}

type WatchedGrpcServer struct {
	watched.UnimplementedWatchedServer

	db     *gorm.DB
	cfg    *cfg.Config
	logger *zap.Logger
}

func New(opts *WatchedGrpcServerOpts) *WatchedGrpcServer {
	return &WatchedGrpcServer{
		db:     opts.Db,
		cfg:    opts.Cfg,
		logger: opts.Logger,
	}
}

func (c *WatchedGrpcServer) IncrementByChannelId(
	ctx context.Context,
	data *watched.Request,
) (*emptypb.Empty, error) {
	twitch := twitch.NewUserClient(twitch.UsersServiceOpts{
		Db:           c.db,
		ClientId:     c.cfg.TwitchClientId,
		ClientSecret: c.cfg.TwitchClientSecret,
	})

	api, err := twitch.CreateBot(data.BotId)
	if err != nil {
		fmt.Println(err)
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

				req, err := api.GetChannelChatChatters(reqParams)
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
				chattersWg.Add(1)

				go func(chatter helix.ChatChatter) {
					defer chattersWg.Done()
					user := model.Users{}
					err := c.db.
						Where(`"users"."id" = ? AND "Stats"."channelId" = ?`, chatter.UserID, channel).
						Joins("Stats").
						Find(&user).Error
					if err != nil {
						c.logger.Sugar().Error(err)
						return
					}

					if user.ID == "" {
						err = c.db.Transaction(func(tx *gorm.DB) error {
							apiKey, _ := uuid.NewV4()
							user := &model.Users{
								ID:     chatter.UserID,
								ApiKey: apiKey.String(),
							}
							if err := tx.Create(user).Error; err != nil {
								return err
							}

							statsId, _ := uuid.NewV4()
							stats := &model.UsersStats{
								ID:        statsId.String(),
								UserID:    chatter.UserID,
								ChannelID: channel,
								Messages:  0,
								Watched:   0,
							}
							if err := tx.Create(stats).Error; err != nil {
								return err
							}

							return nil
						})

						if err != nil {
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
						err := c.db.Model(&model.UsersStats{}).
							Where("id = ?", user.Stats.ID).Select("*").
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
