package handlers

import (
	"log/slog"

	"github.com/nicklaw5/helix/v2"
	"github.com/satont/twir/apps/bots/pkg/utils"
	"github.com/satont/twir/apps/bots/types"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/twitch"
)

type OnUserStateMessageOpts struct {
	Moderator   string
	Broadcaster string
	Channel     string
}

func (c *Handlers) OnUserStateMessage(msg OnUserStateMessageOpts) {
	twitchClient, err := twitch.NewBotClient(c.BotClient.Model.ID, c.cfg, c.tokensGrpc)
	if err != nil {
		panic(err)
	}

	delete(c.BotClient.RateLimiters.Channels.Items, msg.Channel)

	isMod := msg.Moderator == "1" || msg.Broadcaster == "1"

	twitchReq, err := twitchClient.GetUsers(&helix.UsersParams{Logins: []string{msg.Channel}})
	if len(twitchReq.Data.Users) == 0 {
		c.logger.Error("user not found on twitch", slog.String("userName", msg.Channel))
		return
	}

	limiter := utils.CreateBotLimiter(isMod)
	c.BotClient.RateLimiters.Channels.Lock()
	c.BotClient.RateLimiters.Channels.Items[msg.Channel] = &types.Channel{
		IsMod:   isMod,
		Limiter: limiter,
		ID:      twitchReq.Data.Users[0].ID,
	}
	c.BotClient.RateLimiters.Channels.Unlock()

	go func() {
		if err := c.db.Model(&model.Channels{}).Where("id = ?", twitchReq.Data.Users[0].ID).Update(
			`"isBotMod"`,
			isMod,
		).Error; err != nil {
			c.logger.Error(
				"cannot update isMod",
				slog.Any("err", err),
				slog.Group(
					"channel",
					slog.String("id", twitchReq.Data.Users[0].ID),
					slog.String("login", twitchReq.Data.Users[0].Login),
				),
			)
		}
	}()
}
