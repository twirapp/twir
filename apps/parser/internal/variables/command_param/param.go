package command_param

import (
	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/parser/internal/types"
	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"
)

var Variable = types.Variable{
	Name:        "command.param",
	Description: lo.ToPtr("Whats user typed after command name"),
	Handler: func(ctx *variables_cache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
		param := ""
		if ctx.Text != nil {
			param = *ctx.Text
		}

		r := types.VariableHandlerResult{
			Result: param,
		}

		return &r, nil
	},
}
