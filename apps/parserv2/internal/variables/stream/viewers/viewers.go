package streamviewers

import (
	"strconv"
	types "tsuwari/parser/internal/types"
	variablescache "tsuwari/parser/internal/variablescache"
)

const Name = "stream.viewers"

func Handler(ctx *variablescache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
	result := types.VariableHandlerResult{}

	if ctx.Cache.Stream != nil {
		result.Result = strconv.Itoa(ctx.Cache.Stream.ViewerCount)
	} else {
		result.Result = "offline"
	}

	return &result, nil
}
