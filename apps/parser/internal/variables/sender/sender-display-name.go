package sender

import (
	"context"

	"github.com/samber/lo"
	"github.com/satont/twir/apps/parser/internal/types"
)

var SenderDisplayName = &types.Variable{
	Name:        "sender.displayName",
	Description: lo.ToPtr("Username of user, who sended message"),
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := types.VariableHandlerResult{Result: parseCtx.Sender.DisplayName}

		return &result, nil
	},
}
