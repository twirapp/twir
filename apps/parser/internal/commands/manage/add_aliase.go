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
	"gorm.io/gorm"

	model "github.com/twirapp/twir/libs/gomodels"

	"github.com/samber/lo"
)

var AddAliaseCommand = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "commands aliases add",
		Description: null.StringFrom("Add aliase to command"),
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
		aliase := strings.ReplaceAll(
			parseCtx.ArgsParser.Get(commandAliaseArgName).String(),
			"!",
			"",
		)
		aliase = strings.ToLower(aliase)

		var existedCommands []*model.ChannelsCommands
		err := parseCtx.Services.Gorm.
			WithContext(ctx).
			Where(`"channelId" = ?`, parseCtx.Channel.ID).
			Select(`"channelId"`, "name", "aliases").
			Find(&existedCommands).
			Error
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: "cannot get existed commands",
				Err:     err,
			}
		}

		existsError := fmt.Sprintf(`command with "%s" name or aliase already exists`, aliase)
		for _, c := range existedCommands {
			if c.Name == aliase {
				result.Result = append(result.Result, existsError)
				return result, nil
			}

			if lo.Contains(c.Aliases, aliase) {
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
				result.Result = append(result.Result, "Command not found.")
				return result, nil
			} else {
				return nil, &types.CommandHandlerError{
					Message: "cannot get command",
					Err:     err,
				}
			}
		}

		cmd.Aliases = append(cmd.Aliases, aliase)

		err = parseCtx.Services.Gorm.
			WithContext(ctx).
			Save(&cmd).Error

		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: "cannot update command aliases",
				Err:     err,
			}
		}

		parseCtx.Services.CommandsCache.Invalidate(ctx, parseCtx.Channel.ID)

		result.Result = append(result.Result, "✅ Aliase added.")
		return result, nil
	},
}
