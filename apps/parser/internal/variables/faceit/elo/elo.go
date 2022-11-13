package faceitelo

import (
	"strconv"

	"github.com/satont/tsuwari/apps/parser/internal/types"
	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"

	"github.com/samber/lo"
)

var Variable = types.Variable{
	Name:        "faceit.elo",
	Description: lo.ToPtr("Faceit elo"),
	Handler: func(ctx *variables_cache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}

		faceitData, err := ctx.GetFaceitUserData()
		if err != nil {
			result.Result = err.Error()
			return result, nil
		}

		if faceitData == nil {
			return result, nil
		}

		result.Result = strconv.Itoa(faceitData.Elo)

		return result, nil
	},
}
