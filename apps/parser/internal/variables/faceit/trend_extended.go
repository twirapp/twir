package faceit

import (
	"context"
	"fmt"
	"strings"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
)

var TrendExtended = &types.Variable{
	Name:                "faceit.trend.extended",
	Description:         lo.ToPtr(`Faceit latest matches trend in "W +26 | L -23" format`),
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

		trend := make([]string, 0, len(matches))

		for _, match := range matches {
			if match.EloDelta == nil {
				continue
			}
			trend = append(
				trend,
				lo.
					If(match.IsWin, fmt.Sprintf("W +%s", *match.EloDelta)).
					Else(fmt.Sprintf("L -%s", *match.EloDelta)),
			)
		}

		result.Result = strings.Join(trend, " | ")
		return result, nil
	},
}
