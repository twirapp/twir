package types

import (
	"strings"
	"sync"

	"github.com/satont/go-helix/v2"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/twitch"

	ratelimiting "github.com/aidenwallis/go-ratelimiting/local"
	irc "github.com/gempir/go-twitch-irc/v3"
)

type Channel struct {
	IsMod   bool
	Limiter ratelimiting.SlidingWindow
}

type ChannelsMap struct {
	sync.Mutex
	Items map[string]*Channel
}

type RateLimiters struct {
	Global   ratelimiting.SlidingWindow
	Channels ChannelsMap
}

type BotClient struct {
	*irc.Client

	Api          *twitch.Twitch
	RateLimiters RateLimiters
	Model        *model.Bots
	TwitchUser   *helix.User
}

func (c *BotClient) SayWithRateLimiting(channel, text string, replyTo *string) {
	channelLimiter, ok := c.RateLimiters.Channels.Items[strings.ToLower(channel)]
	if !ok {
		return
	}

	if !c.RateLimiters.Global.TryTake() {
		return
	}

	// it should be separately
	if !channelLimiter.Limiter.TryTake() {
		return
	}

	text = validateResponseSlashes(text)

	if replyTo != nil {
		c.Reply(channel, *replyTo, text)
	} else {
		c.Say(channel, text)
	}
}

func validateResponseSlashes(response string) string {
	if strings.HasPrefix(response, "/me") || strings.HasPrefix(response, "/announce") {
		return response
	} else if strings.HasPrefix(response, "/") {
		return "Slash commands except /me and /announce is disallowed. This response wont be ever sended."
	} else if strings.HasPrefix(response, ".") {
		return `Message cannot start with "." symbol.`
	} else {
		return response
	}
}
