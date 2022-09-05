package messages

import (
	"tsuwari/parser/internal/types"
	"tsuwari/parser/internal/variablescache"
)

const Name = "song"

func Handler(ctx *variablescache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
	var page int = 1

	ctx.Context.Text
}
