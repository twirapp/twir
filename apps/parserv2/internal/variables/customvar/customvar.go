package customvar

import (
	"tsuwari/parser/internal/types"
	"tsuwari/parser/internal/variablescache"

	v8 "rogchap.com/v8go"
)

const Name = "customvar"

var Iso = v8.NewIsolate()

func Handler(ctx *variablescache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
	result := &types.VariableHandlerResult{}

	c := v8.NewContext(Iso)
	defer c.Close()
	val, _ := c.RunScript("1 + 1", "value.js")

	result.Result = val.String()

	return result, nil
}
