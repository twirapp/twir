package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"sync"
	testcommand "tsuwari/parser/internal/commands/test"
	"tsuwari/parser/internal/types"
	"tsuwari/parser/internal/variables"
	"tsuwari/parser/internal/variablescache"
	"tsuwari/parser/pkg/helpers"

	"github.com/go-redis/redis/v9"
)

type Commands struct {
	defaultCommands  map[string]*types.DefaultCommand
	redis            *redis.Client
	variablesService variables.Variables
}

func New(redis *redis.Client, variablesService variables.Variables) Commands {
	ctx := Commands{
		redis:            redis,
		defaultCommands:  make(map[string]*types.DefaultCommand),
		variablesService: variablesService,
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

func (c Commands) FindByMessage(input string, cmds *[]types.Command) *types.Command {
	msg := strings.ToLower(input)
	splittedName := regexp.MustCompile(`[^\s]+`).FindAllString(msg, -1)

	var cmd *types.Command

	length := len(splittedName)

	for i := 0; i < length; i++ {
		query := strings.Join(splittedName, " ")
		for _, c := range *cmds {
			if c.Name == query {
				cmd = &c
				break
			}

			if helpers.Contains(c.Aliases, query) {
				cmd = &c
				break
			}
		}

		if cmd != nil {
			break
		} else {

			splittedName = splittedName[:len(splittedName)-1]
			continue
		}
	}

	return cmd
}

func (c Commands) ParseCommandResponses(command *types.Command, data types.HandleProcessCommandData) []string {
	responses := []string{}

	if command.Default && c.defaultCommands[*command.DefaultName] != nil {
		results := c.defaultCommands[*command.DefaultName].Handler(types.VariableHandlerParams{
			Key: "qwe",
		})
		responses = results
	} else {
		responses = command.Responses
	}

	wg := sync.WaitGroup{}
	for i, r := range responses {
		wg.Add(1)
		cacheService := variablescache.New(r, data.Sender.Id, data.Channel.Id, &data.Sender.Name, c.redis, *c.variablesService.Regexp)

		go func(i int, r string) {
			defer wg.Done()
			responses[i] = c.variablesService.ParseInput(cacheService, r)
		}(i, r)
	}
	wg.Wait()

	return responses
}
