package handlers

import (
	"time"

	ratelimiting "github.com/aidenwallis/go-ratelimiting/local"
	irc "github.com/gempir/go-twitch-irc/v3"
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

		var limiter ratelimiting.SlidingWindow

		if isMod {
			l, _ := ratelimiting.NewSlidingWindow(20, 30*time.Second)
			limiter = l
		} else {
			l, _ := ratelimiting.NewSlidingWindow(1, 2*time.Second)
			limiter = l
		}

		c.BotClient.RateLimiters.Channels.Items[msg.Channel].Limiter = limiter
	}
}
