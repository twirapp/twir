package faceitlvl

import (
	"strconv"
	"tsuwari/parser/internal/types"
	variables_cache "tsuwari/parser/internal/variablescache"

	"github.com/samber/lo"
)

var Variable = types.Variable{
	Name:        "faceit.lvl",
	Description: lo.ToPtr("Faceit Lvl"),
	Handler: func(ctx *variables_cache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}

		faceitData := ctx.GetFaceitUserData()

		if faceitData == nil {
			return result, nil
		}

		result.Result = strconv.Itoa(faceitData.Lvl)

		return result, nil
	},
}
