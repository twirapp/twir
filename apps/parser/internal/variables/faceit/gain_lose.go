package faceit

import (
	"context"
	"strconv"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"go.uber.org/zap"
)

var Gain = &types.Variable{
	Name:                "faceit.gain",
	Description:         lo.ToPtr("Faceit match gain elo"),
	CanBeUsedInRegistry: true,
	Handler: func(
		ctx context.Context,
		parseCtx *types.VariableParseContext,
		variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}

		_, err := parseCtx.Cacher.GetFaceitUserData(ctx)
		if err != nil {
			zap.S().Error(err)
			result.Result = "N/A"
			return result, nil
		}

		res, err := parseCtx.Cacher.ComputeFaceitGainLoseEstimate(ctx)
		if err != nil || res == nil {
			if err != nil {
				zap.S().Error(err)
			}
			result.Result = "N/A"
			return result, nil
		}

		result.Result = strconv.Itoa(res.Gain)

		return result, nil
	},
}

var Lose = &types.Variable{
	Name:        "faceit.lose",
	Description: lo.ToPtr("Faceit match lose"),
	Handler: func(
		ctx context.Context,
		parseCtx *types.VariableParseContext,
		variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}

		_, err := parseCtx.Cacher.GetFaceitUserData(ctx)
		if err != nil {
			zap.S().Error(err)
			result.Result = "N/A"
			return result, nil
		}

		res, err := parseCtx.Cacher.ComputeFaceitGainLoseEstimate(ctx)
		if err != nil || res == nil {
			if err != nil {
				zap.S().Error(err)
			}
			result.Result = "N/A"
			return result, nil
		}

		result.Result = strconv.Itoa(res.Lose)

		return result, nil
	},
}
