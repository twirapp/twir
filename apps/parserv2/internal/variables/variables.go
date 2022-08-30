package variables

import (
	"regexp"
	types "tsuwari/parser/internal/types"
	"tsuwari/parser/internal/variables/random"
	sender "tsuwari/parser/internal/variables/sender"
	stream "tsuwari/parser/internal/variables/stream"
	variablescache "tsuwari/parser/internal/variablescache"

	"github.com/go-redis/redis/v9"
)

type Variables struct {
	Store  map[string]types.Variable
	Regexp *regexp.Regexp
	Redis  *redis.Client
}

func New(redis *redis.Client) Variables {
	ctx := Variables{
		Store:  make(map[string]types.Variable),
		Regexp: regexp.MustCompile(`(?m)\$\((?P<all>(?P<main>[^.)|]+)(?:\.[^)|]+)?)(?:\|(?P<params>[^)]+))?\)`),
		Redis:  redis,
	}

	ctx.Store[random.Name] = types.Variable{
		Name:    random.Name,
		Handler: random.Handler,
	}
	ctx.Store[sender.Name] = types.Variable{
		Name:    sender.Name,
		Handler: sender.Handler,
	}
	ctx.Store[stream.Name] = types.Variable{
		Name:    stream.Name,
		Handler: stream.Handler,
	}

	return ctx
}

func (c Variables) ParseInput(cache *variablescache.VariablesCacheService, input string) string {
	result := c.Regexp.ReplaceAllStringFunc(input, func(s string) string {
		v := c.Regexp.FindStringSubmatch(s)
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
