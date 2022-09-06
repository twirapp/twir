package faceitelo

import (
	"strconv"
	"tsuwari/parser/internal/types"
	"tsuwari/parser/internal/variablescache"

	"github.com/samber/lo"
)

var Variable = types.Variable{
	Name:        "faceit.elo",
	Description: lo.ToPtr("Faceit elo"),
	Handler: func(ctx *variablescache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}

		faceitData := ctx.GetFaceitUserData()

		if faceitData == nil {
			return result, nil
		}

		result.Result = strconv.Itoa(faceitData.Elo)

		return result, nil
	},
}
