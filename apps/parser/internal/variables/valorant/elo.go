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

		profile := parseCtx.Cacher.GetValorantMMR(ctx)
		if profile == nil {
			return nil, nil
		}

		result.Result = strconv.Itoa(profile.Data.Current.Elo)

		return &result, nil
	},
}

var EloLastChange = &types.Variable{
	Name:        "valorant.profile.elo.last_change",
	Description: lo.ToPtr(`Last game elo change`),
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := types.VariableHandlerResult{}

		profile := parseCtx.Cacher.GetValorantMMR(ctx)
		if profile == nil {
			return nil, nil
		}

		result.Result = strconv.Itoa(profile.Data.Current.LastChange)

		return &result, nil
	},
}

var RR = &types.Variable{
	Name:        "valorant.profile.rr",
	Description: lo.ToPtr(`Valorant Rank Rating (RR), i.e "42"`),
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := types.VariableHandlerResult{}

		profile := parseCtx.Cacher.GetValorantMMR(ctx)
		if profile == nil {
			return nil, nil
		}

		result.Result = strconv.Itoa(profile.Data.Current.Rr)

		return &result, nil
	},
}
