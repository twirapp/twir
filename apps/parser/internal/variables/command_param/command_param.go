package command_param

import (
	"context"

	"github.com/samber/lo"
	"github.com/satont/twir/apps/parser/internal/types"
)

var Variable = &types.Variable{
	Name:         "command.param",
	Description:  lo.ToPtr("Whats user typed after command name"),
	CommandsOnly: true,
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		return nil, nil
	},
}
