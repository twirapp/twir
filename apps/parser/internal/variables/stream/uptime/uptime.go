package streamuptime

import (
	types "tsuwari/parser/internal/types"
	variables_cache "tsuwari/parser/internal/variablescache"
	"tsuwari/parser/pkg/helpers"

	"github.com/samber/lo"
)

var Variable = types.Variable{
	Name:        "stream.uptime",
	Description: lo.ToPtr("Stream uptime"),
	Handler: func(ctx *variables_cache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
		result := types.VariableHandlerResult{}

		stream := ctx.GetChannelStream()
		if stream == nil {
			result.Result = "offline"
			return &result, nil
		}

		result.Result = helpers.Duration(stream.StartedAt, lo.ToPtr(true))

		return &result, nil
	},
}
