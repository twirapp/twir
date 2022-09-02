package streamuptime

import (
	types "tsuwari/parser/internal/types"
	variablescache "tsuwari/parser/internal/variablescache"
	"tsuwari/parser/pkg/helpers"
)

const Name = "stream.uptime"

func Handler(ctx *variablescache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
	result := types.VariableHandlerResult{}

	stream := ctx.GetChannelStream()
	if stream == nil {
		result.Result = "offline"
		return &result, nil
	}

	result.Result = helpers.Duration(stream.StartedAt)

	return &result, nil
}
