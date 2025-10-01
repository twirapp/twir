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
)

const (
	checkAliasesArgName = "command"
)

var CheckAliasesCommand = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "commands aliases",
		Description: null.StringFrom("Check command aliases"),
		RolesIDS:    pq.StringArray{model.ChannelRoleTypeModerator.String()},
		Module:      "MANAGE",
		IsReply:     true,
	},
	Args: []command_arguments.Arg{
		command_arguments.String{
			Name: checkAliasesArgName,
		},
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		result := &types.CommandsHandlerResult{
			Result: make([]string, 0),
		}

		commandName := parseCtx.ArgsParser.Get(checkAliasesArgName).String()
		commandName = strings.ReplaceAll(strings.ToLower(commandName), "!", "")

		cmd := model.ChannelsCommands{}
		err := parseCtx.Services.Gorm.
			WithContext(ctx).
			Where(`"channelId" = ? AND "name" = ?`, parseCtx.Channel.ID, commandName).
			Find(&cmd).
			Error
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(ctx, locales.Translations.Errors.Generic.CannotGetCommand),
				Err:     err,
			}
		}

		if cmd.ID == "" {
			result.Result = append(result.Result, i18n.GetCtx(ctx, locales.Translations.Commands.Manage.Errors.CommandWithNameCannotFind))
			return result, nil
		}

		if len(cmd.Aliases) == 0 {
			result.Result = append(result.Result, i18n.GetCtx(ctx, locales.Translations.Commands.Manage.Errors.CommandHaveNoAliases))
			return result, nil
		}

		parseCtx.Services.CommandsCache.Invalidate(ctx, parseCtx.Channel.ID)

		result.Result = append(result.Result, strings.Join(cmd.Aliases, ", "))
		return result, nil
	},
}
