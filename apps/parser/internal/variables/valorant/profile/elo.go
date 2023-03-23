package valorant_profile

import (
	"fmt"
	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/parser/internal/types"
	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"
)

var Elo = types.Variable{
	Name:         "valorant.profile.elo",
	Description:  lo.ToPtr(`Current elo, i.e "2419"`),
	CommandsOnly: lo.ToPtr(true),
	Handler: func(ctx *variables_cache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
		result := types.VariableHandlerResult{}

		profile := ctx.GetValorantProfile()
		if profile == nil {
			return nil, nil
		}

		result.Result = fmt.Sprintf("%v", profile.Data.Elo)

		return &result, nil
	},
}
