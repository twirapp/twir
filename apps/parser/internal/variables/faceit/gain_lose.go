package faceit

import (
	"context"
	"fmt"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"go.uber.org/zap"
)

var GainLose = &types.Variable{
	Name:                "faceit.gain",
	Description:         lo.ToPtr("Faceit match gain/lose elo"),
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

		result.Result = fmt.Sprintf("+%d/-%d", res.Gain, res.Lose)

		return result, nil
	},
}
