package streammessages

import (
	"strconv"

	types "github.com/satont/tsuwari/apps/parser/internal/types"
	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"

	"github.com/samber/lo"
)

var Variable = types.Variable{
	Name:        "stream.messages",
	Description: lo.ToPtr("Messages sended by users in this stream"),
	Handler: func(ctx *variables_cache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
		result := types.VariableHandlerResult{}

		stream := ctx.GetChannelStream()
		if stream != nil {
			result.Result = strconv.Itoa(stream.ParsedMessages)
		} else {
			result.Result = "stream offline"
		}

		return &result, nil
	},
}
