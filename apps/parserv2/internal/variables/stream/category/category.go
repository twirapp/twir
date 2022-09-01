package streamcategory

import (
	types "tsuwari/parser/internal/types"
	variablescache "tsuwari/parser/internal/variablescache"
)

const Name = "stream.category"

func Handler(ctx *variablescache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
	result := types.VariableHandlerResult{}

	if ctx.Cache.Stream != nil {
		result.Result = ctx.Cache.Stream.GameName
	} else {
		result.Result = "no stream"
	}

	return &result, nil
}
