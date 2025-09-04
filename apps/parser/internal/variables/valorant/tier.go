package valorant

import (
	"context"
	"strconv"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
)

var TierText = &types.Variable{
	Name:        "valorant.profile.tier.text",
	Description: lo.ToPtr(`Tier in text, i.e "Immortal 3"`),
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := types.VariableHandlerResult{}

		profile := parseCtx.Cacher.GetValorantMMR(ctx)
		if profile == nil {
			return nil, nil
		}

		result.Result = profile.Data.Current.Tier.Name

		return &result, nil
	},
}

var RankInTier = &types.Variable{
	Name:        "valorant.profile.tier.ranking",
	Description: lo.ToPtr(`Ranking in tier, i.e "319"`),
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := types.VariableHandlerResult{}

		profile := parseCtx.Cacher.GetValorantMMR(ctx)
		if profile == nil {
			return nil, nil
		}

		result.Result = strconv.Itoa(profile.Data.Current.LeaderboardPlacement.Rank)

		return &result, nil
	},
}
