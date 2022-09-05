package streamcategory

import (
	types "tsuwari/parser/internal/types"
	variablescache "tsuwari/parser/internal/variablescache"
)

const Name = "stream.category"
const Description = "Stream category"

func Handler(ctx *variablescache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
	result := types.VariableHandlerResult{}

	stream := ctx.GetChannelStream()
	if stream != nil {
		result.Result = stream.GameName
	} else {
		result.Result = "no stream"
	}

	return &result, nil
}
