package bots

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/samber/do"
	"github.com/satont/twir/apps/bots/internal/di"
	"github.com/satont/twir/libs/grpc/generated/events"
	"github.com/satont/twir/libs/grpc/generated/tokens"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/samber/lo"
	cfg "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/grpc/generated/parser"

	model "github.com/satont/twir/libs/gomodels"

	"github.com/satont/twir/libs/twitch"

	ratelimiting "github.com/aidenwallis/go-ratelimiting/local"
	irc "github.com/gempir/go-twitch-irc/v3"
	"github.com/nicklaw5/helix/v2"
	"github.com/satont/twir/apps/bots/internal/bots/handlers"
	"github.com/satont/twir/apps/bots/types"
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
	eventsGrpc := do.MustInvoke[events.EventsClient](di.Provider)

	globalRateLimiter, _ := ratelimiting.NewSlidingWindow(100, 30*time.Second)

	client := types.BotClient{
		RateLimiters: types.RateLimiters{
			Global: globalRateLimiter,
		},
		Model: opts.Model,
	}

	twitchClient, err := twitch.NewBotClient(opts.Model.ID, *opts.Cfg, tokensGrpc)

	meReq, err := twitchClient.GetUsers(
		&helix.UsersParams{
			IDs: []string{opts.Model.ID},
		},
	)
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

	messagesCounter := promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "bots_messages_counter",
			Help: "The total number of processed messages",
			ConstLabels: prometheus.Labels{
				"botName": meReq.Data.Users[0].Login,
				"botId":   meReq.Data.Users[0].ID,
			},
		},
	)

	prometheus.Register(messagesCounter)

	client.Client = irc.NewClient(me.Login, "")
	client.TwitchUser = &me
	client.Client.Capabilities = []string{irc.TagsCapability, irc.MembershipCapability, irc.CommandsCapability}
	client.RateLimiters.Channels = types.ChannelsMap{
		Items: make(map[string]*types.Channel),
	}

	botHandlers := handlers.CreateHandlers(
		&handlers.HandlersOpts{
			DB:         opts.DB,
			Logger:     opts.Logger,
			Cfg:        opts.Cfg,
			BotClient:  &client,
			ParserGrpc: opts.ParserGrpc,
			EventsGrpc: eventsGrpc,
			TokensGrpc: tokensGrpc,
		},
	)

	go func() {
		for {
			token, err := tokensGrpc.RequestBotToken(
				context.Background(), &tokens.GetBotTokenRequest{
					BotId: opts.Model.ID,
				},
			)
			if err != nil {
				panic(err)
			}

			twitchClient.SetUserAccessToken(token.AccessToken)

			//joinChannels(opts.DB, opts.Cfg, opts.Logger, &client)
			client.Client.SetIRCToken(fmt.Sprintf("oauth:%s", token.AccessToken))
			meReq, err := twitchClient.GetUsers(
				&helix.UsersParams{
					IDs: []string{opts.Model.ID},
				},
			)
			if err != nil {
				opts.Logger.Sugar().Error(err)
				return
			}

			if len(meReq.Data.Users) == 0 {
				opts.Logger.Sugar().Error("No user found for bot " + opts.Model.ID)
				return
			}

			client.OnConnect(botHandlers.OnConnect)
			client.OnSelfJoinMessage(botHandlers.OnSelfJoin)
			client.OnUserStateMessage(
				func(message irc.UserStateMessage) {
					defer messagesCounter.Inc()
					if message.User.ID == me.ID && opts.Cfg.AppEnv != "development" {
						return
					}
					botHandlers.OnUserStateMessage(message)
				},
			)
			client.OnUserNoticeMessage(
				func(message irc.UserNoticeMessage) {
					defer messagesCounter.Inc()
					if message.User.ID == me.ID && opts.Cfg.AppEnv != "development" {
						return
					}
					botHandlers.OnMessage(
						&handlers.Message{
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
							Tags:    message.Tags,
						},
					)
				},
			)
			client.OnPrivateMessage(
				func(message irc.PrivateMessage) {
					defer messagesCounter.Inc()
					if message.User.ID == me.ID && opts.Cfg.AppEnv != "development" {
						return
					}
					botHandlers.OnMessage(
						&handlers.Message{
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
							Tags:    message.Tags,
						},
					)
				},
			)
			client.OnClearChatMessage(
				func(message irc.ClearChatMessage) {
					if message.TargetUserID != "" {
						return
					}
					opts.DB.Create(
						&model.ChannelsEventsListItem{
							ID:        uuid.New().String(),
							ChannelID: message.RoomID,
							Type:      model.ChannelEventListItemTypeChatClear,
							CreatedAt: time.Now().UTC(),
							Data:      &model.ChannelsEventsListItemData{},
						},
					)
					eventsGrpc.ChatClear(
						context.Background(), &events.ChatClearMessage{
							BaseInfo: &events.BaseInfo{ChannelId: message.RoomID},
						},
					)
				},
			)
			client.OnUserNoticeMessage(botHandlers.OnNotice)
			client.OnUserJoinMessage(botHandlers.OnUserJoin)
			joinChannels(opts.DB, opts.Cfg, opts.Logger, &client)

			if err = client.Connect(); err != nil {
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

	var botChannels []model.Channels
	db.
		Where(
			`"botId" = ? AND "isEnabled" = ? AND "isBanned" = ? AND "isTwitchBanned" = ?`,
			botClient.Model.ID,
			true,
			false,
			false,
		).Find(&botChannels)
	channelsChunks := lo.Chunk(botChannels, 100)

	var twitchUsers []helix.User
	var twitchUsersMU sync.Mutex
	var wg sync.WaitGroup
	wg.Add(len(channelsChunks))

	for _, chunk := range channelsChunks {
		go func(chunk []model.Channels) {
			defer wg.Done()
			usersIds := lo.Map(
				chunk,
				func(item model.Channels, _ int) string {
					return item.ID
				},
			)

			twitchUsersReq, err := twitchClient.GetUsers(
				&helix.UsersParams{
					IDs: usersIds,
				},
			)
			if err != nil {
				panic(err)
			}
			twitchUsersMU.Lock()
			twitchUsers = append(twitchUsers, twitchUsersReq.Data.Users...)
			twitchUsersMU.Unlock()
		}(chunk)
	}

	wg.Wait()

	for _, u := range twitchUsers {
		botClient.Join(u.Login)
	}
}
