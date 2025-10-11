package manage

import (
	"context"
	"errors"
	"strings"

	"github.com/guregu/null"
	"github.com/lib/pq"
	command_arguments "github.com/twirapp/twir/apps/parser/internal/command-arguments"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"
	"gorm.io/gorm"

	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/i18n"
)

var EditCommand = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "commands edit",
		Description: null.StringFrom("Edit command response"),
		RolesIDS:    pq.StringArray{model.ChannelRoleTypeModerator.String()},
		Module:      "MANAGE",
		IsReply:     true,
	},
	SkipToxicityCheck: true,
	Args: []command_arguments.Arg{
		command_arguments.String{
			Name: commandNameArgName,
		},
		command_arguments.VariadicString{
			Name: commandTextArgName,
		},
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		result := &types.CommandsHandlerResult{
			Result: make([]string, 0),
		}

		name := strings.ToLower(
			strings.ReplaceAll(
				parseCtx.ArgsParser.Get(commandNameArgName).String(),
				"!",
				"",
			),
		)
		text := parseCtx.ArgsParser.Get(commandTextArgName).String()

		cmd := model.ChannelsCommands{}
		err := parseCtx.Services.Gorm.
			WithContext(ctx).
			Where(`"channelId" = ? AND name = ?`, parseCtx.Channel.ID, name).
			Preload(`Responses`).
			First(&cmd).Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				result.Result = append(result.Result, i18n.GetCtx(ctx, locales.Translations.Commands.Manage.Errors.CommandNotFound))
				return result, nil
			} else {
				return nil, &types.CommandHandlerError{
					Message: i18n.GetCtx(ctx, locales.Translations.Commands.Manage.Errors.CommandCannotGet),
					Err:     err,
				}
			}
		}

		if cmd.Default {
			result.Result = append(result.Result, i18n.GetCtx(ctx, locales.Translations.Commands.Manage.Errors.CommandCannotDeleteDefault))
			return result, nil
		}

		if cmd.Responses != nil && len(cmd.Responses) > 1 {
			result.Result = append(
				result.Result,
				i18n.GetCtx(ctx, locales.Translations.Commands.Manage.Errors.CommandCannotUpdateResponse),
			)
			return result, nil
		}

		err = parseCtx.Services.Gorm.
			WithContext(ctx).
			Model(&model.ChannelsCommandsResponses{}).
			Where(`"commandId" = ?`, cmd.ID).
			Update(`text`, text).Error

		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(ctx, locales.Translations.Commands.Manage.Errors.CommandCannotUpdate),
				Err:     err,
			}
		}

		parseCtx.Services.CommandsCache.Invalidate(ctx, parseCtx.Channel.ID)

		result.Result = append(result.Result, i18n.GetCtx(ctx, locales.Translations.Commands.Manage.Edit.CommandEdited))
		return result, nil
	},
}
