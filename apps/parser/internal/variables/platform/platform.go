package platform

import (
	"context"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
)

var Platform = &types.Variable{
	Name:        "platform",
	Description: lo.ToPtr("Platform name"),
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		return &types.VariableHandlerResult{Result: parseCtx.Platform}, nil
	},
}
