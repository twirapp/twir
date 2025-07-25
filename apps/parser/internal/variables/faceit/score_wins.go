package faceit

import (
	"context"
	"fmt"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
)

var ScoreWins = &types.Variable{
	Name:                "faceit.score.wins",
	Description:         lo.ToPtr(`Faceit wins on stream`),
	CanBeUsedInRegistry: true,
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}

		matches, err := parseCtx.Cacher.GetFaceitLatestMatches(ctx)
		if err != nil {
			result.Result = err.Error()
			return result, nil
		}

		wins := lo.Reduce(
			matches, func(agg int, item *types.FaceitMatch, _ int) int {
				if item.IsWin {
					return agg + 1
				} else {
					return agg
				}
			}, 0,
		)

		result.Result = fmt.Sprintf("%v", wins)
		return result, nil
	},
}
