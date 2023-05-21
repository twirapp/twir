package variables

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/satont/tsuwari/apps/parser-new/internal/types/services"

	"github.com/satont/tsuwari/apps/parser-new/internal/types"
	"github.com/satont/tsuwari/apps/parser-new/internal/variables/command_param"
	commands_list "github.com/satont/tsuwari/apps/parser-new/internal/variables/commands"
	command_counters "github.com/satont/tsuwari/apps/parser-new/internal/variables/commands/counters"
	"github.com/satont/tsuwari/apps/parser-new/internal/variables/custom_var"
	"github.com/satont/tsuwari/libs/gopool"
)

type Variables struct {
	Store          map[string]*types.Variable
	services       *services.Services
	goroutinesPool *gopool.Pool
}

type Opts struct {
	Services *services.Services
}

var Regexp = regexp.MustCompile(
	`(?m)\$\((?P<all>(?P<main>[^.)|]+)(?:\.[^)|]+)?)(?:\|(?P<params>[^)]+))?\)`,
)

func New(opts *Opts) *Variables {
	store := make(map[string]*types.Variable)

	store[command_param.Variable.Name] = command_param.Variable
	store[commands_list.Variable.Name] = commands_list.Variable
	store[command_counters.CommandCounter.Name] = command_counters.CommandCounter
	store[command_counters.CommandFromOtherCounter.Name] = command_counters.CommandFromOtherCounter
	store[command_counters.CommandUserCounter.Name] = command_counters.CommandUserCounter
	store[custom_var.CustomVar.Name] = custom_var.CustomVar

	variables := &Variables{
		services:       opts.Services,
		goroutinesPool: gopool.NewPool(500),
	}

	return variables
}

func (c *Variables) ParseVariablesInText(ctx context.Context, parseCtx *types.ParseContext, input string) string {
	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}

	variablesParseCtx := &types.VariableParseContext{
		ParseContext: parseCtx,
	}

	for _, s := range Regexp.FindAllString(input, len(input)) {
		wg.Add(1)
		v := Regexp.FindStringSubmatch(s)
		all := v[1]
		params := v[3]

		variable, ok := c.Store[all]

		if !ok {
			wg.Done()
			continue
		}

		if variable.CommandsOnly && !parseCtx.IsCommand {
			mu.Lock()
			input = strings.ReplaceAll(input, s, fmt.Sprintf("$(%s)", all))
			mu.Unlock()
			wg.Done()
			continue
		}

		str := s
		c.goroutinesPool.Submit(func() {
			defer wg.Done()
			res, err := variable.Handler(ctx, variablesParseCtx, &types.VariableData{
				Key:    all,
				Params: &params,
			})

			if err == nil {
				mu.Lock()
				input = strings.ReplaceAll(input, str, res.Result)
				mu.Unlock()
			}
		})
	}

	wg.Wait()
	return input
}
