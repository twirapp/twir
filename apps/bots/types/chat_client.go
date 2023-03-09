package types

import (
	"strings"
	"sync"
	"unicode/utf8"

	ratelimiting "github.com/aidenwallis/go-ratelimiting/local"
	irc "github.com/gempir/go-twitch-irc/v3"
	"github.com/satont/go-helix/v2"
	model "github.com/satont/tsuwari/libs/gomodels"
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

	text = strings.ReplaceAll(text, "\n", " ")

	parts := splitTextByLength(text)

	if replyTo != nil {
		for _, part := range parts {
			text = validateResponseSlashes(text)
			c.Reply(channel, *replyTo, part)
		}
	} else {
		for _, part := range parts {
			text = validateResponseSlashes(text)
			c.Say(channel, part)
		}
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

func splitTextByLength(text string) []string {
	var parts []string
	for len(text) > 0 {
		if len(text) < 500 {
			parts = append(parts, text)
			break
		}
		// Find the last complete UTF-8 character within the next 500 bytes
		i := 500
		for !utf8.ValidString(text[:i]) {
			i--
		}
		parts = append(parts, text[:i])
		text = text[i:]
	}

	return parts
}