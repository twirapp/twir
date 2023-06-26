package command_counters

import (
	"context"

	"github.com/samber/lo"
	"github.com/satont/twir/apps/parser/internal/types"
)

var CommandCounter = &types.Variable{
	Name:         "command.counter",
	Description:  lo.ToPtr("Counter saying how many times command was used"),
	CommandsOnly: true,
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}

		count, err := getCount(parseCtx.Services.Gorm, parseCtx.Command.ID, nil)
		if err != nil {
			parseCtx.Services.Logger.Sugar().Error(err)

			result.Result = "cannot get count"
			return result, nil
		}

		result.Result = count

		return result, nil
	},
}
