package manage

import (
	"context"
	"strings"

	"github.com/guregu/null"
	"github.com/lib/pq"
	command_arguments "github.com/twirapp/twir/apps/parser/internal/command-arguments"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"

	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/i18n"

	"github.com/samber/lo"
)

var RemoveAliaseCommand = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "commands aliases remove",
		Description: null.StringFrom("Remove alias from command"),
		RolesIDS:    pq.StringArray{model.ChannelRoleTypeModerator.String()},
		Module:      "MANAGE",
		IsReply:     true,
	},
	Args: []command_arguments.Arg{
		command_arguments.String{
			Name: commandNameArgName,
		},
		command_arguments.String{
			Name: commandAliaseArgName,
		},
	},
	SkipToxicityCheck: true,
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		result := &types.CommandsHandlerResult{
			Result: make([]string, 0),
		}

		commandName := strings.ReplaceAll(
			parseCtx.ArgsParser.Get(commandNameArgName).String(),
			"!",
			"",
		)
		commandName = strings.ToLower(commandName)
		alias := strings.ReplaceAll(
			parseCtx.ArgsParser.Get(commandAliaseArgName).String(),
			"!",
			"",
		)
		alias = strings.ToLower(alias)

		cmd := model.ChannelsCommands{}
		err := parseCtx.Services.Gorm.
			WithContext(ctx).
			Where(`"channelId" = ? AND name = ?`, parseCtx.Channel.ID, commandName).
			First(&cmd).Error

		if err != nil || cmd.ID == "" {
			result.Result = append(result.Result, i18n.GetCtx(ctx, locales.Translations.Commands.Manage.Errors.CommandNotFound))
			return result, nil
		}

		if !lo.Contains(cmd.Aliases, alias) {
			result.Result = append(result.Result, i18n.GetCtx(ctx, locales.Translations.Commands.Manage.Errors.AliasNotCommand))
			return result, nil
		}

		index := lo.IndexOf(cmd.Aliases, alias)
		cmd.Aliases = append(cmd.Aliases[:index], cmd.Aliases[index+1:]...)

		err = parseCtx.Services.Gorm.
			WithContext(ctx).
			Save(&cmd).Error

		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(ctx, locales.Translations.Commands.Manage.Errors.CommandCannotSave),
				Err:     err,
			}
		}

		parseCtx.Services.CommandsCache.Invalidate(ctx, parseCtx.Channel.ID)

		result.Result = append(result.Result, i18n.GetCtx(ctx, locales.Translations.Commands.Manage.Remove.AliasRemoved))
		return result, nil
	},
}
