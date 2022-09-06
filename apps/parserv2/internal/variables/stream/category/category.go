package streamcategory

import (
	types "tsuwari/parser/internal/types"
	variablescache "tsuwari/parser/internal/variablescache"

	"github.com/samber/lo"
)

var Variable = types.Variable{
	Name:        "stream.category",
	Description: lo.ToPtr("Stream category"),
	Handler: func(ctx *variablescache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
		result := types.VariableHandlerResult{}

		stream := ctx.GetChannelStream()
		if stream != nil {
			result.Result = stream.GameName
		} else {
			result.Result = "no stream"
		}

		return &result, nil
	},
}
