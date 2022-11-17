package handlers

import (
	"sync"
	"time"

	"github.com/satont/tsuwari/apps/bots/types"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/twitch"

	ratelimiting "github.com/aidenwallis/go-ratelimiting/local"
	"github.com/samber/lo"
	"github.com/satont/go-helix/v2"
)

func (c *Handlers) OnConnect() {
	usersApiService := twitch.NewUserClient(twitch.UsersServiceOpts{
		Db:           c.db,
		ClientId:     c.cfg.TwitchClientId,
		ClientSecret: c.cfg.TwitchClientSecret,
	})

	c.logger.Sugar().
		Infow("Bot connected to twitch", "botId", c.BotClient.Model.ID, "botName", c.BotClient.TwitchUser.Login)

	c.BotClient.RateLimiters.Channels = types.ChannelsMap{
		Items: make(map[string]ratelimiting.SlidingWindow),
	}

	twitchUsers := []helix.User{}
	twitchUsersMU := sync.Mutex{}

	botChannels := []model.Channels{}

	c.db.
		Where(
			`"botId" = ? AND "isEnabled" = ? AND "isBanned" = ? AND "isTwitchBanned" = ?`,
			c.BotClient.Model.ID,
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

			twitchUsersReq, err := c.BotClient.Api.Client.GetUsers(&helix.UsersParams{
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
			err := c.db.Where("id = ?", u.ID).Preload("Token").Find(&dbUser).Error
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
					UserID:        c.BotClient.Model.ID,
				})

				if err != nil || botModRequest.ResponseCommon.StatusCode != 200 {
					isMod = false
				}

				if len(botModRequest.Data.Mods) == 1 {
					isMod = true
				}
			}

			var limiter ratelimiting.SlidingWindow
			if isMod {
				l, _ := ratelimiting.NewSlidingWindow(20, 30*time.Second)
				limiter = l
			} else {
				l, _ := ratelimiting.NewSlidingWindow(1, 2*time.Second)
				limiter = l
			}

			c.BotClient.RateLimiters.Channels.Lock()
			c.BotClient.RateLimiters.Channels.Items[u.Login] = limiter
			c.BotClient.RateLimiters.Channels.Unlock()

			c.BotClient.Join(u.Login)
		}(u)
	}

	wg.Wait()
}
