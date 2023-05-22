package stream

import (
	"context"

	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/parser/internal/types"
	"github.com/satont/tsuwari/apps/parser/pkg/helpers"
)

var Uptime = &types.Variable{
	Name:        "stream.uptime",
	Description: lo.ToPtr("Prints uptime of stream"),
	Handler: func(ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData) (*types.VariableHandlerResult, error) {
		result := types.VariableHandlerResult{}

		stream := parseCtx.Cacher.GetChannelStream(ctx)
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
