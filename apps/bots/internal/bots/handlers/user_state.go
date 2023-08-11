package handlers

import (
	irc "github.com/gempir/go-twitch-irc/v3"
	"github.com/nicklaw5/helix/v2"
	"github.com/satont/twir/apps/bots/pkg/utils"
	"github.com/satont/twir/apps/bots/types"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/twitch"
	"go.uber.org/zap"
)

func (c *Handlers) OnUserStateMessage(msg irc.UserStateMessage) {
	moderatorBadge, _ := msg.User.Badges["moderator"]
	broadcasterBadge, _ := msg.User.Badges["broadcaster"]

	twitchClient, err := twitch.NewBotClient(c.BotClient.Model.ID, *c.cfg, c.tokensGrpc)
	if err != nil {
		panic(err)
	}

	delete(c.BotClient.RateLimiters.Channels.Items, msg.Channel)

	isMod := moderatorBadge == 1 || broadcasterBadge == 1

	twitchReq, err := twitchClient.GetUsers(&helix.UsersParams{Logins: []string{msg.Channel}})
	if len(twitchReq.Data.Users) == 0 {
		zap.S().Error("user not found on twitch", zap.String("userName", msg.Channel))
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
			zap.S().Error(err)
		}
	}()
}
