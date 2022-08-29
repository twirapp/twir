package variables

import (
	"regexp"
	"tsuwari/parser/internal/types"
	"tsuwari/parser/internal/variables/random"
	"tsuwari/parser/internal/variables/sender"

	"github.com/go-redis/redis/v9"
)

type Variables struct {
	store  map[string]types.Variable
	regexp *regexp.Regexp
	redis  *redis.Client
}

func New(redis *redis.Client) Variables {
	ctx := Variables{
		store:  make(map[string]types.Variable),
		regexp: regexp.MustCompile(`\$\(([^)|]+)(?:\|([^)]+))?\)`),
		redis:  redis,
	}

	ctx.load()

	return ctx
}

func (c Variables) load() {
	c.store[random.Name] = types.Variable{
		Name:    random.Name,
		Handler: random.Handler,
	}
	c.store[sender.Name] = types.Variable{
		Name:    sender.Name,
		Handler: sender.Handler,
	}
}

func (c Variables) ParseInput(input string) string {
	result := c.regexp.ReplaceAllStringFunc(input, func(s string) string {
		v := c.regexp.FindStringSubmatchIndex(s)
		matchedVarName := s[v[2]:v[3]]

		var params *string

		if v[4] != -1 {
			p := s[v[4]:v[5]]
			params = &p
		}

		if val, ok := c.store[matchedVarName]; ok {
			res, err := val.Handler(types.VariableHandlerParams{
				Key:    matchedVarName,
				Params: params,
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
