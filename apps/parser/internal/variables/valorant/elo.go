package valorant

import (
	"context"
	"strconv"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
)

var Elo = &types.Variable{
	Name:                "valorant.profile.elo",
	Description:         lo.ToPtr(`Current elo, i.e "2419"`),
	CanBeUsedInRegistry: true,
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := types.VariableHandlerResult{}

		profile := parseCtx.Cacher.GetValorantProfile(ctx)
		if profile == nil {
			return nil, nil
		}

		result.Result = strconv.Itoa(profile.Data.CurrentData.Elo)

		return &result, nil
	},
}
