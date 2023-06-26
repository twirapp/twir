package handlers

import (
	irc "github.com/gempir/go-twitch-irc/v3"
	"github.com/satont/twir/apps/bots/pkg/utils"
)

func (c *Handlers) OnUserStateMessage(msg irc.UserStateMessage) {
	moderatorBadge, isModeratorBadge := msg.User.Badges["moderator"]
	broadcasterBadge, isBroadcasterBadge := msg.User.Badges["broadcaster"]
	if isModeratorBadge || isBroadcasterBadge {
		channel := c.BotClient.RateLimiters.Channels.Items[msg.Channel]

		c.BotClient.RateLimiters.Channels.Lock()
		defer c.BotClient.RateLimiters.Channels.Unlock()

		isMod := moderatorBadge == 1 || broadcasterBadge == 1

		if channel.IsMod && isMod {
			return
		}

		if !channel.IsMod && !isMod {
			return
		}

		channel.IsMod = isMod
		limiter := utils.CreateBotLimiter(isMod)

		c.BotClient.RateLimiters.Channels.Items[msg.Channel].Limiter = limiter
	}
}
