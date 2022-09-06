package faceitlvl

import (
	"strconv"
	"tsuwari/parser/internal/types"
	"tsuwari/parser/internal/variablescache"

	"github.com/samber/lo"
)

var Variable = types.Variable{
	Name:        "faceit.lvl",
	Description: lo.ToPtr("Faceit Lvl"),
	Handler: func(ctx *variablescache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}

		faceitData := ctx.GetFaceitUserData()

		if faceitData == nil {
			return result, nil
		}

		result.Result = strconv.Itoa(faceitData.Lvl)

		return result, nil
	},
}
