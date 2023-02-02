package bots

import (
	"context"
	"fmt"
	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/bots/internal/di"
	"github.com/satont/tsuwari/libs/grpc/generated/tokens"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/samber/lo"
	cfg "github.com/satont/tsuwari/libs/config"
	"github.com/satont/tsuwari/libs/grpc/generated/parser"

	model "github.com/satont/tsuwari/libs/gomodels"

	"github.com/satont/tsuwari/libs/twitch"

	ratelimiting "github.com/aidenwallis/go-ratelimiting/local"
	irc "github.com/gempir/go-twitch-irc/v3"
	"github.com/satont/go-helix/v2"
	"github.com/satont/tsuwari/apps/bots/internal/bots/handlers"
	"github.com/satont/tsuwari/apps/bots/pkg/utils"
	"github.com/satont/tsuwari/apps/bots/types"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ClientOpts struct {
	DB         *gorm.DB
	Cfg        *cfg.Config
	Logger     *zap.Logger
	Model      *model.Bots
	ParserGrpc parser.ParserClient
}

func newBot(opts *ClientOpts) *types.BotClient {
	tokensGrpc := do.MustInvoke[tokens.TokensClient](di.Provider)

	globalRateLimiter, _ := ratelimiting.NewSlidingWindow(100, 30*time.Second)

	client := types.BotClient{
		RateLimiters: types.RateLimiters{
			Global: globalRateLimiter,
		},
		Model: opts.Model,
	}

	twitchClient, err := twitch.NewBotClient(opts.Model.ID, *opts.Cfg, tokensGrpc)

	meReq, err := twitchClient.GetUsers(&helix.UsersParams{
		IDs: []string{opts.Model.ID},
	})
	if err != nil {
		panic(err)
	}

	if meReq.Error != "" {
		fmt.Println(opts.Model.ID + " " + meReq.ErrorMessage)
		panic(meReq.Error)
	}

	if len(meReq.Data.Users) == 0 {
		panic("No user found for bot " + opts.Model.ID)
	}

	me := meReq.Data.Users[0]

	messagesCounter := promauto.NewCounter(prometheus.CounterOpts{
		Name: "bots_messages_counter",
		Help: "The total number of processed messages",
		ConstLabels: prometheus.Labels{
			"botName": meReq.Data.Users[0].Login,
			"botId":   meReq.Data.Users[0].ID,
		},
	})

	prometheus.Register(messagesCounter)

	client.Client = irc.NewClient(me.Login, "")
	client.TwitchUser = &me
	botHandlers := handlers.CreateHandlers(&handlers.HandlersOpts{
		DB:         opts.DB,
		Logger:     opts.Logger,
		Cfg:        opts.Cfg,
		BotClient:  &client,
		ParserGrpc: opts.ParserGrpc,
	})

	go func() {
		for {
			token, err := tokensGrpc.RequestBotToken(context.Background(), &tokens.GetBotTokenRequest{
				BotId: opts.Model.ID,
			})

			twitchClient.SetUserAccessToken(token.AccessToken)

			if err != nil {
				panic(err)
			}

			joinChannels(opts.DB, opts.Cfg, opts.Logger, &client)
			client.Client.SetIRCToken(fmt.Sprintf("oauth:%s", token.AccessToken))
			meReq, err := twitchClient.GetUsers(&helix.UsersParams{
				IDs: []string{opts.Model.ID},
			})
			if err != nil {
				return
			}

			if len(meReq.Data.Users) == 0 {
				panic("No user found for bot " + opts.Model.ID)
			}

			client.OnConnect(botHandlers.OnConnect)
			client.OnSelfJoinMessage(botHandlers.OnSelfJoin)
			client.OnUserStateMessage(func(message irc.UserStateMessage) {
				defer messagesCounter.Inc()
				if message.User.ID == me.ID && opts.Cfg.AppEnv != "development" {
					return
				}
				botHandlers.OnUserStateMessage(message)
			})
			client.OnUserNoticeMessage(func(message irc.UserNoticeMessage) {
				defer messagesCounter.Inc()
				if message.User.ID == me.ID && opts.Cfg.AppEnv != "development" {
					return
				}
				botHandlers.OnMessage(handlers.Message{
					ID: message.ID,
					Channel: handlers.MessageChannel{
						ID:   message.RoomID,
						Name: message.Channel,
					},
					User: handlers.MessageUser{
						ID:          message.User.ID,
						Name:        message.User.Name,
						DisplayName: message.User.DisplayName,
						Badges:      message.User.Badges,
					},
					Message: message.Message,
					Emotes:  message.Emotes,
				})
			})
			client.OnPrivateMessage(func(message irc.PrivateMessage) {
				defer messagesCounter.Inc()
				if message.User.ID == me.ID && opts.Cfg.AppEnv != "development" {
					return
				}
				botHandlers.OnMessage(handlers.Message{
					ID: message.ID,
					Channel: handlers.MessageChannel{
						ID:   message.RoomID,
						Name: message.Channel,
					},
					User: handlers.MessageUser{
						ID:          message.User.ID,
						Name:        message.User.Name,
						DisplayName: message.User.DisplayName,
						Badges:      message.User.Badges,
					},
					Message: message.Message,
					Emotes:  message.Emotes,
				})
			})

			err = client.Connect()
			if err != nil {
				opts.Logger.Sugar().Error(err)
			}
		}
	}()

	return &client
}

func joinChannels(db *gorm.DB, cfg *cfg.Config, logger *zap.Logger, botClient *types.BotClient) {
	tokensGrpc := do.MustInvoke[tokens.TokensClient](di.Provider)

	twitchClient, err := twitch.NewBotClient(botClient.Model.ID, *cfg, tokensGrpc)

	if err != nil {
		panic(err)
	}

	botClient.RateLimiters.Channels = types.ChannelsMap{
		Items: make(map[string]*types.Channel),
	}

	twitchUsers := []helix.User{}
	twitchUsersMU := sync.Mutex{}

	botChannels := []model.Channels{}

	db.
		Where(
			`"botId" = ? AND "isEnabled" = ? AND "isBanned" = ? AND "isTwitchBanned" = ?`,
			botClient.Model.ID,
			true,
			false,
			false,
		).Find(&botChannels)

	channelsChunks := lo.Chunk(botChannels, 100)
	wg := sync.WaitGroup{}
	wg.Add(len(channelsChunks))

	for _, chunk := range channelsChunks {
		go func(chunk []model.Channels) {
			defer wg.Done()
			usersIds := lo.Map(chunk, func(item model.Channels, _ int) string {
				return item.ID
			})

			twitchUsersReq, err := twitchClient.GetUsers(&helix.UsersParams{
				IDs: usersIds,
			})
			if err != nil {
				panic(err)
			}
			twitchUsersMU.Lock()
			twitchUsers = append(twitchUsers, twitchUsersReq.Data.Users...)
			twitchUsersMU.Unlock()
		}(chunk)
	}

	wg.Wait()

	wg = sync.WaitGroup{}

	for _, u := range twitchUsers {
		wg.Add(1)
		go func(u helix.User) {
			defer wg.Done()

			dbUser := model.Users{}
			err := db.Where("id = ?", u.ID).Preload("Token").Find(&dbUser).Error
			if err != nil {
				return
			}
			isMod := false

			if dbUser.ID != "" && dbUser.Token != nil {
				botModRequest, err := twitchClient.GetModerators(&helix.GetModeratorsParams{
					BroadcasterID: u.ID,
					UserIDs:       []string{botClient.Model.ID},
				})

				if err != nil || botModRequest.ResponseCommon.StatusCode != 200 {
					isMod = false
				}

				if len(botModRequest.Data.Moderators) == 1 {
					isMod = true
				}
			}

			limiter := utils.CreateBotLimiter(isMod)

			botClient.RateLimiters.Channels.Lock()
			botClient.RateLimiters.Channels.Items[u.Login] = &types.Channel{
				IsMod:   isMod,
				Limiter: limiter,
			}
			botClient.RateLimiters.Channels.Unlock()

			botClient.Join(u.Login)
		}(u)
	}

	wg.Wait()
}
