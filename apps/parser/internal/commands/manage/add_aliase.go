package manage

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/guregu/null"
	"github.com/lib/pq"
	command_arguments "github.com/twirapp/twir/apps/parser/internal/command-arguments"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"
	"gorm.io/gorm"

	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/i18n"

	"github.com/samber/lo"
)

var AddAliaseCommand = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "commands aliases add",
		Description: null.StringFrom("Add alias to command"),
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

		var existedCommands []*model.ChannelsCommands
		err := parseCtx.Services.Gorm.
			WithContext(ctx).
			Where(`"channelId" = ?`, parseCtx.Channel.ID).
			Select(`"channelId"`, "name", "aliases").
			Find(&existedCommands).
			Error
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(ctx, locales.Translations.Commands.Manage.Errors.AliasCannotGetExistedCommands),
				Err:     err,
			}
		}

		existsError := fmt.Sprintf(i18n.GetCtx(
			ctx,
			locales.Translations.Commands.Manage.Errors.AliasAlreadyExist.
				SetVars(locales.KeysCommandsManageErrorsAliasAlreadyExistVars{Alias: alias}),
		))
		for _, c := range existedCommands {
			if c.Name == alias {
				result.Result = append(result.Result, existsError)
				return result, nil
			}

			if lo.Contains(c.Aliases, alias) {
				result.Result = append(result.Result, existsError)
				return result, nil
			}
		}

		cmd := model.ChannelsCommands{}
		err = parseCtx.Services.Gorm.
			WithContext(ctx).
			Where(`"channelId" = ? AND name = ?`, parseCtx.Channel.ID, commandName).
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

		cmd.Aliases = append(cmd.Aliases, alias)

		err = parseCtx.Services.Gorm.
			WithContext(ctx).
			Save(&cmd).Error

		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(ctx, locales.Translations.Commands.Manage.Errors.AliasCannotUpdate),
				Err:     err,
			}
		}

		parseCtx.Services.CommandsCache.Invalidate(ctx, parseCtx.Channel.ID)

		result.Result = append(result.Result, i18n.GetCtx(ctx, locales.Translations.Commands.Manage.Add.AliasAdd))
		return result, nil
	},
}
