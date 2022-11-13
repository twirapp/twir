package streamviewers

import (
	"strconv"

	types "github.com/satont/tsuwari/apps/parser/internal/types"
	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"

	"github.com/samber/lo"
)

var Variable = types.Variable{
	Name:        "stream.viewers",
	Description: lo.ToPtr("Stream viewers"),
	Handler: func(ctx *variables_cache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
		result := types.VariableHandlerResult{}

		stream := ctx.GetChannelStream()
		if stream != nil {
			result.Result = strconv.Itoa(stream.ViewerCount)
		} else {
			result.Result = "offline"
		}

		return &result, nil
	},
}
