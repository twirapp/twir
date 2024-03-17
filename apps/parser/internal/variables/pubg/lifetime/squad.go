package lifetime

import (
	"context"
	"fmt"

	"github.com/NovikovRoman/pubg"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/parser/internal/types"
)

var LifetimeKDSquad = &types.Variable{
	Name:        "pubg.lifetime.kdsquad",
	Description: lo.ToPtr("K/D squad"),
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}

		data, err := parseCtx.Cacher.GetPubgLifetimeData(ctx)
		if err != nil {
			result.Result = err.Error()
			return result, nil
		}

		if len(data.Data.Attributes.GameModeStats) == 0 {
			return result, nil
		}

		stats := data.Data.Attributes.GameModeStats[pubg.SquadFPPMode]

		result.Result = fmt.Sprintf("%.2f", float64(stats.Kills)/float64(stats.Losses))

		return result, nil
	},
}
