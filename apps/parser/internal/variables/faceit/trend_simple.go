package faceit

import (
	"context"

	"github.com/samber/lo"
	"github.com/satont/twir/apps/parser/internal/types"
)

var TrendSimple = &types.Variable{
	Name:        "faceit.trend.simple",
	Description: lo.ToPtr(`Faceit latest matches trend in "WWW" format`),
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}

		matches, err := parseCtx.Cacher.GetFaceitLatestMatches(ctx)
		if err != nil {
			result.Result = err.Error()
			return result, nil
		}

		trend := ""

		for _, match := range matches {
			trend += lo.If(match.IsWin, "W").Else("L")
		}

		result.Result = trend
		return result, nil
	},
}
