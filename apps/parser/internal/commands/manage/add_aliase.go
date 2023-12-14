package manage

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/satont/twir/apps/parser/internal/types"
	"gorm.io/gorm"

	model "github.com/satont/twir/libs/gomodels"

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
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		result := &types.CommandsHandlerResult{
			Result: make([]string, 0),
		}

		if parseCtx.Text == nil {
			result.Result = append(result.Result, incorrectUsage)
			return result, nil
		}

		args := strings.Split(*parseCtx.Text, " ")

		if len(args) < 2 {
			result.Result = append(result.Result, incorrectUsage)
			return result, nil
		}

		commandName := strings.ToLower(strings.ReplaceAll(args[0], "!", ""))
		aliase := strings.ToLower(strings.ReplaceAll(strings.Join(args[1:], " "), "!", ""))

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

		result.Result = append(result.Result, "✅ Aliase added.")
		return result, nil
	},
}
