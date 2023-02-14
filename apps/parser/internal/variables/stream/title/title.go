package streamid

import (
	types "github.com/satont/tsuwari/apps/parser/internal/types"
	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"

	"github.com/samber/lo"
)

var Variable = types.Variable{
	Name:        "stream.title",
	Description: lo.ToPtr("Stream title"),
	Handler: func(ctx *variables_cache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
		result := types.VariableHandlerResult{}

		stream := ctx.GetChannelStream()

		if stream != nil {
			result.Result = stream.Title
		} else {
			channelInfo := ctx.GetTwitchChannel()
			if channelInfo != nil {
				result.Result = channelInfo.Title
			} else {
				result.Result = "error"
			}
		}

		return &result, nil
	},
}
