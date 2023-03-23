package valorant_profile

import (
	"fmt"
	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/parser/internal/types"
	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"
)

var Tier = types.Variable{
	Name:         "valorant.profile.tier",
	Description:  lo.ToPtr(`Tier in number, i.e "26"`),
	CommandsOnly: lo.ToPtr(true),
	Handler: func(ctx *variables_cache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
		result := types.VariableHandlerResult{}

		profile := ctx.GetValorantProfile()
		if profile == nil {
			return nil, nil
		}

		result.Result = fmt.Sprintf("%v", profile.Data.CurrentTier)

		return &result, nil
	},
}

var TierPatched = types.Variable{
	Name:         "valorant.profile.tier.text",
	Description:  lo.ToPtr(`Tier in text, i.e "Immortal 3"`),
	CommandsOnly: lo.ToPtr(true),
	Handler: func(ctx *variables_cache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
		result := types.VariableHandlerResult{}

		profile := ctx.GetValorantProfile()
		if profile == nil {
			return nil, nil
		}

		result.Result = profile.Data.CurrentTierpatched

		return &result, nil
	},
}

var RankingInTier = types.Variable{
	Name:         "valorant.profile.tier.ranking",
	Description:  lo.ToPtr(`Ranking in tier, i.e "319"`),
	CommandsOnly: lo.ToPtr(true),
	Handler: func(ctx *variables_cache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
		result := types.VariableHandlerResult{}

		profile := ctx.GetValorantProfile()
		if profile == nil {
			return nil, nil
		}

		result.Result = fmt.Sprintf("%v", profile.Data.RankingInTier)

		return &result, nil
	},
}
