package faceitelo

import (
	"strconv"
	"tsuwari/parser/internal/types"
	"tsuwari/parser/internal/variablescache"
)

const Name = "faceit.elo"

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

	result.Result = strconv.Itoa(faceitData.Elo)

	return result, nil
}
