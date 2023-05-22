package variables

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/parser-new/internal/types/services"
	"github.com/satont/tsuwari/apps/parser-new/internal/variables/keywords"
	"github.com/satont/tsuwari/apps/parser-new/internal/variables/random"
	"github.com/satont/tsuwari/apps/parser-new/internal/variables/sender"
	"github.com/satont/tsuwari/apps/parser-new/internal/variables/song"
	"github.com/satont/tsuwari/apps/parser-new/internal/variables/stream"
	"github.com/satont/tsuwari/apps/parser-new/internal/variables/top"

	"github.com/satont/tsuwari/apps/parser-new/internal/types"
	"github.com/satont/tsuwari/apps/parser-new/internal/variables/command_param"
	commands_list "github.com/satont/tsuwari/apps/parser-new/internal/variables/commands"
	command_counters "github.com/satont/tsuwari/apps/parser-new/internal/variables/commands/counters"
	"github.com/satont/tsuwari/apps/parser-new/internal/variables/custom_var"
	"github.com/satont/tsuwari/apps/parser-new/internal/variables/emotes"
	"github.com/satont/tsuwari/apps/parser-new/internal/variables/faceit"
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
	store := lo.SliceToMap([]*types.Variable{
		command_param.Variable,
		commands_list.Variable,
		command_counters.CommandCounter,
		command_counters.CommandFromOtherCounter,
		command_counters.CommandUserCounter,
		custom_var.CustomVar,
		emotes.SevenTv,
		emotes.BetterTTV,
		emotes.FrankerFaceZ,
		faceit.Elo,
		faceit.EloDiff,
		faceit.LVL,
		faceit.ScoreLoses,
		faceit.ScoreWins,
		faceit.TrendExtended,
		faceit.TrendSimple,
		keywords.Counter,
		random.Number,
		random.OnlineUser,
		random.Phrase,
		sender.Sender,
		song.Song,
		stream.Category,
		stream.Title,
		stream.Uptime,
		stream.Viewers,
		top.Watched,
		top.Messages,
		top.SongRequesters,
		top.ChannelPoints,
		top.Emotes,
		top.EmotesUsers,
	}, func(v *types.Variable) (string, *types.Variable) {
		return v.Name, v
	})

	variables := &Variables{
		services:       opts.Services,
		goroutinesPool: gopool.NewPool(500),
		Store:          store,
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
