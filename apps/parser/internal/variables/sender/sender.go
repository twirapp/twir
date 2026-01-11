package sender

import (
	"context"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
)

var Sender = &types.Variable{
	Name:         "sender",
	Description:  lo.ToPtr("Username of user, who sended message"),
	CommandsOnly: true,
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := types.VariableHandlerResult{Result: parseCtx.Sender.Name}

		return &result, nil
	},
}
