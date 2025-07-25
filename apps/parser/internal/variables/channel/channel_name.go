package channel

import (
	"context"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
)

var Name = &types.Variable{
	Name:         "channel.name",
	Description:  lo.ToPtr("Name of twitch channel"),
	CommandsOnly: true,
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		return &types.VariableHandlerResult{Result: parseCtx.Channel.Name}, nil
	},
}
