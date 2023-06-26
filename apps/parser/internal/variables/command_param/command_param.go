package command_param

import (
	"context"

	"github.com/samber/lo"
	"github.com/satont/twir/apps/parser/internal/types"
)

var Variable = &types.Variable{
	Name:        "command.param",
	Description: lo.ToPtr("Whats user typed after command name"),
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		param := ""
		if parseCtx.Text != nil {
			param = *parseCtx.Text
		}

		r := types.VariableHandlerResult{
			Result: param,
		}

		return &r, nil
	},
}
