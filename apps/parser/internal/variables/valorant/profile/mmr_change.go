package valorant_profile

import (
	"fmt"
	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/parser/internal/types"
	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"
)

var MmrChangeToLastGame = types.Variable{
	Name:         "valorant.profile.mmrChangeFromLastGame",
	Description:  lo.ToPtr(`Mmr change from last game, i.e "38"`),
	CommandsOnly: lo.ToPtr(true),
	Handler: func(ctx *variables_cache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
		result := types.VariableHandlerResult{}

		profile := ctx.GetValorantProfile()
		if profile == nil {
			return nil, nil
		}

		result.Result = fmt.Sprintf("%v", profile.Data.MmrChangeToLastGame)

		return &result, nil
	},
}
