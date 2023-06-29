package valorant

import (
	"context"
	"fmt"

	"github.com/samber/lo"
	"github.com/satont/twir/apps/parser/internal/types"
)

var Elo = &types.Variable{
	Name:        "valorant.profile.elo",
	Description: lo.ToPtr(`Current elo, i.e "2419"`),
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := types.VariableHandlerResult{}

		profile := parseCtx.Cacher.GetValorantProfile(ctx)
		if profile == nil {
			return nil, nil
		}

		result.Result = fmt.Sprintf("%v", profile.Data.Elo)

		return &result, nil
	},
}
