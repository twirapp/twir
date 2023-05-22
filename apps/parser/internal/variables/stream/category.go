package stream

import (
	"context"

	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/parser/internal/types"
)

var Category = &types.Variable{
	Name:        "stream.category",
	Description: lo.ToPtr("Stream category"),
	Handler: func(ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData) (*types.VariableHandlerResult, error) {
		result := types.VariableHandlerResult{}

		stream := parseCtx.Cacher.GetChannelStream(ctx)
		if stream != nil {
			result.Result = stream.GameName
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
