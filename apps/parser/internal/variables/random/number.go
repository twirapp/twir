package random

import (
	"context"
	"errors"
	"math/rand"
	"strconv"
	"strings"

	"github.com/samber/lo"
	"github.com/satont/twir/apps/parser/internal/types"
)

var Number = &types.Variable{
	Name:                "random",
	Description:         lo.ToPtr("Random number from N to N"),
	Example:             lo.ToPtr("random|1-100"),
	CanBeUsedInRegistry: true,
	NotCachable:         true,
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}

		params := [2]int{}
		if variableData.Params == nil {
			result.Result = "Have not passed params to random variable."
			return result, nil
		}

		splittedArgs := strings.Split(*variableData.Params, "-")
		if len(splittedArgs) != 2 {
			result.Result = "Wrong number of arguments passed to random."
			return result, nil
		}

		if variableData.Params != nil {
			first, err := strconv.Atoi(splittedArgs[0])
			if err == nil {
				params[0] = first
			} else {
				result.Result = "cannot parse first number from arguments."
				return result, nil
			}
			second, err := strconv.Atoi(splittedArgs[1])
			if err == nil {
				params[1] = second
			} else {
				result.Result = "cannot parse second number from arguments. "
				return result, nil
			}
		}

		if params[0] > params[1] {
			return nil, errors.New("first number cannot be larger then second")
		}

		if params[0] < 0 || params[1] < 0 {
			return nil, errors.New("numbers cannot be lower then 0")
		}

		random := params[0] + rand.Intn(params[1]-params[0]+1)
		result.Result = strconv.Itoa(random)

		return result, nil
	},
}
