package streamuptime

import (
	types "tsuwari/parser/internal/types"
	variablescache "tsuwari/parser/internal/variablescache"
	"tsuwari/parser/pkg/helpers"
)

const Name = "stream.uptime"

func Handler(ctx *variablescache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
	result := types.VariableHandlerResult{}

	if ctx.Cache.Stream == nil {
		result.Result = "offline"
		return &result, nil
	}

	result.Result = helpers.Duration(ctx.Cache.Stream.StartedAt)

	return &result, nil
}
