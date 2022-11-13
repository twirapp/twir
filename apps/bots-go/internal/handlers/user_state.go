package handlers

import (
	"time"

	ratelimiting "github.com/aidenwallis/go-ratelimiting/local"
	irc "github.com/gempir/go-twitch-irc/v3"
)

func (c *Handlers) OnUserStateMessage(msg irc.UserStateMessage) {
	badge, botModBadge := msg.User.Badges["moderator"]
	if botModBadge {
		var limiter ratelimiting.SlidingWindow
		if badge == 1 {
			l, _ := ratelimiting.NewSlidingWindow(20, 30*time.Second)
			limiter = l
		} else {
			l, _ := ratelimiting.NewSlidingWindow(1, 2*time.Second)
			limiter = l
		}

		c.BotClient.RateLimiters.Channels[msg.Channel] = limiter
	}
}
