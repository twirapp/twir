package shorturl

import (
	"context"

	"github.com/guregu/null"
	"github.com/lib/pq"
	command_arguments "github.com/twirapp/twir/apps/parser/internal/command-arguments"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/i18n"
)

const (
	urlArgName = "url"
)

var Command = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "shorturl",
		Description: null.StringFrom("Create short url"),
		RolesIDS:    pq.StringArray{},
		Module:      "UTILITY",
		IsReply:     true,
	},
	SkipToxicityCheck: true,
	Args: []command_arguments.Arg{
		command_arguments.String{
			Name: urlArgName,
		},
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		result := &types.CommandsHandlerResult{}

		urlArg := parseCtx.ArgsParser.Get(urlArgName).String()

		link, err := parseCtx.Services.ShortUrlServices.FindOrCreate(
			ctx,
			urlArg,
			parseCtx.Sender.ID,
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(
					ctx,
					locales.Translations.Commands.Shorturl.Errors.CannotCreateShortUrl.SetVars(
						locales.KeysCommandsShorturlErrorsCannotCreateShortUrlVars{
							Error: err.Error(),
						},
					),
				),
				Err: err,
			}
		}

		result.Result = []string{
			i18n.GetCtx(
				ctx,
				locales.Translations.Commands.Shorturl.Success.ShortUrlCreated.SetVars(
					locales.KeysCommandsShorturlSuccessShortUrlCreatedVars{
						Url: link.Short,
					},
				),
			),
		}
		return result, nil
	},
}
