package manage

import (
	"context"
	"fmt"
	"log"
	"strings"
	model "tsuwari/models"
	"tsuwari/parser/internal/types"

	variables_cache "tsuwari/parser/internal/variablescache"

	"github.com/samber/lo"
)

var EditCommand = types.DefaultCommand{
	Command: types.Command{
		Name:        "commands edit",
		Description: lo.ToPtr("Edit command response"),
		Permission:  "MODERATOR",
		Visible:     false,
		Module:      lo.ToPtr("MANAGE"),
	},
	Handler: func(ctx variables_cache.ExecutionContext) *types.CommandsHandlerResult {
		result := &types.CommandsHandlerResult{
			Result: make([]string, 0),
		}

		if ctx.Text == nil {
			result.Result = append(result.Result, incorrectUsage)
			return result
		}

		args := strings.Split(*ctx.Text, " ")

		if len(args) < 2 {
			result.Result = append(result.Result, incorrectUsage)
			return result
		}

		name := args[0]
		text := strings.Join(args[1:], " ")

		cmd := model.ChannelsCommands{}
		err := ctx.Services.Db.
			Where(`"channelId" = ? AND name = ?`, ctx.ChannelId, name).
			Preload(`Responses`).
			First(&cmd).Error
		
		if err != nil || &cmd == nil {
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

		err = ctx.Services.Db.
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

		bytes, err := CreateRedisBytes(cmd, text, nil)

		if err != nil {
			log.Fatalln(err)
			result.Result = append(
				result.Result, 
				"Cannot update command response. This is internal bug, please report it.",
			)
			return result
		}

		ctx.Services.Redis.Set(
			context.TODO(),
			fmt.Sprintf("commands:%s:%s", ctx.ChannelId, cmd.Name),
			*bytes,
			0,
		)

		ctx.Services.Redis.Del(
			context.TODO(), 
			fmt.Sprintf("nest:cache:v1/channels/%s/commands", ctx.ChannelId), 
		)

		result.Result = append(result.Result, "✅ Command edited.")
		return result
	},
}