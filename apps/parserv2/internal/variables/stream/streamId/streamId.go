package streamId

import (
	types "tsuwari/parser/internal/types"
	variablescache "tsuwari/parser/internal/variablescache"
)

const Name = "streamId"

func Handler(ctx *variablescache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
	result := types.VariableHandlerResult{}

	if ctx.Cache.StreamId != nil {
		result.Result = *ctx.Cache.StreamId
	} else {
		result.Result = "no stream"
	}

	return &result, nil
}
