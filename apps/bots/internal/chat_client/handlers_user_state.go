package chat_client

import (
	"log/slog"

	irc "github.com/gempir/go-twitch-irc/v3"
	"github.com/nicklaw5/helix/v2"
	"github.com/satont/twir/apps/bots/pkg/utils"
	model "github.com/satont/twir/libs/gomodels"
)

func (c *ChatClient) OnUserStateMessage(msg irc.UserStateMessage) {
	moderatorBadge, _ := msg.User.Badges["moderator"]
	broadcasterBadge, _ := msg.User.Badges["broadcaster"]

	delete(c.RateLimiters.Channels.Items, msg.Channel)

	isMod := moderatorBadge == 1 || broadcasterBadge == 1

	twitchReq, _ := c.services.TwitchClient.GetUsers(
		&helix.UsersParams{
			Logins: []string{msg.Channel},
		},
	)
	if len(twitchReq.Data.Users) == 0 {
		c.services.Logger.Error("user not found on twitch", slog.String("userName", msg.Channel))
		return
	}

	limiter := utils.CreateBotLimiter(isMod)
	c.RateLimiters.Channels.Lock()
	c.RateLimiters.Channels.Items[msg.Channel] = &Channel{
		IsMod:   isMod,
		Limiter: limiter,
		ID:      twitchReq.Data.Users[0].ID,
	}
	c.RateLimiters.Channels.Unlock()

	go func() {
		if err := c.services.DB.Model(&model.Channels{}).Where(
			"id = ?",
			twitchReq.Data.Users[0].ID,
		).Update(
			`"isBotMod"`,
			isMod,
		).Error; err != nil {
			c.services.Logger.Error(
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
