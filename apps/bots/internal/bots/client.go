package bots

import (
	"fmt"
	"time"

	cfg "github.com/satont/tsuwari/libs/config"

	model "github.com/satont/tsuwari/libs/gomodels"

	"github.com/satont/tsuwari/libs/twitch"

	ratelimiting "github.com/aidenwallis/go-ratelimiting/local"
	irc "github.com/gempir/go-twitch-irc/v3"
	"github.com/nats-io/nats.go"
	"github.com/satont/go-helix/v2"
	"github.com/satont/tsuwari/apps/bots/internal/bots/handlers"
	"github.com/satont/tsuwari/apps/bots/types"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ClientOpts struct {
	DB     *gorm.DB
	Cfg    *cfg.Config
	Logger *zap.Logger
	Model  *model.Bots
	Twitch *twitch.Twitch
	Nats   *nats.Conn
}

func newBot(opts *ClientOpts) *types.BotClient {
	onRefresh := func(newToken helix.RefreshTokenResponse) {
		opts.DB.Model(&model.Tokens{}).
			Where(`id = ?`, opts.Model.ID).
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
		UserAccessToken:  opts.Model.Token.AccessToken,
		UserRefreshToken: opts.Model.Token.RefreshToken,
		OnRefresh:        &onRefresh,
	})

	meReq, err := api.Client.GetUsers(&helix.UsersParams{
		IDs: []string{opts.Model.ID},
	})
	if err != nil {
		panic(err)
	}

	if len(meReq.Data.Users) == 0 {
		panic("No user found for bot " + opts.Model.ID)
	}

	me := meReq.Data.Users[0]

	token := fmt.Sprintf("oauth:%s", api.Client.GetUserAccessToken())
	globalRateLimiter, _ := ratelimiting.NewSlidingWindow(100, 30*time.Second)

	client := types.BotClient{
		Client: irc.NewClient(me.Login, token),
		Api:    api,
		RateLimiters: types.RateLimiters{
			Global: globalRateLimiter,
		},
		Model: opts.Model,
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

	go client.Connect()
	return &client
}
