package shorturl

import (
	"context"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"
	"github.com/twirapp/twir/libs/i18n"
)

var Variable = &types.Variable{
	Name:         "shorturl",
	Description:  lo.ToPtr("Create short url from your link"),
	Example:      lo.ToPtr("shorturl|https://example.com"),
	CommandsOnly: false,
	NotCachable:  true,
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		if variableData.Params == nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(ctx, locales.Translations.Variables.Shorturl.Errors.UrlRequired),
			}
		}

		link, err := parseCtx.Services.ShortUrlServices.FindOrCreate(
			ctx,
			*variableData.Params,
			parseCtx.Sender.ID,
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(
					ctx,
					locales.Translations.Variables.Shorturl.Errors.CreateShortUrl.
						SetVars(locales.KeysVariablesShorturlErrorsCreateShortUrlVars{Reason: err.Error()}),
				),
				Err: err,
			}
		}

		result := types.VariableHandlerResult{Result: link.Short}

		return &result, nil
	},
}
