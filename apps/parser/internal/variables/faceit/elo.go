package faceit

import (
	"context"
	"strconv"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
)

var Elo = &types.Variable{
	Name:                "faceit.elo",
	Description:         lo.ToPtr("Faceit elo"),
	CanBeUsedInRegistry: true,
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}

		faceitData, err := parseCtx.Cacher.GetFaceitUserData(ctx)
		if err != nil {
			result.Result = err.Error()
			return result, nil
		}

		if faceitData == nil {
			return result, nil
		}

		result.Result = strconv.Itoa(faceitData.Elo)

		return result, nil
	},
}
