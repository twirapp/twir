package categories_aliases

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

var Remove = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "game alias remove",
		Description: null.StringFrom("Remove category alias"),
		RolesIDS: pq.StringArray{
			model.ChannelRoleTypeModerator.String(),
		},
		Module: "MODERATION",
		Aliases: pq.StringArray{
			"game alias delete",
		},
		Visible: true,
		IsReply: true,
	},
	Args: []command_arguments.Arg{
		command_arguments.String{
			Name: "alias",
		},
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		alias := parseCtx.ArgsParser.Get("alias").String()

		categories, err := parseCtx.Services.CategoriesAliasesRepo.GetManyByChannelID(
			ctx,
			parseCtx.Channel.ID,
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: "cannot get categories",
				Err:     err,
			}
		}

		var found bool

		for _, category := range categories {
			if category.Alias == alias {
				found = true
				err = parseCtx.Services.CategoriesAliasesRepo.Delete(ctx, category.ID)
				if err != nil {
					return nil, &types.CommandHandlerError{
						Message: i18n.GetCtx(
							ctx,
							locales.Translations.Commands.CategoriesAliases.Errors.CategoryCannotDelete,
						),
						Err: err,
					}
				}
				break
			}
		}

		if !found {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(
					ctx,
					locales.Translations.Commands.CategoriesAliases.Errors.AliasNotFound,
				),
			}
		}

		result := &types.CommandsHandlerResult{
			Result: []string{
				i18n.GetCtx(
					ctx,
					locales.Translations.Commands.CategoriesAliases.Errors.AliasRemoved.
						SetVars(locales.KeysCommandsCategoriesAliasesErrorsAliasRemovedVars{AliasName: alias}),
				),
			},
		}

		return result, nil
	},
}
