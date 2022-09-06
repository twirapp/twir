package commandslist

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"tsuwari/parser/internal/types"
	"tsuwari/parser/internal/variablescache"
	"tsuwari/parser/pkg/helpers"

	"github.com/samber/lo"
)

const Name = "commands.list"
const Description = "List of commands"

var Variable = types.Variable{
	Name: "commands.list",
	Handler: func(ctx *variablescache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
		rCtx := context.TODO()
		keys, err := ctx.Services.Redis.Keys(rCtx, fmt.Sprintf("commands:%s:*", ctx.Context.ChannelId)).Result()

		if err != nil {
			return nil, err
		}

		var cmds = make([]types.Command, len(keys))
		rCmds, err := ctx.Services.Redis.MGet(rCtx, keys...).Result()

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

		mapped := helpers.Map(cmds, func(c types.Command) string {
			return "!" + c.Name
		})

		r := types.VariableHandlerResult{
			Result: strings.Join(mapped, ", "),
		}

		return &r, nil
	},
	Description: lo.ToPtr("Command list"),
}

/* func Handler(ctx *variablescache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
	rCtx := context.TODO()
	keys, err := ctx.Services.Redis.Keys(rCtx, fmt.Sprintf("commands:%s:*", ctx.Context.ChannelId)).Result()

	if err != nil {
		return nil, err
	}

	var cmds = make([]types.Command, len(keys))
	rCmds, err := ctx.Services.Redis.MGet(rCtx, keys...).Result()

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

	mapped := helpers.Map(cmds, func(c types.Command) string {
		return "!" + c.Name
	})

	r := types.VariableHandlerResult{
		Result: strings.Join(mapped, ", "),
	}

	return &r, nil
} */
