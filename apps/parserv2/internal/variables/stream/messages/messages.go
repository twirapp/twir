package streammessages

import (
	"strconv"
	types "tsuwari/parser/internal/types"
	variables_cache "tsuwari/parser/internal/variablescache"

	"github.com/samber/lo"
)

var Variable = types.Variable{
	Name:        "stream.messages",
	Description: lo.ToPtr("Messages sended by users in this stream"),
	Handler: func(ctx *variables_cache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
		result := types.VariableHandlerResult{}

		stream := ctx.GetChannelStream()
		if stream != nil {
			result.Result = strconv.Itoa(stream.Messages)
		} else {
			result.Result = "stream offline"
		}

		return &result, nil
	},
}
