package faceit

import (
	"context"
	"fmt"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
)

var ScoreLoses = &types.Variable{
	Name:                "faceit.score.loses",
	Description:         lo.ToPtr(`Faceit loses on stream`),
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

		loses := lo.Reduce(
			matches, func(agg int, item *types.FaceitMatch, _ int) int {
				if item.IsWin {
					return agg
				} else {
					return agg + 1
				}
			}, 0,
		)

		result.Result = fmt.Sprintf("%v", loses)
		return result, nil
	},
}
