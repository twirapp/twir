package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"sync"
	testcommand "tsuwari/parser/internal/commands/test"
	testproto "tsuwari/parser/internal/proto"
	"tsuwari/parser/internal/types"
	"tsuwari/parser/internal/variables"
	"tsuwari/parser/internal/variablescache"
	"tsuwari/parser/pkg/helpers"

	"github.com/go-redis/redis/v9"
	"gorm.io/gorm"
)

type Commands struct {
	defaultCommands  map[string]*types.DefaultCommand
	redis            *redis.Client
	variablesService variables.Variables
	Db               *gorm.DB
}

func New(redis *redis.Client, variablesService variables.Variables, db *gorm.DB) Commands {
	ctx := Commands{
		redis:            redis,
		defaultCommands:  make(map[string]*types.DefaultCommand),
		variablesService: variablesService,
		Db:               db,
	}

	ctx.defaultCommands[testcommand.Command.Name] = &testcommand.Command
	return ctx
}

func (c Commands) GetChannelCommands(channelId string) (*[]types.Command, error) {
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

func (c Commands) FindByMessage(input string, cmds *[]types.Command) FindByMessageResult {
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

func (c Commands) ParseCommandResponses(command FindByMessageResult, data testproto.Request) []string {
	responses := []string{}

	cmd := *command.Cmd

	if cmd.Default && c.defaultCommands[*cmd.DefaultName] != nil {
		results := c.defaultCommands[*cmd.DefaultName].Handler(types.VariableHandlerParams{
			Key: "qwe",
		})
		responses = results
	} else {
		responses = cmd.Responses
	}

	cmdParams := strings.TrimSpace(data.Message.Text[len(command.FoundBy)+1:])

	fmt.Println("cmdParams:", cmdParams)

	wg := sync.WaitGroup{}
	for i, r := range responses {
		wg.Add(1)
		// TODO: concatenate all responses into one slice and use it for cache
		cacheService := variablescache.New(cmdParams, data.Sender.Id, data.Channel.Id, &data.Sender.Name, c.redis, *variables.Regexp, c.variablesService.Twitch, c.Db)

		go func(i int, r string) {
			defer wg.Done()

			responses[i] = c.variablesService.ParseInput(cacheService, r)
		}(i, r)
	}
	wg.Wait()

	return responses
}
