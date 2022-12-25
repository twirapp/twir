package command_counter

import (
	"github.com/satont/tsuwari/apps/parser/internal/types"
	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"

	"github.com/samber/lo"
)

var CommandVariable = types.Variable{
	Name:         "command.counter",
	Description:  lo.ToPtr("Counter saying how many times command was used"),
	CommandsOnly: lo.ToPtr(true),
	Handler: func(ctx *variables_cache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}

		count, err := getCount(ctx.Services.Db, ctx.Command.ID, nil)
		if err != nil {
			result.Result = "cannot get count"
			return result, nil
		}
		result.Result = count

		return result, nil
	},
}
