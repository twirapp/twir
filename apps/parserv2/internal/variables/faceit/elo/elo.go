package faceitelo

import (
	"strconv"
	"tsuwari/parser/internal/types"
	"tsuwari/parser/internal/variablescache"
)

const Name = "faceit.elo"
const Description = "Faceit elo"

func Handler(ctx *variablescache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
	result := &types.VariableHandlerResult{}

	faceitData := ctx.GetFaceitUserData()

	if faceitData == nil {
		return result, nil
	}

	result.Result = strconv.Itoa(faceitData.Elo)

	return result, nil
}
