package variables

import (
	"regexp"
	types "tsuwari/parser/internal/types"
	"tsuwari/parser/internal/variables/random"
	sender "tsuwari/parser/internal/variables/sender"
	"tsuwari/parser/internal/variables/stream/streamId"
	variablescache "tsuwari/parser/internal/variablescache"

	"github.com/go-redis/redis/v9"
)

type Variables struct {
	Store map[string]types.Variable
	Redis *redis.Client
}

var Regexp = regexp.MustCompile(`(?m)\$\((?P<all>(?P<main>[^.)|]+)(?:\.[^)|]+)?)(?:\|(?P<params>[^)]+))?\)`)

func New(redis *redis.Client) Variables {
	ctx := Variables{
		Store: make(map[string]types.Variable),
		Redis: redis,
	}

	ctx.Store[random.Name] = types.Variable{
		Name:    random.Name,
		Handler: random.Handler,
	}
	ctx.Store[sender.Name] = types.Variable{
		Name:    sender.Name,
		Handler: sender.Handler,
	}
	ctx.Store[streamId.Name] = types.Variable{
		Name:    streamId.Name,
		Handler: streamId.Handler,
	}

	return ctx
}

func (c Variables) ParseInput(cache *variablescache.VariablesCacheService, input string) string {
	result := Regexp.ReplaceAllStringFunc(input, func(s string) string {
		v := Regexp.FindStringSubmatch(s)
		// main := v[1]
		all := v[2]
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
