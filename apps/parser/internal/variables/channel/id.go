package channel

import (
	"context"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
)

var ID = &types.Variable{
	Name:         "channel.id",
	Description:  lo.ToPtr("Twitch ID of channel"),
	CommandsOnly: true,
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		return &types.VariableHandlerResult{Result: parseCtx.Channel.ID}, nil
	},
}
