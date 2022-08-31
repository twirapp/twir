package streamid

import (
	"fmt"
	types "tsuwari/parser/internal/types"
	variablescache "tsuwari/parser/internal/variablescache"
)

const Name = "stream.title"

func Handler(ctx *variablescache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
	result := types.VariableHandlerResult{}

	fmt.Println(ctx.Cache.Stream)

	if ctx.Cache.Stream != nil {
		result.Result = ctx.Cache.Stream.Title
	} else {
		result.Result = "no stream"
	}

	return &result, nil
}
