package lifetime

import (
	"context"
	"fmt"
	"strconv"

	"github.com/NovikovRoman/pubg"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/parser/internal/types"
)

var LifetimeKDSolo = &types.Variable{
	Name:        "pubg.lifetime.kdsolo",
	Description: lo.ToPtr("K/D solo"),
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}

		data, err := parseCtx.Cacher.GetPubgLifetimeData(ctx)
		if err != nil {
			result.Result = err.Error()
			return result, nil
		}

		if data == nil {
			return result, nil
		}

		if len(data.Data.Attributes.GameModeStats) == 0 {
			return result, nil
		}

		stats := data.Data.Attributes.GameModeStats[pubg.SoloFPPMode]

		result.Result = fmt.Sprintf("%.2f", float64(stats.Kills)/float64(stats.Losses))

		return result, nil
	},
}

var LifetimeWinsSolo = &types.Variable{
	Name:        "pubg.lifetime.winssolo",
	Description: lo.ToPtr("Wins solo"),
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}

		data, err := parseCtx.Cacher.GetPubgLifetimeData(ctx)
		if err != nil {
			result.Result = err.Error()
			return result, nil
		}

		if data == nil {
			return result, nil
		}

		if len(data.Data.Attributes.GameModeStats) == 0 {
			return result, nil
		}

		stats := data.Data.Attributes.GameModeStats[pubg.SoloFPPMode]

		result.Result = strconv.Itoa(stats.Wins)

		return result, nil
	},
}

var LifetimeMaxKillsSolo = &types.Variable{
	Name:        "pubg.lifetime.maxkillssolo",
	Description: lo.ToPtr("Max kills solo"),
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}

		data, err := parseCtx.Cacher.GetPubgLifetimeData(ctx)
		if err != nil {
			result.Result = err.Error()
			return result, nil
		}

		if data == nil {
			return result, nil
		}

		if len(data.Data.Attributes.GameModeStats) == 0 {
			return result, nil
		}

		stats := data.Data.Attributes.GameModeStats[pubg.SoloFPPMode]

		result.Result = strconv.Itoa(stats.RoundMostKills)

		return result, nil
	},
}

var LifetimeWinrateSolo = &types.Variable{
	Name:        "pubg.lifetime.winratesolo",
	Description: lo.ToPtr("Winrate solo"),
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}

		data, err := parseCtx.Cacher.GetPubgLifetimeData(ctx)
		if err != nil {
			result.Result = err.Error()
			return result, nil
		}

		if data == nil {
			return result, nil
		}

		if len(data.Data.Attributes.GameModeStats) == 0 {
			return result, nil
		}

		stats := data.Data.Attributes.GameModeStats[pubg.SoloFPPMode]

		result.Result = fmt.Sprintf("%.2f", float64(stats.Wins)/float64(stats.Losses)*100)

		return result, nil
	},
}

var LifetimeAverageDamageSolo = &types.Variable{
	Name:        "pubg.lifetime.averagedamagesolo",
	Description: lo.ToPtr("Average damage solo"),
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}

		data, err := parseCtx.Cacher.GetPubgLifetimeData(ctx)
		if err != nil {
			result.Result = err.Error()
			return result, nil
		}

		if data == nil {
			return result, nil
		}

		if len(data.Data.Attributes.GameModeStats) == 0 {
			return result, nil
		}

		stats := data.Data.Attributes.GameModeStats[pubg.SoloFPPMode]

		result.Result = fmt.Sprintf("%.2f", float64(stats.DamageDealt)/float64(stats.RoundsPlayed))

		return result, nil
	},
}
