package sender

import (
	"context"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
)

var ID = &types.Variable{
	Name:        "sender.id",
	Description: lo.ToPtr("Sender id"),
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := types.VariableHandlerResult{Result: parseCtx.Sender.ID}

		return &result, nil
	},
}
