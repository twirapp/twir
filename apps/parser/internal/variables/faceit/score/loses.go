package score

import (
	"fmt"
	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/parser/internal/types"
	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"
)

var Loses = types.Variable{
	Name:        "faceit.score.loses",
	Description: lo.ToPtr(`Faceit loses on stream`),
	Handler: func(ctx *variables_cache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}

		ctx.GetFaceitUserData()
		matches, err := ctx.GetFaceitLatestMatches()
		if err != nil {
			result.Result = err.Error()
			return result, nil
		}

		wins := lo.Reduce(matches, func(agg int, item variables_cache.FaceitMatch, _ int) int {
			if item.IsWin {
				return agg
			} else {
				return agg + 1
			}
		}, 0)

		result.Result = fmt.Sprintf("%v", wins)
		return result, nil
	},
}
