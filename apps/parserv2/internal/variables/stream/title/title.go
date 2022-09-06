package streamid

import (
	types "tsuwari/parser/internal/types"
	variablescache "tsuwari/parser/internal/variablescache"

	"github.com/samber/lo"
)

var Variable = types.Variable{
	Name:        "stream.title",
	Description: lo.ToPtr("Stream title"),
	Handler: func(ctx *variablescache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
		result := types.VariableHandlerResult{}

		stream := ctx.GetChannelStream()
		if stream != nil {
			result.Result = stream.Title
		} else {
			result.Result = "no stream"
		}

		return &result, nil
	},
}
