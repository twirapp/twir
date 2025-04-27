package stream

import (
	"context"
	"strconv"

	"github.com/samber/lo"
	"github.com/satont/twir/apps/parser/internal/types"
)

var Viewers = &types.Variable{
	Name:                "stream.viewers",
	Description:         lo.ToPtr("Stream viewers"),
	CanBeUsedInRegistry: true,
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := types.VariableHandlerResult{}

		if parseCtx.ChannelStream != nil {
			result.Result = strconv.Itoa(parseCtx.ChannelStream.ViewerCount)
		} else {
			result.Result = "offline"
		}

		return &result, nil
	},
}
