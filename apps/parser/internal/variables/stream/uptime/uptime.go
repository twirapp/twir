package streamuptime

import (
	types "github.com/satont/tsuwari/apps/parser/internal/types"
	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"
	"github.com/satont/tsuwari/apps/parser/pkg/helpers"

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

		result.Result = helpers.Duration(stream.StartedAt, &helpers.DurationOpts{
			UseUtc: true,
			Hide:   helpers.DurationOptsHide{},
		})

		return &result, nil
	},
}
