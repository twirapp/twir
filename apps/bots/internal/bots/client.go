package bots

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/Adeithe/go-twitch/irc"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/satont/twir/libs/grpc/generated/websockets"
	"github.com/satont/twir/libs/logger"
	"github.com/satont/twir/libs/utils"

	"github.com/satont/twir/libs/grpc/generated/events"
	"github.com/satont/twir/libs/grpc/generated/tokens"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/samber/lo"
	cfg "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/grpc/generated/parser"

	model "github.com/satont/twir/libs/gomodels"

	tokensLib "github.com/satont/twir/libs/twitch"

	"github.com/Adeithe/go-twitch"
	ratelimiting "github.com/aidenwallis/go-ratelimiting/local"
	"github.com/nicklaw5/helix/v2"
	"github.com/satont/twir/apps/bots/internal/bots/handlers"
	"github.com/satont/twir/apps/bots/types"
	"gorm.io/gorm"
)

type ClientOpts struct {
	DB             *gorm.DB
	Cfg            cfg.Config
	Logger         logger.Logger
	Model          *model.Bots
	ParserGrpc     parser.ParserClient
	TokensGrpc     tokens.TokensClient
	EventsGrpc     events.EventsClient
	WebsocketsGrpc websockets.WebsocketClient
	Redis          *redis.Client
}

func newBot(opts ClientOpts) *types.BotClient {
	globalRateLimiter, _ := ratelimiting.NewSlidingWindow(100, 30*time.Second)

	client := types.BotClient{
		RateLimiters: types.RateLimiters{
			Global: globalRateLimiter,
		},
		Model: opts.Model,
	}

	twitchClient, err := tokensLib.NewBotClient(opts.Model.ID, opts.Cfg, opts.TokensGrpc)
	if err != nil {
		panic(err)
	}

	meReq, err := twitchClient.GetUsers(
		&helix.UsersParams{
			IDs: []string{opts.Model.ID},
		},
	)
	if err != nil || len(meReq.Data.Users) == 0 {
		panic("No user found for bot " + opts.Model.ID)
	}

	me := meReq.Data.Users[0]
	messagesCounter := createPromCounter(me)

	token, err := opts.TokensGrpc.RequestBotToken(
		context.TODO(),
		&tokens.GetBotTokenRequest{
			BotId: opts.Model.ID,
		},
	)
	if err != nil {
		panic(err)
	}

	client.Writer = &irc.Conn{}
	client.Writer.SetLogin(me.Login, fmt.Sprintf("oauth:%s", token.AccessToken))
	if err := client.Writer.Connect(); err != nil {
		panic("failed to start writer")
	}
	opts.Logger.Info("IRC writer connected", slog.String("botName", me.Login))

	client.TwitchUser = &me
	client.RateLimiters.Channels = types.ChannelsMap{
		Items: make(map[string]*types.Channel),
	}

	joinOpts := joinChannelOpts{
		db:         opts.DB,
		config:     opts.Cfg,
		logger:     opts.Logger,
		botModel:   opts.Model,
		tokensGrpc: opts.TokensGrpc,
	}

	botHandlers := handlers.CreateHandlers(
		&handlers.Opts{
			DB:             opts.DB,
			Logger:         opts.Logger,
			Cfg:            opts.Cfg,
			BotClient:      &client,
			ParserGrpc:     opts.ParserGrpc,
			EventsGrpc:     opts.EventsGrpc,
			TokensGrpc:     opts.TokensGrpc,
			WebsocketsGrpc: opts.WebsocketsGrpc,
			Redis:          opts.Redis,
		},
	)

	client.Reader = twitch.IRC()
	client.Reader.SetMaxChannelsPerShard(50)

	client.Reader.OnShardRawMessage(
		func(i int, m irc.Message) {
			if m.Command == "CLEARCHAT" && m.Tags["target-user-id"] == "" {
				channelId := m.Tags["room-id"]
				opts.DB.Create(
					&model.ChannelsEventsListItem{
						ID:        uuid.New().String(),
						ChannelID: channelId,
						Type:      model.ChannelEventListItemTypeChatClear,
						CreatedAt: time.Now().UTC(),
						Data:      &model.ChannelsEventsListItemData{},
					},
				)
				opts.EventsGrpc.ChatClear(
					context.Background(),
					&events.ChatClearMessage{
						BaseInfo: &events.BaseInfo{ChannelId: channelId},
					},
				)
			}
		},
	)
	client.Reader.OnShardMessage(
		func(i int, m irc.ChatMessage) {
			if fmt.Sprint(m.Sender.ID) == me.ID && opts.Cfg.AppEnv != "development" {
				return
			}

			defer messagesCounter.Inc()
			botHandlers.OnMessage(
				&handlers.Message{
					ID: m.ID,
					Channel: handlers.MessageChannel{
						ID:   fmt.Sprint(m.ChannelID),
						Name: m.Channel,
					},
					User: handlers.MessageUser{
						ID:          fmt.Sprint(m.Sender.ID),
						Name:        m.Sender.Username,
						DisplayName: m.Sender.DisplayName,
						Badges:      m.Sender.Badges,
					},
					Message: m.Text,
					Emotes:  botHandlers.ParseEmotes(m.Text, m.IRCMessage.Tags["emotes"]),
					Tags:    m.IRCMessage.Tags,
				},
			)
		},
	)
	client.Reader.OnShardChannelJoin(
		func(i int, channel, user string) {
			botHandlers.OnUserJoin(
				handlers.OnUserJoinOpts{
					Channel: channel,
					User:    user,
				},
			)
		},
	)
	client.Reader.OnShardChannelUserNotice(
		func(i int, m irc.UserNotice) {
			botHandlers.OnNotice(
				handlers.OnNoticeOpts{
					Type:              m.Type,
					Message:           m.Message,
					UserID:            fmt.Sprint(m.Sender.ID),
					ChannelID:         m.IRCMessage.Tags["room-id"],
					SenderUserLogin:   m.Sender.Username,
					SenderDisplayName: m.Sender.DisplayName,
				},
				m.IRCMessage.Tags,
			)
		},
	)
	client.Reader.OnShardLatencyUpdate(
		func(i int, duration time.Duration) {
			opts.Logger.Info(
				"Shard latency update",
				slog.Int("shardId", i),
				slog.Duration(
					"latency",
					duration,
				),
			)
		},
	)
	client.Reader.OnShardReconnect(
		func(_ int) {
			joinChannels(joinOpts, client.Reader.Join)
		},
	)

	client.Writer.OnRawMessage(
		func(message irc.Message) {
			if message.Command == "USERSTATE" {
				botHandlers.OnUserStateMessage(
					handlers.OnUserStateMessageOpts{
						Moderator:   message.Tags["mod"],
						Broadcaster: message.Tags["broadcaster"],
						Channel:     message.Params[0][1:],
					},
				)
			}
		},
	)
	client.Writer.OnReconnect(
		func() {
			joinChannels(joinOpts, client.Writer.Join)
		},
	)

	expiresIn := token.ExpiresIn

	go func() {
		for {
			time.Sleep(time.Duration(expiresIn-60) * time.Second)
			newToken, err := opts.TokensGrpc.RequestBotToken(
				context.TODO(),
				&tokens.GetBotTokenRequest{
					BotId: opts.Model.ID,
				},
			)
			if err != nil {
				opts.Logger.Error("cannot fetch token", slog.Any("err", err))
				expiresIn = 5
				return
			}
			expiresIn = newToken.ExpiresIn
			client.Writer.SetLogin(me.Login, fmt.Sprintf("oauth:%s", newToken.AccessToken))
			opts.Logger.Info("Updated writer token")
		}
	}()

	joinChannels(joinOpts, client.Reader.Join)
	joinChannels(joinOpts, client.Writer.Join)

	opts.Logger.Info(
		"IRC reader connected", slog.String("botName", me.Login),
	)

	// reader.OnShardMessage(onShardMessage)

	// go func() {
	// 	for {
	// 		client.OnUserStateMessage(
	// 			func(message irc.UserStateMessage) {
	// 				defer messagesCounter.Inc()
	// 				if message.User.ID == me.ID && opts.Cfg.AppEnv != "development" {
	// 					return
	// 				}
	// 				botHandlers.OnUserStateMessage(message)
	// 			},
	// 		)
	// 		client.OnUserNoticeMessage(botHandlers.OnNotice)
	// 	}
	// }()

	return &client
}

