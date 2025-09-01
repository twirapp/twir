package variables

import (
	"context"
	"fmt"
	"regexp"
	"slices"
	"strings"
	"sync"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types/services"
	seventv "github.com/twirapp/twir/apps/parser/internal/variables/7tv"
	"github.com/twirapp/twir/apps/parser/internal/variables/channel"
	"github.com/twirapp/twir/apps/parser/internal/variables/chat_eval"
	"github.com/twirapp/twir/apps/parser/internal/variables/counttime"
	"github.com/twirapp/twir/apps/parser/internal/variables/donations/last_donate"
	"github.com/twirapp/twir/apps/parser/internal/variables/donations/top_donate"
	"github.com/twirapp/twir/apps/parser/internal/variables/donations/top_donate_stream"
	"github.com/twirapp/twir/apps/parser/internal/variables/followers"
	"github.com/twirapp/twir/apps/parser/internal/variables/keywords"
	"github.com/twirapp/twir/apps/parser/internal/variables/random"
	"github.com/twirapp/twir/apps/parser/internal/variables/repeat"
	"github.com/twirapp/twir/apps/parser/internal/variables/request"
	"github.com/twirapp/twir/apps/parser/internal/variables/sender"
	"github.com/twirapp/twir/apps/parser/internal/variables/shorturl"
	"github.com/twirapp/twir/apps/parser/internal/variables/song"
	"github.com/twirapp/twir/apps/parser/internal/variables/stream"
	"github.com/twirapp/twir/apps/parser/internal/variables/subscribers"
	"github.com/twirapp/twir/apps/parser/internal/variables/to_user"
	"github.com/twirapp/twir/apps/parser/internal/variables/top"
	"github.com/twirapp/twir/apps/parser/internal/variables/tracer"
	"github.com/twirapp/twir/apps/parser/internal/variables/user"
	"github.com/twirapp/twir/apps/parser/internal/variables/valorant"
	"github.com/twirapp/twir/apps/parser/internal/variables/weather"

	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/internal/variables/command_param"
	commands_list "github.com/twirapp/twir/apps/parser/internal/variables/commands"
	command_counters "github.com/twirapp/twir/apps/parser/internal/variables/commands/counters"
	"github.com/twirapp/twir/apps/parser/internal/variables/custom_var"
	"github.com/twirapp/twir/apps/parser/internal/variables/emotes"
	"github.com/twirapp/twir/apps/parser/internal/variables/faceit"
)

type Variables struct {
	Store    map[string]*types.Variable
	services *services.Services
}

type Opts struct {
	Services *services.Services
}

var Regexp = regexp.MustCompile(
	`(?m)\$\((?P<all>(?P<main>[^.)|]+)(?:\.[^)|]+)?)(?:\|(?P<params>[^)]+))?\)`,
)

func New(opts *Opts) *Variables {
	store := lo.SliceToMap(
		[]*types.Variable{
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
			faceit.Gain,
			faceit.Lose,
			keywords.Counter,
			random.Number,
			random.OnlineUser,
			random.Phrase,
			sender.Sender,
			sender.SenderDisplayName,
			sender.ID,
			song.CurrentSong,
			song.Cover,
			song.History,
			song.HistorySpotify,
			song.HistoryLastfm,
			stream.Category,
			stream.Title,
			stream.Uptime,
			stream.Viewers,
			stream.CategoryTime,
			stream.Followers,
			// stream.CategoriesTime,
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
			user.Commands,
			user.ID,
			last_donate.Amount,
			last_donate.Currency,
			last_donate.UserName,
			top_donate.UserName,
			top_donate.Amount,
			top_donate.Currency,
			top_donate_stream.Amount,
			top_donate_stream.Currency,
			top_donate_stream.UserName,
			followers.LatestFollowerUsername,
			followers.Count,
			subscribers.Count,
			subscribers.LatestSubscriberUsername,
			request.Request,
			chat_eval.ChatEval,
			user.Reputation,
			weather.Weather,
			channel.Name,
			seventv.ProfileLink,
			seventv.EmoteSetLink,
			seventv.EmoteSetName,
			seventv.EmoteSetCount,
			seventv.Roles,
			seventv.EmoteSetCapacity,
			seventv.Paint,
			seventv.EditorForCount,
			seventv.ProfileCreatedAt,
			seventv.UnlockedPaints,
			repeat.Variable,
			shorturl.Variable,
			counttime.CountDown,
			counttime.CountUp,
		}, func(v *types.Variable) (string, *types.Variable) {
			return v.Name, v
		},
	)

	variables := &Variables{
		services: opts.Services,
		Store:    store,
	}

	return variables
}

func (c *Variables) ParseVariablesInText(
	ctx context.Context,
	parseCtx *types.ParseContext,
	input string,
) string {
	newCtx, span := tracer.VariablesTracer.Start(ctx, "ParseVariablesInText")
	defer span.End()
	ctx = newCtx

	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}

	variablesParseCtx := &types.VariableParseContext{
		ParseContext: parseCtx,
	}

	variablesCopy := make([]*types.Variable, 0, len(c.Store))
	for _, v := range c.Store {
		variablesCopy = append(variablesCopy, v)
	}

	slices.SortFunc(
		variablesCopy, func(a, b *types.Variable) int {
			return b.Priority - a.Priority
		},
	)

	if parseCtx.Text != nil && len(*parseCtx.Text) > 0 {
		input = strings.ReplaceAll(input, "$(command.param)", *parseCtx.Text)
	}

	for _, s := range Regexp.FindAllString(input, len(input)) {
		v := Regexp.FindStringSubmatch(s)
		all := v[1]
		params := v[3]

		variable, ok := c.Store[all]

		if !ok {
			continue
		}

		if variable.CommandsOnly && !parseCtx.IsCommand {
			// mu.Lock()
			// input = strings.ReplaceAll(input, s, fmt.Sprintf("$(%s)", all))
			// mu.Unlock()
			continue
		}

		if variable.DisableInCustomVariables && parseCtx.IsInCustomVar {
			// mu.Lock()
			// input = strings.ReplaceAll(input, s, fmt.Sprintf("$(%s)", all))
			// mu.Unlock()
			continue
		}

		str := s

		wg.Add(1)

		variableCtx, variableSpan := tracer.VariablesTracer.Start(
			ctx,
			fmt.Sprintf("VariableHandler.%s", all),
		)

		go func() {
			defer wg.Done()
			defer variableSpan.End()

			res, err := variable.Handler(
				variableCtx,
				variablesParseCtx,
				&types.VariableData{
					Key:    all,
					Params: &params,
				},
			)

			if res != nil && err == nil {
				mu.Lock()

				if variable.NotCachable {
					input = strings.Replace(input, str, res.Result, 1)
				} else {
					input = strings.ReplaceAll(input, str, res.Result)
				}

				mu.Unlock()
			}
		}()
	}

	wg.Wait()
	return input
}
