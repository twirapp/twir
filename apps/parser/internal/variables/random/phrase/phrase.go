package random

import (
	"strings"
	types "tsuwari/parser/internal/types"
	variables_cache "tsuwari/parser/internal/variablescache"

	"github.com/samber/lo"
)

const exampleStr = "Please check example."

var Variable = types.Variable{
	Name:        "random.phrase",
	Description: lo.ToPtr("Random phrase from list"),
	Example:     lo.ToPtr("random.phrase|Hi there|Kappa|Another Phrase"),
	Handler: func(ctx *variables_cache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}

		if data.Params == nil {
			result.Result = "Parameters are not specified"
			return result, nil
		}

		params := strings.Split(*data.Params, "|")
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
