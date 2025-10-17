package random

import (
	"context"
	"errors"
	"math/rand"
	"strconv"
	"strings"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"
	"github.com/twirapp/twir/libs/i18n"
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
			result.Result = i18n.GetCtx(ctx, locales.Translations.Variables.Random.Errors.NotPassedParams)
			return result, nil
		}

		splittedArgs := strings.Split(*variableData.Params, "-")
		if len(splittedArgs) != 2 {
			result.Result = i18n.GetCtx(ctx, locales.Translations.Variables.Random.Errors.WrongNumber)
			return result, nil
		}

		if variableData.Params != nil {
			first, err := strconv.Atoi(splittedArgs[0])
			if err == nil {
				params[0] = first
			} else {
				result.Result = i18n.GetCtx(ctx, locales.Translations.Variables.Random.Errors.ParseFirstNumber)
				return result, nil
			}
			second, err := strconv.Atoi(splittedArgs[1])
			if err == nil {
				params[1] = second
			} else {
				result.Result = i18n.GetCtx(ctx, locales.Translations.Variables.Random.Errors.ParseSecondNumber)
				return result, nil
			}
		}

		if params[0] > params[1] {
			return nil, errors.New(i18n.GetCtx(ctx, locales.Translations.Variables.Random.Errors.FirstLargerSecond))
		}

		if params[0] < 0 || params[1] < 0 {
			return nil, errors.New(i18n.GetCtx(ctx, locales.Translations.Variables.Random.Errors.LowerNumbers))
		}

		random := params[0] + rand.Intn(params[1]-params[0]+1)
		result.Result = strconv.Itoa(random)

		return result, nil
	},
}
