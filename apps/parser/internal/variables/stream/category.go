package stream

import (
	"context"

	"github.com/samber/lo"
	"github.com/satont/twir/apps/parser/internal/types"
)

var Category = &types.Variable{
	Name:                "stream.category",
	Description:         lo.ToPtr("Stream category"),
	CanBeUsedInRegistry: true,
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := types.VariableHandlerResult{}

		if parseCtx.ChannelStream != nil {
			result.Result = parseCtx.ChannelStream.GameName
		} else {
			channelInfo := parseCtx.Cacher.GetTwitchChannel(ctx)
			if channelInfo != nil {
				result.Result = channelInfo.GameName
			} else {
				result.Result = "error"
			}
		}

		return &result, nil
	},
}
