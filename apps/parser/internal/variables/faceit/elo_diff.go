package faceit

import (
	"context"
	"strconv"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
)

var EloDiff = &types.Variable{
	Name:                "faceit.todayEloDiff",
	Description:         lo.ToPtr("Faceit today elo earned"),
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

		diff := parseCtx.Cacher.GetFaceitTodayEloDiff(ctx, matches)

		result.Result = strconv.Itoa(diff)

		return result, nil
	},
}
