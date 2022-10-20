package touser

import (
	"tsuwari/parser/internal/types"
	variables_cache "tsuwari/parser/internal/variablescache"

	"github.com/samber/lo"
)

var Variable = types.Variable{
	Name:         "touser",
	Description:  lo.ToPtr("Mention user"),
	CommandsOnly: lo.ToPtr(true),
	Handler: func(ctx *variables_cache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
		result := types.VariableHandlerResult{
			Result: ctx.SenderName,
		}

		if ctx.Text != nil {
			result.Result = *ctx.Text
		}

		return &result, nil
	},
}
