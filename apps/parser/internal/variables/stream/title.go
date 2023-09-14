package stream

import (
	"context"

	"github.com/samber/lo"
	"github.com/satont/twir/apps/parser/internal/types"
)

var Title = &types.Variable{
	Name:                "stream.title",
	Description:         lo.ToPtr("Stream title"),
	CanBeUsedInRegistry: true,
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := types.VariableHandlerResult{}

		stream := parseCtx.Cacher.GetChannelStream(ctx)
		if stream != nil {
			result.Result = stream.Title
		} else {
			channelInfo := parseCtx.Cacher.GetTwitchChannel(ctx)
			if channelInfo != nil {
				result.Result = channelInfo.Title
			} else {
				result.Result = "error"
			}
		}

		return &result, nil
	},
}
