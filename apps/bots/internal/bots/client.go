package bots

import (
	"fmt"
	"sync"
	"time"

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
	Twitch     *twitch.Twitch
	ParserGrpc parser.ParserClient
}

func newBot(opts *ClientOpts) *types.BotClient {
	globalRateLimiter, _ := ratelimiting.NewSlidingWindow(100, 30*time.Second)

	client := types.BotClient{
		RateLimiters: types.RateLimiters{
			Global: globalRateLimiter,
		},
		Model: opts.Model,
	}

	botHandlers := handlers.CreateHandlers(&handlers.HandlersOpts{
		DB:         opts.DB,
		Logger:     opts.Logger,
		Cfg:        opts.Cfg,
		BotClient:  &client,
		ParserGrpc: opts.ParserGrpc,
	})

	onRefresh := func(newToken helix.RefreshTokenResponse) {
		opts.DB.Model(&model.Tokens{}).
			Where(`id = ?`, opts.Model.TokenID.String).
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

	client.Api = api

	go func() {
		for {
			token := fmt.Sprintf("oauth:%s", api.Client.GetUserAccessToken())
			if client.Client == nil {
				meReq, err := api.Client.GetUsers(&helix.UsersParams{
					IDs: []string{opts.Model.ID},
				})
				if err != nil {
					panic(err)
				}

				if meReq.Error != "" {
					fmt.Println(meReq.ErrorMessage)
					panic(meReq.Error)
				}

				if len(meReq.Data.Users) == 0 {
					panic("No user found for bot " + opts.Model.ID)
				}

				me := meReq.Data.Users[0]
				client.TwitchUser = &me
				client.Client = irc.NewClient(me.Login, token)
				joinChannels(opts.DB, opts.Cfg, opts.Logger, &client)
			} else {
				client.Client.SetIRCToken(token)
			}
			meReq, err := api.Client.GetUsers(&helix.UsersParams{
				IDs: []string{opts.Model.ID},
			})
			if err != nil {
				panic(err)
			}

			if len(meReq.Data.Users) == 0 {
				panic("No user found for bot " + opts.Model.ID)
			}

			client.OnConnect(botHandlers.OnConnect)
			client.OnSelfJoinMessage(botHandlers.OnSelfJoin)
			client.OnUserStateMessage(botHandlers.OnUserStateMessage)
			client.OnPrivateMessage(botHandlers.OnPrivateMessage)

			err = client.Connect()
			if err != nil {
				opts.Logger.Sugar().Error(err)
			}
		}
	}()

	return &client
}

func joinChannels(db *gorm.DB, cfg *cfg.Config, logger *zap.Logger, botClient *types.BotClient) {
	usersApiService := twitch.NewUserClient(twitch.UsersServiceOpts{
		Db:           db,
		ClientId:     cfg.TwitchClientId,
		ClientSecret: cfg.TwitchClientSecret,
	})

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

			twitchUsersReq, err := botClient.Api.Client.GetUsers(&helix.UsersParams{
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
				api, err := usersApiService.Create(u.ID)
				if err != nil {
					return
				}

				botModRequest, err := api.GetChannelMods(&helix.GetChannelModsParams{
					BroadcasterID: u.ID,
					UserID:        botClient.Model.ID,
				})

				if err != nil || botModRequest.ResponseCommon.StatusCode != 200 {
					isMod = false
				}

				if len(botModRequest.Data.Mods) == 1 {
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
