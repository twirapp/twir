package command_counter

import (
	"tsuwari/parser/internal/types"
	variables_cache "tsuwari/parser/internal/variablescache"

	"github.com/samber/lo"
)

var UserVariable = types.Variable{
	Name:         "command.counter.user",
	Description:  lo.ToPtr("Counter saying how many times command was used by sender user"),
	CommandsOnly: lo.ToPtr(true),
	Handler: func(ctx *variables_cache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}
		count, err := getCount(ctx.Services.Db, ctx.Command.ID, &ctx.SenderId)
		if err != nil {
			result.Result = "cannot get count"
			return result, nil
		}
		result.Result = count

		return result, nil
	},
}
