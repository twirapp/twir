package variables

import (
	"fmt"
	"github.com/satont/tsuwari/apps/parser/internal/variables/command_param"
	top_channel_points "github.com/satont/tsuwari/apps/parser/internal/variables/top/channel_points"
	top_song_requesters "github.com/satont/tsuwari/apps/parser/internal/variables/top/song_requesters"
	user_emotes "github.com/satont/tsuwari/apps/parser/internal/variables/user/emotes"
	user_songs_requested "github.com/satont/tsuwari/apps/parser/internal/variables/user/songs_requested"
	user_top "github.com/satont/tsuwari/apps/parser/internal/variables/user/top"
	user_used_channel_points "github.com/satont/tsuwari/apps/parser/internal/variables/user/used_channel_points"
	valorant_matches "github.com/satont/tsuwari/apps/parser/internal/variables/valorant"
	valorant_profile "github.com/satont/tsuwari/apps/parser/internal/variables/valorant/profile"
	"github.com/satont/tsuwari/libs/gopool"
	"regexp"
	"strings"
	"sync"

	"github.com/satont/tsuwari/apps/parser/internal/types"
	commandsvariable "github.com/satont/tsuwari/apps/parser/internal/variables/commands"
	command_counter "github.com/satont/tsuwari/apps/parser/internal/variables/commands/counter"
	"github.com/satont/tsuwari/apps/parser/internal/variables/customvar"
	emotes7tv "github.com/satont/tsuwari/apps/parser/internal/variables/emotes/7tv"
	emotesbttv "github.com/satont/tsuwari/apps/parser/internal/variables/emotes/bttv"
	emotesffz "github.com/satont/tsuwari/apps/parser/internal/variables/emotes/ffz"
	faceit_elo "github.com/satont/tsuwari/apps/parser/internal/variables/faceit/elo"
	faceit_elo_diff "github.com/satont/tsuwari/apps/parser/internal/variables/faceit/elodiff"
	faceit_lvl "github.com/satont/tsuwari/apps/parser/internal/variables/faceit/lvl"
	faceit_score "github.com/satont/tsuwari/apps/parser/internal/variables/faceit/score"
	faceit_trand "github.com/satont/tsuwari/apps/parser/internal/variables/faceit/trend"
	"github.com/satont/tsuwari/apps/parser/internal/variables/keywords"
	random "github.com/satont/tsuwari/apps/parser/internal/variables/random/number"
	randomonlineuser "github.com/satont/tsuwari/apps/parser/internal/variables/random/online/user"
	phrase "github.com/satont/tsuwari/apps/parser/internal/variables/random/phrase"
	"github.com/satont/tsuwari/apps/parser/internal/variables/sender"
	"github.com/satont/tsuwari/apps/parser/internal/variables/song"
	streamcategory "github.com/satont/tsuwari/apps/parser/internal/variables/stream/category"
	streammessages "github.com/satont/tsuwari/apps/parser/internal/variables/stream/messages"
	streamtitle "github.com/satont/tsuwari/apps/parser/internal/variables/stream/title"
	streamuptime "github.com/satont/tsuwari/apps/parser/internal/variables/stream/uptime"
	streamviewers "github.com/satont/tsuwari/apps/parser/internal/variables/stream/viewers"
	top_emotes "github.com/satont/tsuwari/apps/parser/internal/variables/top/emotes"
	top_messages "github.com/satont/tsuwari/apps/parser/internal/variables/top/messages"
	top_watched "github.com/satont/tsuwari/apps/parser/internal/variables/top/watched"
	"github.com/satont/tsuwari/apps/parser/internal/variables/touser"
	userage "github.com/satont/tsuwari/apps/parser/internal/variables/user/age"
	user_follow "github.com/satont/tsuwari/apps/parser/internal/variables/user/follow"
	usermessages "github.com/satont/tsuwari/apps/parser/internal/variables/user/messages"
	userwatched "github.com/satont/tsuwari/apps/parser/internal/variables/user/watched"
	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"

	"github.com/samber/lo"
)

type Variables struct {
	Store          []types.Variable
	goroutinesPool *gopool.Pool
}

var Regexp = regexp.MustCompile(
	`(?m)\$\((?P<all>(?P<main>[^.)|]+)(?:\.[^)|]+)?)(?:\|(?P<params>[^)]+))?\)`,
)

func New() *Variables {
	store := []types.Variable{
		commandsvariable.Variable,
		customvar.Variable,
		emotes7tv.Variable,
		emotesbttv.Variable,
		emotesffz.Variable,
		faceit_elo.Variable,
		faceit_elo_diff.Variable,
		faceit_lvl.Variable,
		faceit_trand.SimpleTrend,
		faceit_trand.ExtendedTrend,
		faceit_score.Wins,
		faceit_score.Loses,
		random.Variable,
		randomonlineuser.Variable,
		sender.Variable,
		song.Variable,
		streamcategory.Variable,
		streammessages.Variable,
		streamtitle.Variable,
		streamuptime.Variable,
		streamviewers.Variable,
		top_messages.Variable,
		top_watched.Variable,
		top_emotes.Variable,
		userage.Variable,
		user_follow.FollowageVariable,
		user_follow.FollowsinceVariable,
		usermessages.Variable,
		userwatched.Variable,
		touser.Variable,
		phrase.Variable,
		user_top.TopEmotesVariable,
		command_counter.CommandVariable,
		command_counter.UserVariable,
		command_counter.CommandVariableFromOther,
		keywords.Counter,
		user_emotes.Variable,
		command_param.Variable,
		top_emotes.UsersVariable,
		user_used_channel_points.Variable,
		top_channel_points.Variable,
		top_song_requesters.Variable,
		user_songs_requested.CountVariable,
		user_songs_requested.DurationVariable,
		valorant_profile.Tier,
		valorant_profile.TierPatched,
		valorant_profile.RankingInTier,
		valorant_profile.Elo,
		valorant_profile.MmrChangeToLastGame,
		valorant_matches.Trend,
	}

	ctx := &Variables{
		Store:          store,
		goroutinesPool: gopool.NewPool(100),
	}

	return ctx
}

func (c *Variables) ParseInput(cache *variables_cache.VariablesCacheService, input string) string {
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}

	for _, s := range Regexp.FindAllString(input, len(input)) {
		wg.Add(1)
		v := Regexp.FindStringSubmatch(s)
		all := v[1]
		params := v[3]

		variable, ok := lo.Find(c.Store, func(el types.Variable) bool {
			return el.Name == all
		})

		if !ok {
			wg.Done()
			continue
		}

		if variable.CommandsOnly != nil && *variable.CommandsOnly && !cache.IsCommand {
			mu.Lock()
			input = strings.ReplaceAll(input, s, fmt.Sprintf("$(%s)", all))
			mu.Unlock()
			wg.Done()
			continue
		}

		str := s
		c.goroutinesPool.Submit(func() {
			defer wg.Done()
			res, err := variable.Handler(cache, types.VariableHandlerParams{
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
