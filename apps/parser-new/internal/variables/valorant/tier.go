package valorant

import (
	"context"
	"fmt"

	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/parser-new/internal/types"
)

var Tier = &types.Variable{
	Name:        "valorant.profile.tier",
	Description: lo.ToPtr(`Tier in number, i.e "26"`),
	Handler: func(ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData) (*types.VariableHandlerResult, error) {
		result := types.VariableHandlerResult{}

		profile := parseCtx.Cacher.GetValorantProfile(ctx)
		if profile == nil {
			return nil, nil
		}

		result.Result = fmt.Sprintf("%v", profile.Data.CurrentTier)

		return &result, nil
	},
}

var TierText = &types.Variable{
	Name:        "valorant.profile.tier.text",
	Description: lo.ToPtr(`Tier in text, i.e "Immortal 3"`),
	Handler: func(ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData) (*types.VariableHandlerResult, error) {
		result := types.VariableHandlerResult{}

		profile := parseCtx.Cacher.GetValorantProfile(ctx)
		if profile == nil {
			return nil, nil
		}

		result.Result = profile.Data.CurrentTierpatched

		return &result, nil
	},
}

var RankInTier = &types.Variable{
	Name:        "valorant.profile.tier.ranking",
	Description: lo.ToPtr(`Ranking in tier, i.e "319"`),
	Handler: func(ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData) (*types.VariableHandlerResult, error) {
		result := types.VariableHandlerResult{}

		profile := parseCtx.Cacher.GetValorantProfile(ctx)
		if profile == nil {
			return nil, nil
		}

		result.Result = fmt.Sprintf("%v", profile.Data.RankingInTier)

		return &result, nil
	},
}
