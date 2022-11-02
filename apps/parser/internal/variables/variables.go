package variables

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
	types "tsuwari/parser/internal/types"
	commandsvariable "tsuwari/parser/internal/variables/commands"
	command_counter "tsuwari/parser/internal/variables/commands/counter"
	customvar "tsuwari/parser/internal/variables/customvar"
	emotes7tv "tsuwari/parser/internal/variables/emotes/7tv"
	emotesbttv "tsuwari/parser/internal/variables/emotes/bttv"
	emotesffz "tsuwari/parser/internal/variables/emotes/ffz"
	faceitelo "tsuwari/parser/internal/variables/faceit/elo"
	faceitelodiff "tsuwari/parser/internal/variables/faceit/elodiff"
	faceitlvl "tsuwari/parser/internal/variables/faceit/lvl"
	random "tsuwari/parser/internal/variables/random/number"
	randomonlineuser "tsuwari/parser/internal/variables/random/online/user"
	phrase "tsuwari/parser/internal/variables/random/phrase"
	sender "tsuwari/parser/internal/variables/sender"
	song "tsuwari/parser/internal/variables/song"
	streamcategory "tsuwari/parser/internal/variables/stream/category"
	streammessages "tsuwari/parser/internal/variables/stream/messages"
	streamtitle "tsuwari/parser/internal/variables/stream/title"
	streamuptime "tsuwari/parser/internal/variables/stream/uptime"
	streamviewers "tsuwari/parser/internal/variables/stream/viewers"
	topMessages "tsuwari/parser/internal/variables/top/messages"
	"tsuwari/parser/internal/variables/touser"
	userage "tsuwari/parser/internal/variables/user/age"
	userfollowage "tsuwari/parser/internal/variables/user/followage"
	usermessages "tsuwari/parser/internal/variables/user/messages"
	variables_cache "tsuwari/parser/internal/variablescache"

	"github.com/samber/lo"
)

type Variables struct {
	Store []types.Variable
}

var Regexp = regexp.MustCompile(
	`(?m)\$\((?P<all>(?P<main>[^.)|]+)(?:\.[^)|]+)?)(?:\|(?P<params>[^)]+))?\)`,
)

func New() Variables {
	store := []types.Variable{
		commandsvariable.Variable,
		customvar.Variable,
		emotes7tv.Variable,
		emotesbttv.Variable,
		emotesffz.Variable,
		faceitelo.Variable,
		faceitelodiff.Variable,
		faceitlvl.Variable,
		random.Variable,
		randomonlineuser.Variable,
		sender.Variable,
		song.Variable,
		streamcategory.Variable,
		streammessages.Variable,
		streamtitle.Variable,
		streamuptime.Variable,
		streamviewers.Variable,
		topMessages.Variable,
		userage.Variable,
		userfollowage.Variable,
		usermessages.Variable,
		touser.Variable,
		phrase.Variable,
		command_counter.CommandVariable,
		command_counter.UserVariable,
	}

	ctx := Variables{
		Store: store,
	}

	return ctx
}

func (c Variables) ParseInput(cache *variables_cache.VariablesCacheService, input string) string {
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
			defer wg.Done()
			continue
		}

		if variable.CommandsOnly != nil && *variable.CommandsOnly && !cache.IsCommand {
			mu.Lock()
			input = strings.ReplaceAll(input, s, fmt.Sprintf("$(%s)", all))
			mu.Unlock()
			wg.Done()
			continue
		}

		go func(s string) {
			defer wg.Done()
			res, err := variable.Handler(cache, types.VariableHandlerParams{
				Key:    all,
				Params: &params,
			})

			if err == nil {
				mu.Lock()
				input = strings.ReplaceAll(input, s, res.Result)
				mu.Unlock()
			}
		}(s)
	}

	wg.Wait()
	return input
}
