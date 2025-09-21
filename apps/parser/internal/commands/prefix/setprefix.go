package prefix

import (
	"context"
	"errors"
	"unicode/utf8"

	"github.com/guregu/null"
	"github.com/lib/pq"
	command_arguments "github.com/twirapp/twir/apps/parser/internal/command-arguments"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/i18n"
	channelscommandsprefixrepository "github.com/twirapp/twir/libs/repositories/channels_commands_prefix"
	channelscommandsprefixmodel "github.com/twirapp/twir/libs/repositories/channels_commands_prefix/model"
)

const setPrefixArgName = "prefix"

var SetPrefix = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name: "prefix set",
		Description: null.StringFrom(
			"Set prefix for commands",
		),
		RolesIDS: pq.StringArray{model.ChannelRoleTypeModerator.String()},
		Module:   "MODERATION",
		IsReply:  true,
	},
	SkipToxicityCheck: true,
	Args: []command_arguments.Arg{
		command_arguments.String{
			Name: setPrefixArgName,
		},
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		prefixArg := parseCtx.ArgsParser.Get(setPrefixArgName)
		if prefixArg == nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(
					ctx,
					locales.Translations.Commands.Prefix.Errors.Required,
				),
			}
		}

		if utf8.RuneCountInString(prefixArg.String()) > 10 {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(
					ctx,
					locales.Translations.Commands.Prefix.Errors.TooLong,
				),
			}
		}

		currentPrefix, err := parseCtx.Services.CommandsPrefixRepository.GetByChannelID(
			ctx,
			parseCtx.Channel.ID,
		)
		if err != nil && !errors.Is(err, channelscommandsprefixrepository.ErrNotFound) {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(
					ctx,
					locales.Translations.Commands.Prefix.Errors.CannotGetCurrent,
				),
				Err: err,
			}
		}

		if currentPrefix == channelscommandsprefixmodel.Nil {
			_, err = parseCtx.Services.CommandsPrefixRepository.Create(
				ctx,
				channelscommandsprefixrepository.CreateInput{
					ChannelID: parseCtx.Channel.ID,
					Prefix:    prefixArg.String(),
				},
			)
			if err != nil {
				return nil, &types.CommandHandlerError{
					Message: i18n.GetCtx(
						ctx,
						locales.Translations.Commands.Prefix.Errors.CannotCreate,
					),
					Err: err,
				}
			}
		} else {
			_, err = parseCtx.Services.CommandsPrefixRepository.Update(
				ctx,
				currentPrefix.ID,
				channelscommandsprefixrepository.UpdateInput{
					Prefix: prefixArg.String(),
				},
			)
			if err != nil {
				return nil, &types.CommandHandlerError{
					Message: i18n.GetCtx(
						ctx,
						locales.Translations.Commands.Prefix.Errors.CannotUpdate,
					),
					Err: err,
				}
			}
		}

		parseCtx.Services.CommandsPrefixCache.Invalidate(ctx, parseCtx.Channel.ID)

		return &types.CommandsHandlerResult{
			Result: []string{
				i18n.GetCtx(
					ctx,
					locales.Translations.Commands.Prefix.Success.Updated,
				),
			},
		}, nil
	},
}
