package random

import (
	"context"
	"strings"

	"github.com/samber/lo"
	"github.com/satont/twir/apps/parser/internal/types"
)

var Phrase = &types.Variable{
	Name:                "random.phrase",
	Description:         lo.ToPtr("Random phrase from list"),
	Example:             lo.ToPtr("random.phrase|Hi there|Kappa|Another Phrase"),
	CanBeUsedInRegistry: true,
	NotCachable:         true,
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}

		if variableData.Params == nil {
			result.Result = "Parameters are not specified"
			return result, nil
		}

		params := strings.Split(*variableData.Params, "|")
		if params == nil {
			result.Result = "Something is wrong with your params"
			return result, nil
		}

		for i, str := range params {
			params[i] = strings.TrimSpace(str)
			if len(params[i]) == 0 {
				result.Result = "Your phrases contains empty phrase, check you writed commas correctly."
				return result, nil
			}
		}

		result.Result = lo.Sample(params)

		return result, nil
	},
}
