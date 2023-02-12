package trend

import (
	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/parser/internal/types"
	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"
)

var SimpleTrend = types.Variable{
	Name:        "faceit.trend.simple",
	Description: lo.ToPtr(`Faceit latest matches trend in "WWW" format`),
	Handler: func(ctx *variables_cache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}

		matches, err := ctx.GetFaceitLatestMatches()
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
