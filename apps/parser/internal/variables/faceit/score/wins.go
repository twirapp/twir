package score

import (
	"fmt"
	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/parser/internal/types"
	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"
)

var Wins = types.Variable{
	Name:        "faceit.score.wins",
	Description: lo.ToPtr(`Faceit wins on stream`),
	Handler: func(ctx *variables_cache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}

		matches, err := ctx.GetFaceitLatestMatches()
		if err != nil {
			result.Result = err.Error()
			return result, nil
		}

		wins := lo.Reduce(matches, func(agg int, item variables_cache.FaceitMatch, _ int) int {
			if item.IsWin {
				return agg + 1
			} else {
				return agg
			}
		}, 0)

		result.Result = fmt.Sprintf("%v", wins)
		return result, nil
	},
}
