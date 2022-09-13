package faceitelodiff

import (
	"strconv"
	"tsuwari/parser/internal/types"
	variables_cache "tsuwari/parser/internal/variablescache"

	"github.com/samber/lo"
)

var Variable = types.Variable{
	Name:        "faceit.todayEloDiff",
	Description: lo.ToPtr("Faceit today elo earned"),
	Handler: func(ctx *variables_cache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}

		matches, err := ctx.GetFaceitLatestMatches()
		
		if err != nil {
			result.Result = err.Error()
			return result, nil
		}

		diff := ctx.GetFaceitTodayEloDiff(matches)

		result.Result = strconv.Itoa(diff)

		return result, nil
	},
}
