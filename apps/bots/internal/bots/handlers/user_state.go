package handlers

import (
	irc "github.com/gempir/go-twitch-irc/v3"
	"github.com/satont/tsuwari/apps/bots/pkg/utils"
)

func (c *Handlers) OnUserStateMessage(msg irc.UserStateMessage) {
	badge, botModBadge := msg.User.Badges["moderator"]
	if botModBadge {
		channel := c.BotClient.RateLimiters.Channels.Items[msg.Channel]

		c.BotClient.RateLimiters.Channels.Lock()
		defer c.BotClient.RateLimiters.Channels.Unlock()

		isMod := badge == 1

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
