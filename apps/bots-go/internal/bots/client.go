package bots

import (
	"fmt"
	"time"
	cfg "tsuwari/config"
	model "tsuwari/models"
	"tsuwari/twitch"

	ratelimiting "github.com/aidenwallis/go-ratelimiting/local"
	irc "github.com/gempir/go-twitch-irc/v3"
	"github.com/nats-io/nats.go"
	"github.com/satont/go-helix/v2"
	"github.com/satont/tsuwari/apps/bots/internal/handlers"
	"github.com/satont/tsuwari/apps/bots/types"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ClientOpts struct {
	DB     *gorm.DB
	Cfg    *cfg.Config
	Logger *zap.Logger
	Bot    *model.Bots
	Twitch *twitch.Twitch
	Nats   *nats.Conn
}

func newBot(opts *ClientOpts) *types.BotClient {
	onRefresh := func(newToken helix.RefreshTokenResponse) {
		opts.DB.Model(&model.Tokens{}).
			Where(`id = ?`, opts.Bot.ID).
			Select("*").
			Updates(map[string]any{
				"accessToken":         newToken.Data.AccessToken,
				"refreshToken":        newToken.Data.RefreshToken,
				"expiresIn":           int32(newToken.Data.ExpiresIn),
				"obtainmentTimestamp": time.Now(),
			})
	}
	api := twitch.NewClient(&helix.Options{
		ClientID:         opts.Cfg.TwitchClientId,
		ClientSecret:     opts.Cfg.TwitchClientSecret,
		UserAccessToken:  opts.Bot.Token.AccessToken,
		UserRefreshToken: opts.Bot.Token.RefreshToken,
		OnRefresh:        &onRefresh,
	})

	meReq, err := api.Client.GetUsers(&helix.UsersParams{
		IDs: []string{opts.Bot.ID},
	})
	if err != nil {
		panic(err)
	}

	if len(meReq.Data.Users) == 0 {
		panic("No user found for bot " + opts.Bot.ID)
	}

	me := meReq.Data.Users[0]

	token := fmt.Sprintf("oauth:%s", api.Client.GetUserAccessToken())
	globalRateLimiter, _ := ratelimiting.NewSlidingWindow(100, 30*time.Second)

	client := types.BotClient{
		Client: irc.NewClient(me.Login, token),
		Api:    api,
		RateLimiters: types.RateLimiters{
			Global:   globalRateLimiter,
			Channels: make(map[string]ratelimiting.SlidingWindow),
		},
	}

	botHandlers := handlers.CreateHandlers(&handlers.HandlersOpts{
		DB:        opts.DB,
		Logger:    opts.Logger,
		Cfg:       opts.Cfg,
		BotClient: &client,
		Nats:      opts.Nats,
	})

	client.OnConnect(botHandlers.OnConnect)
	client.OnSelfJoinMessage(botHandlers.OnSelfJoin)
	client.OnUserStateMessage(botHandlers.OnUserStateMessage)
	client.OnPrivateMessage(botHandlers.OnPrivateMessage)

	err = client.Connect()
	if err != nil {
		panic(err)
	}

	return &client
}
