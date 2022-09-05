package variables

import (
	"regexp"
	"strings"
	"sync"
	"tsuwari/parser/internal/config/twitch"
	types "tsuwari/parser/internal/types"
	commandsvariable "tsuwari/parser/internal/variables/commands"
	customvar "tsuwari/parser/internal/variables/customvar"
	emotes7tv "tsuwari/parser/internal/variables/emotes/7tv"
	emotesbttv "tsuwari/parser/internal/variables/emotes/bttv"
	emotesffz "tsuwari/parser/internal/variables/emotes/ffz"
	faceitelo "tsuwari/parser/internal/variables/faceit/elo"
	faceitelodiff "tsuwari/parser/internal/variables/faceit/elodiff"
	faceitlvl "tsuwari/parser/internal/variables/faceit/lvl"
	"tsuwari/parser/internal/variables/random"
	sender "tsuwari/parser/internal/variables/sender"
	song "tsuwari/parser/internal/variables/song"
	streamcategory "tsuwari/parser/internal/variables/stream/category"
	streammessages "tsuwari/parser/internal/variables/stream/messages"
	streamtitle "tsuwari/parser/internal/variables/stream/title"
	streamuptime "tsuwari/parser/internal/variables/stream/uptime"
	streamviewers "tsuwari/parser/internal/variables/stream/viewers"
	topMessages "tsuwari/parser/internal/variables/top/messages"
	userage "tsuwari/parser/internal/variables/user/age"
	userfollowage "tsuwari/parser/internal/variables/user/followage"
	usermessages "tsuwari/parser/internal/variables/user/messages"
	variablescache "tsuwari/parser/internal/variablescache"

	"github.com/go-redis/redis/v9"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type Variables struct {
	Store  map[string]types.Variable
	Redis  *redis.Client
	Twitch *twitch.Twitch
	Db     *gorm.DB
}

var Regexp = regexp.MustCompile(`(?m)\$\((?P<all>(?P<main>[^.)|]+)(?:\.[^)|]+)?)(?:\|(?P<params>[^)]+))?\)`)

func New(redis *redis.Client, twitchApi *twitch.Twitch, db *gorm.DB) Variables {
	ctx := Variables{
		Store:  make(map[string]types.Variable),
		Redis:  redis,
		Twitch: twitchApi,
		Db:     db,
	}

	ctx.Store[random.Name] = types.Variable{
		Name:        random.Name,
		Handler:     random.Handler,
		Description: lo.ToPtr(random.Description),
	}
	ctx.Store[sender.Name] = types.Variable{
		Name:        sender.Name,
		Handler:     sender.Handler,
		Description: lo.ToPtr(sender.Description),
	}
	ctx.Store[streamuptime.Name] = types.Variable{
		Name:        streamuptime.Name,
		Handler:     streamuptime.Handler,
		Description: lo.ToPtr(streamuptime.Description),
	}
	ctx.Store[streamtitle.Name] = types.Variable{
		Name:        streamtitle.Name,
		Handler:     streamtitle.Handler,
		Description: lo.ToPtr(streamtitle.Description),
	}
	ctx.Store[streamviewers.Name] = types.Variable{
		Name:        streamviewers.Name,
		Handler:     streamviewers.Handler,
		Description: lo.ToPtr(streamviewers.Description),
	}
	ctx.Store[streamcategory.Name] = types.Variable{
		Name:        streamcategory.Name,
		Handler:     streamcategory.Handler,
		Description: lo.ToPtr(streamcategory.Description),
	}
	ctx.Store[streammessages.Name] = types.Variable{
		Name:        streammessages.Name,
		Handler:     streammessages.Handler,
		Description: lo.ToPtr(streammessages.Description),
	}
	ctx.Store[emotesffz.Name] = types.Variable{
		Name:        emotesffz.Name,
		Handler:     emotesffz.Handler,
		Description: lo.ToPtr(emotesffz.Description),
	}
	ctx.Store[emotes7tv.Name] = types.Variable{
		Name:        emotes7tv.Name,
		Handler:     emotes7tv.Handler,
		Description: lo.ToPtr(emotes7tv.Description),
	}
	ctx.Store[emotesbttv.Name] = types.Variable{
		Name:        emotesbttv.Name,
		Handler:     emotesbttv.Handler,
		Description: lo.ToPtr(emotesbttv.Description),
	}
	ctx.Store[usermessages.Name] = types.Variable{
		Name:        usermessages.Name,
		Handler:     usermessages.Handler,
		Description: lo.ToPtr(usermessages.Description),
	}
	ctx.Store[userfollowage.Name] = types.Variable{
		Name:        userfollowage.Name,
		Handler:     userfollowage.Handler,
		Description: lo.ToPtr(userfollowage.Description),
	}
	ctx.Store[userage.Name] = types.Variable{
		Name:        userage.Name,
		Handler:     userage.Handler,
		Description: lo.ToPtr(userage.Description),
	}
	ctx.Store[faceitelo.Name] = types.Variable{
		Name:        faceitelo.Name,
		Handler:     faceitelo.Handler,
		Description: lo.ToPtr(faceitelo.Description),
	}
	ctx.Store[faceitlvl.Name] = types.Variable{
		Name:        faceitlvl.Name,
		Handler:     faceitlvl.Handler,
		Description: lo.ToPtr(faceitelo.Description),
	}
	ctx.Store[faceitelodiff.Name] = types.Variable{
		Name:        faceitelodiff.Name,
		Handler:     faceitelodiff.Handler,
		Description: lo.ToPtr(faceitelodiff.Description),
	}
	ctx.Store[customvar.Name] = types.Variable{
		Name:        customvar.Name,
		Handler:     customvar.Handler,
		Description: lo.ToPtr(customvar.Description),
	}
	ctx.Store[song.Name] = types.Variable{
		Name:        song.Name,
		Handler:     song.Handler,
		Description: lo.ToPtr(song.Description),
	}
	ctx.Store[topMessages.Name] = types.Variable{
		Name:        topMessages.Name,
		Handler:     topMessages.Handler,
		Description: lo.ToPtr(topMessages.Description),
	}
	ctx.Store[commandsvariable.Name] = types.Variable{
		Name:        commandsvariable.Name,
		Handler:     commandsvariable.Handler,
		Description: lo.ToPtr(commandsvariable.Description),
	}

	return ctx
}

func (c Variables) ParseInput(cache *variablescache.VariablesCacheService, input string) string {
	wg := sync.WaitGroup{}

	for _, s := range Regexp.FindAllString(input, len(input)) {
		wg.Add(1)
		v := Regexp.FindStringSubmatch(s)
		all := v[1]
		params := v[3]

		if val, ok := c.Store[all]; ok {
			go func(s string) {
				defer wg.Done()
				res, err := val.Handler(cache, types.VariableHandlerParams{
					Key:    all,
					Params: &params,
				})

				if err == nil {
					input = strings.ReplaceAll(input, s, res.Result)
				}
			}(s)
		} else {
			wg.Done()
		}
	}

	wg.Wait()
	return input
}
