package types

import (
	"context"
	"strings"

	"github.com/satont/tsuwari/libs/twitch"

	ratelimiting "github.com/aidenwallis/go-ratelimiting/local"
	irc "github.com/gempir/go-twitch-irc/v3"
)

type RateLimiters struct {
	Global   ratelimiting.SlidingWindow
	Channels map[string]ratelimiting.SlidingWindow
}

type BotClient struct {
	*irc.Client

	Api          *twitch.Twitch
	RateLimiters RateLimiters
}

func (c *BotClient) SayWithRateLimiting(channel, text string, replyTo *string) {
	channelLimiter, ok := c.RateLimiters.Channels[strings.ToLower(channel)]
	if !ok {
		return
	}

	ctx := context.Background()

	c.RateLimiters.Global.WaitFunc(ctx, func() {
		channelLimiter.WaitFunc(ctx, func() {
			if replyTo != nil {
				c.Reply(channel, *replyTo, text)
			} else {
				c.Say(channel, text)
			}
		})
	})
}
