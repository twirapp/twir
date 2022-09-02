package streammessages

import (
	"strconv"
	types "tsuwari/parser/internal/types"
	variablescache "tsuwari/parser/internal/variablescache"
)

const Name = "stream.messages"

func Handler(ctx *variablescache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
	result := types.VariableHandlerResult{}

	stream := ctx.GetChannelStream()
	if stream != nil {
		result.Result = strconv.Itoa(stream.Messages)
	} else {
		result.Result = "no stream"
	}

	return &result, nil
}
