package trend

import (
	"fmt"
	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/parser/internal/types"
	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"
	"strings"
)

var ExtendedTrend = types.Variable{
	Name:        "faceit.trend.extended",
	Description: lo.ToPtr(`Faceit latest matches trend in "W +26 | L -23" format`),
	Handler: func(ctx *variables_cache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}

		matches, err := ctx.GetFaceitLatestMatches()
		if err != nil {
			result.Result = err.Error()
			return result, nil
		}

		trend := make([]string, len(matches))

		for _, match := range matches {
			if match.EloDiff == nil {
				continue
			}
			trend = append(
				trend,
				lo.
					If(match.IsWin, fmt.Sprintf("W +%s", *match.EloDiff)).
					Else(fmt.Sprintf("L -%s", *match.EloDiff)),
			)
		}

		result.Result = strings.Join(trend, " | ")
		return result, nil
	},
}
