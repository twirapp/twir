package variables

import (
	"regexp"
	"tsuwari/parser/internal/config/twitch"
	types "tsuwari/parser/internal/types"
	emotes7tv "tsuwari/parser/internal/variables/emotes/7tv"
	emotesbttv "tsuwari/parser/internal/variables/emotes/bttv"
	emotesffz "tsuwari/parser/internal/variables/emotes/ffz"
	"tsuwari/parser/internal/variables/random"
	sender "tsuwari/parser/internal/variables/sender"
	streamcategory "tsuwari/parser/internal/variables/stream/category"
	streamtitle "tsuwari/parser/internal/variables/stream/title"
	streamuptime "tsuwari/parser/internal/variables/stream/uptime"
	streamviewers "tsuwari/parser/internal/variables/stream/viewers"
	variablescache "tsuwari/parser/internal/variablescache"

	"github.com/go-redis/redis/v9"
)

type Variables struct {
	Store  map[string]types.Variable
	Redis  *redis.Client
	Twitch *twitch.Twitch
}

var Regexp = regexp.MustCompile(`(?m)\$\((?P<all>(?P<main>[^.)|]+)(?:\.[^)|]+)?)(?:\|(?P<params>[^)]+))?\)`)

func New(redis *redis.Client, twitchApi *twitch.Twitch) Variables {
	ctx := Variables{
		Store:  make(map[string]types.Variable),
		Redis:  redis,
		Twitch: twitchApi,
	}

	ctx.Store[random.Name] = types.Variable{
		Name:    random.Name,
		Handler: random.Handler,
	}
	ctx.Store[sender.Name] = types.Variable{
		Name:    sender.Name,
		Handler: sender.Handler,
	}
	ctx.Store[streamuptime.Name] = types.Variable{
		Name:    streamuptime.Name,
		Handler: streamuptime.Handler,
	}
	ctx.Store[streamtitle.Name] = types.Variable{
		Name:    streamtitle.Name,
		Handler: streamtitle.Handler,
	}
	ctx.Store[streamviewers.Name] = types.Variable{
		Name:    streamviewers.Name,
		Handler: streamviewers.Handler,
	}
	ctx.Store[streamcategory.Name] = types.Variable{
		Name:    streamcategory.Name,
		Handler: streamcategory.Handler,
	}
	ctx.Store[emotesffz.Name] = types.Variable{
		Name:    emotesffz.Name,
		Handler: emotesffz.Handler,
	}
	ctx.Store[emotes7tv.Name] = types.Variable{
		Name:    emotes7tv.Name,
		Handler: emotes7tv.Handler,
	}
	ctx.Store[emotesbttv.Name] = types.Variable{
		Name:    emotesbttv.Name,
		Handler: emotesbttv.Handler,
	}

	return ctx
}

func (c Variables) ParseInput(cache *variablescache.VariablesCacheService, input string) string {
	result := Regexp.ReplaceAllStringFunc(input, func(s string) string {
		v := Regexp.FindStringSubmatch(s)
		all := v[1]
		// main := v[2]
		params := v[3]

		if val, ok := c.Store[all]; ok {
			res, err := val.Handler(cache, types.VariableHandlerParams{
				Key:    all,
				Params: &params,
			})

			if err != nil {
				return string(err.Error())
			} else {
				return res.Result
			}
		}

		return s
	})

	return result
}
