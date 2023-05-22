package manage

import (
	"context"
	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/satont/tsuwari/apps/parser/internal/types"

	"log"
	"strings"

	model "github.com/satont/tsuwari/libs/gomodels"
)

var EditCommand = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "commands edit",
		Description: null.StringFrom("Edit command response"),
		RolesIDS:    pq.StringArray{model.ChannelRoleTypeModerator.String()},
		Module:      "MANAGE",
		IsReply:     true,
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) *types.CommandsHandlerResult {
		result := &types.CommandsHandlerResult{
			Result: make([]string, 0),
		}

		if parseCtx.Text == nil {
			result.Result = append(result.Result, incorrectUsage)
			return result
		}

		args := strings.Split(*parseCtx.Text, " ")

		if len(args) < 2 {
			result.Result = append(result.Result, incorrectUsage)
			return result
		}

		name := strings.ToLower(strings.ReplaceAll(args[0], "!", ""))
		text := strings.Join(args[1:], " ")

		cmd := model.ChannelsCommands{}
		err := parseCtx.Services.Gorm.
			WithContext(ctx).
			Where(`"channelId" = ? AND name = ?`, parseCtx.Channel.ID, name).
			Preload(`Responses`).
			First(&cmd).Error

		if err != nil || cmd.ID == "" {
			log.Fatalln(err)
			result.Result = append(result.Result, "Command not found.")
			return result
		}

		if cmd.Default {
			result.Result = append(result.Result, "Cannot delete default command.")
			return result
		}

		if cmd.Responses != nil && len(cmd.Responses) > 1 {
			result.Result = append(
				result.Result,
				"Cannot update response because you have more then 1 responses in command. Please use UI.",
			)
			return result
		}

		err = parseCtx.Services.Gorm.
			WithContext(ctx).
			Model(&model.ChannelsCommandsResponses{}).
			Where(`"commandId" = ?`, cmd.ID).
			Update(`text`, text).Error

		if err != nil {
			log.Fatalln(err)
			result.Result = append(
				result.Result,
				"Cannot update command response. This is internal bug, please report it.",
			)
			return result
		}

		result.Result = append(result.Result, "✅ Command edited.")
		return result
	},
}
