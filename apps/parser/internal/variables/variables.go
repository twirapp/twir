package variables

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/parser/internal/types/services"
	"github.com/satont/tsuwari/apps/parser/internal/variables/keywords"
	"github.com/satont/tsuwari/apps/parser/internal/variables/random"
	"github.com/satont/tsuwari/apps/parser/internal/variables/rocket_league"
	"github.com/satont/tsuwari/apps/parser/internal/variables/sender"
	"github.com/satont/tsuwari/apps/parser/internal/variables/song"
	"github.com/satont/tsuwari/apps/parser/internal/variables/stream"
	"github.com/satont/tsuwari/apps/parser/internal/variables/to_user"
	"github.com/satont/tsuwari/apps/parser/internal/variables/top"
	"github.com/satont/tsuwari/apps/parser/internal/variables/user"
	"github.com/satont/tsuwari/apps/parser/internal/variables/valorant"

	"github.com/satont/tsuwari/apps/parser/internal/types"
	"github.com/satont/tsuwari/apps/parser/internal/variables/command_param"
	commands_list "github.com/satont/tsuwari/apps/parser/internal/variables/commands"
	command_counters "github.com/satont/tsuwari/apps/parser/internal/variables/commands/counters"
	"github.com/satont/tsuwari/apps/parser/internal/variables/custom_var"
	"github.com/satont/tsuwari/apps/parser/internal/variables/emotes"
	"github.com/satont/tsuwari/apps/parser/internal/variables/faceit"
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
		to_user.ToUser,
		valorant.Matches,
		valorant.Elo,
		valorant.Tier,
		valorant.TierText,
		valorant.RankInTier,
		user.Age,
		user.FollowAge,
		user.FollowSince,
		user.Emotes,
		user.Messages,
		user.UsedChannelPoints,
		user.SongsRequested,
		user.SongsRequestedDuration,
		user.EmotesTop,
		user.Watched,
		rocket_league.Rank2v2,
		rocket_league.Rank1v1,
		rocket_league.Rank3v3,
		rocket_league.RankDropshot,
		rocket_league.RankRumble,
		rocket_league.Rating2v2,
		rocket_league.Rating1v1,
		rocket_league.Rating3v3,
		rocket_league.RatingDropshot,
		rocket_league.RatingRumble,
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
			res, err := variable.Handler(ctx, variablesParseCtx, &types.VariableData{
				Key:    all,
				Params: &params,
			})

			if err == nil {
				mu.Lock()
				input = strings.ReplaceAll(input, str, res.Result)
				mu.Unlock()
			}
			wg.Done()
		})
	}

	wg.Wait()
	return input
}
