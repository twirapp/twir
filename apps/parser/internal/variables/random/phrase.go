package random

import (
	"context"
	"strings"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"
	"github.com/twirapp/twir/libs/i18n"
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
			result.Result = i18n.GetCtx(ctx, locales.Translations.Variables.Random.Errors.ParametersNotSpecified)
			return result, nil
		}

		params := strings.Split(*variableData.Params, "|")
		if params == nil {
			result.Result = i18n.GetCtx(ctx, locales.Translations.Variables.Random.Errors.WrongWithParams)
			return result, nil
		}

		for i, str := range params {
			params[i] = strings.TrimSpace(str)
			if len(params[i]) == 0 {
				result.Result = i18n.GetCtx(ctx, locales.Translations.Variables.Random.Errors.EmptyPhrase)
				return result, nil
			}
		}

		result.Result = lo.Sample(params)

		return result, nil
	},
}
