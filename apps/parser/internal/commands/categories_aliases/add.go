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
	categoriesaliasesrepository "github.com/twirapp/twir/libs/repositories/channels_categories_aliases"
	"github.com/twirapp/twir/libs/twitch"
)

var Add = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "game alias add",
		Description: null.StringFrom("Add alias for category for shorten usage of set command"),
		RolesIDS: pq.StringArray{
			model.ChannelRoleTypeModerator.String(),
		},
		Module: "MODERATION",
		Aliases: pq.StringArray{
			"game alias create",
		},
		Visible: true,
		IsReply: true,
	},
	Args: []command_arguments.Arg{
		command_arguments.String{
			Name: "alias",
		},
		command_arguments.VariadicString{
			Name: "category",
		},
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		aliasArg := parseCtx.ArgsParser.Get("alias")
		categoryArg := parseCtx.ArgsParser.Get("category")

		if aliasArg == nil || categoryArg == nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(
					ctx,
					locales.Translations.Commands.CategoriesAliases.Errors.CategoryRequired,
				),
			}
		}

		alias := aliasArg.String()
		category := categoryArg.String()

		categories, err := parseCtx.Services.CategoriesAliasesRepo.GetManyByChannelID(
			ctx,
			parseCtx.Channel.ID,
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(
					ctx,
					locales.Translations.Commands.CategoriesAliases.Errors.CategoryFailedToGet,
				),
				Err: err,
			}
		}

		for _, c := range categories {
			if c.Alias == alias {
				return nil, &types.CommandHandlerError{
					Message: i18n.GetCtx(
						ctx,
						locales.Translations.Commands.CategoriesAliases.Errors.AliasAlreadyExists.
							SetVars(locales.KeysCommandsCategoriesAliasesErrorsAliasAlreadyExistsVars{AliasName: alias}),
					),
				}
			}
		}

		foundTwitchCategory, err := twitch.SearchCategory(ctx, category)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(
					ctx,
					locales.Translations.Commands.CategoriesAliases.Errors.CategoryFailedToGet,
				),
				Err: err,
			}
		}

		if foundTwitchCategory == nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(
					ctx,
					locales.Translations.Commands.CategoriesAliases.Errors.CategoryNotFound,
				),
			}
		}

		err = parseCtx.Services.CategoriesAliasesRepo.Create(
			ctx,
			categoriesaliasesrepository.CreateInput{
				ChannelID:  parseCtx.Channel.ID,
				Alias:      alias,
				CategoryID: foundTwitchCategory.ID,
			},
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(
					ctx,
					locales.Translations.Commands.CategoriesAliases.Errors.CategoryFailedToCreate,
				),
				Err: err,
			}
		}

		result := &types.CommandsHandlerResult{
			Result: []string{
				i18n.GetCtx(
					ctx,
					locales.Translations.Commands.CategoriesAliases.Add.AliasAddToCategory.
						SetVars(locales.KeysCommandsCategoriesAliasesAddAliasAddToCategoryVars{AliasName: alias, CategoryName: foundTwitchCategory.Name}),
				),
			},
		}

		return result, nil
	},
}
