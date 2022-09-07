package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"sync"
	channel_game "tsuwari/parser/internal/commands/channel/game"
	channel_title "tsuwari/parser/internal/commands/channel/title"
	"tsuwari/parser/internal/commands/nuke"
	"tsuwari/parser/internal/commands/permit"
	"tsuwari/parser/internal/commands/spam"
	usersauth "tsuwari/parser/internal/twitch/user"
	"tsuwari/parser/internal/types"
	"tsuwari/parser/internal/variables"
	variables_cache "tsuwari/parser/internal/variablescache"
	"tsuwari/parser/pkg/helpers"

	"github.com/nats-io/nats.go"
	"github.com/samber/lo"
	parserproto "github.com/satont/tsuwari/nats/parser"

	"github.com/go-redis/redis/v9"
	"gorm.io/gorm"
)

type Commands struct {
	DefaultCommands  []types.DefaultCommand
	redis            *redis.Client
	variablesService variables.Variables
	Db               *gorm.DB
	UsersAuth        *usersauth.UsersTokensService
	Nats             *nats.Conn
}

type CommandsOpts struct {
	Redis            *redis.Client
	VariablesService variables.Variables
	Db               *gorm.DB
	UsersAuth        *usersauth.UsersTokensService
	Nats             *nats.Conn
}

func New(opts CommandsOpts) Commands {
	commands := []types.DefaultCommand{
		channel_title.Command,
		channel_game.Command,
		permit.Command,
		spam.Command,
		nuke.Command,
	}

	ctx := Commands{
		redis:            opts.Redis,
		DefaultCommands:  commands,
		variablesService: opts.VariablesService,
		Db:               opts.Db,
		UsersAuth:        opts.UsersAuth,
		Nats:             opts.Nats,
	}

	return ctx
}

func (c *Commands) GetChannelCommands(channelId string) (*[]types.Command, error) {
	rCtx := context.TODO()
	keys, err := c.redis.Keys(rCtx, fmt.Sprintf("commands:%s:*", channelId)).Result()

	if err != nil {
		return nil, err
	}

	var cmds = make([]types.Command, len(keys))
	rCmds, err := c.redis.MGet(rCtx, keys...).Result()

	if err != nil {
		return nil, err
	}

	for i, cmd := range rCmds {
		parsedCmd := types.Command{}

		err := json.Unmarshal([]byte(cmd.(string)), &parsedCmd)

		if err == nil {
			cmds[i] = parsedCmd
		}
	}

	return &cmds, nil
}

var splittedNameRegexp = regexp.MustCompile(`[^\s]+`)

type FindByMessageResult struct {
	Cmd     *types.Command
	FoundBy string
}

func (c *Commands) FindByMessage(input string, cmds *[]types.Command) FindByMessageResult {
	msg := strings.ToLower(input)
	splittedName := splittedNameRegexp.FindAllString(msg, -1)

	res := FindByMessageResult{}

	length := len(splittedName)

	for i := 0; i < length; i++ {
		query := strings.Join(splittedName, " ")
		for _, c := range *cmds {
			if c.Name == query {
				res.FoundBy = query
				res.Cmd = &c
				break
			}

			if helpers.Contains(c.Aliases, query) {
				res.FoundBy = query
				res.Cmd = &c
				break
			}
		}

		if res.Cmd != nil {
			break
		} else {

			splittedName = splittedName[:len(splittedName)-1]
			continue
		}
	}

	return res
}

func (c *Commands) ParseCommandResponses(command FindByMessageResult, data parserproto.Request) []string {
	responses := []string{}

	cmd := *command.Cmd
	var cmdParams *string
	params := strings.TrimSpace(data.Message.Text[len(command.FoundBy):])
	if len(params) > 0 {
		cmdParams = &params
	}

	defaultCommand, isDefaultExists := lo.Find(c.DefaultCommands, func(command types.DefaultCommand) bool {
		if cmd.DefaultName != nil {
			return command.Name == *cmd.DefaultName
		} else {
			return false
		}
	})

	if cmd.Default && isDefaultExists {
		results := defaultCommand.Handler(variables_cache.ExecutionContext{
			ChannelName: data.Channel.Name,
			ChannelId:   data.Channel.Id,
			SenderId:    data.Sender.Id,
			SenderName:  data.Sender.Name,
			Text:        cmdParams,
			Services: variables_cache.ExecutionServices{
				Redis:     c.redis,
				Regexp:    nil,
				Twitch:    c.variablesService.Twitch,
				Db:        c.Db,
				UsersAuth: c.UsersAuth,
				Nats:      c.Nats,
			},
		})
		responses = results
	} else {
		responses = cmd.Responses
	}

	wg := sync.WaitGroup{}
	for i, r := range responses {
		wg.Add(1)
		// TODO: concatenate all responses into one slice and use it for cache
		cacheService := variables_cache.New(variables_cache.VariablesCacheOpts{
			Text:       cmdParams,
			SenderId:   data.Sender.Id,
			ChannelId:  data.Channel.Id,
			SenderName: &data.Sender.DisplayName,
			Redis:      c.redis,
			Regexp:     variables.Regexp,
			Twitch:     c.variablesService.Twitch,
			DB:         c.Db,
		})

		go func(i int, r string) {
			defer wg.Done()

			responses[i] = c.variablesService.ParseInput(cacheService, r)
		}(i, r)
	}
	wg.Wait()

	return responses
}