type joinChannelOpts struct {
	db         *gorm.DB
	config     cfg.Config
	logger     logger.Logger
	botModel   *model.Bots
	tokensGrpc tokens.TokensClient
}

func getChannelsNames(opts joinChannelOpts) []string {
	ctx, closeFunc := context.WithTimeout(context.TODO(), 10*time.Second)
	defer closeFunc()
	twitchClient, err := tokensLib.NewBotClientWithContext(
		ctx,
		opts.botModel.ID,
		opts.config,
		opts.tokensGrpc,
	)
	if err != nil {
		panic(err)
	}

	var botChannels []model.Channels
	opts.db.
		Where(
			`"botId" = ? AND "isEnabled" = ? AND "isBanned" = ? AND "isTwitchBanned" = ?`,
			opts.botModel.ID,
			true,
			false,
			false,
		).Find(&botChannels)
	channelsChunks := lo.Chunk(botChannels, 100)

	var twitchUsers []helix.User
	var twitchUsersMU sync.Mutex

	wg := utils.NewGoroutinesGroup()

	for _, chunk := range channelsChunks {
		chunk := chunk

		wg.Go(
			func() {
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
					opts.logger.Error("cannot fetch users", slog.Any("err", err))
					return
				}
				twitchUsersMU.Lock()
				twitchUsers = append(twitchUsers, twitchUsersReq.Data.Users...)
				twitchUsersMU.Unlock()
			},
		)
	}

	wg.Wait()

	names := make([]string, len(twitchUsers))

	for i, u := range twitchUsers {
		names[i] = u.Login
	}

	return names
}

func createPromCounter(user helix.User) prometheus.Counter {
	messagesCounter := promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "bots_messages_counter",
			Help: "The total number of processed messages",
			ConstLabels: prometheus.Labels{
				"botName": user.Login,
				"botId":   user.ID,
			},
		},
	)

	prometheus.Register(messagesCounter)

	return messagesCounter
}

func joinChannels(opts joinChannelOpts, joinFunc func(names ...string) error) {
	channelsNames := getChannelsNames(opts)

	if err := joinFunc(channelsNames...); err != nil {
		opts.logger.Error("cannot join channels", slog.Any("error", err))
	}
}
