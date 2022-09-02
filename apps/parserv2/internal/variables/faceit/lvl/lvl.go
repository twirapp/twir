package faceitlvl

import (
	"strconv"
	"tsuwari/parser/internal/types"
	"tsuwari/parser/internal/variablescache"
)

const Name = "faceit.lvl"

func Handler(ctx *variablescache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
	result := &types.VariableHandlerResult{}

	faceitData, err := ctx.GetFaceitData()

	if err != nil {
		result.Result = err.Error()
		return result, nil
	}

	if faceitData == nil {
		return result, nil
	}

	result.Result = strconv.Itoa(faceitData.Lvl)

	return result, nil
}
