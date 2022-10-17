package manage

import (
	model "tsuwari/models"
	"tsuwari/parser/internal/types"

	variables_cache "tsuwari/parser/internal/variablescache"

	"github.com/samber/lo"
)

var DelCommand = types.DefaultCommand{
	Command: types.Command{
		Name:        "commands remove",
		Description: lo.ToPtr("Remove command"),
		Permission:  "MODERATOR",
		Visible:     false,
		Module:      lo.ToPtr("MANAGE"),
		IsReply:     true,
	},
	Handler: func(ctx variables_cache.ExecutionContext) *types.CommandsHandlerResult {
		result := &types.CommandsHandlerResult{
			Result: make([]string, 0),
		}

		if ctx.Text == nil {
			result.Result = append(result.Result, incorrectUsage)
			return result
		}

		var cmd *model.ChannelsCommands = nil
		err := ctx.Services.Db.Where(`"channelId" = ? AND name = ?`, ctx.ChannelId, *ctx.Text).
			First(&cmd).
			Error

		if err != nil || cmd == nil {
			result.Result = append(result.Result, "Command not found.")
			return result
		}

		if cmd.Default {
			result.Result = append(result.Result, "Cannot delete default command.")
			return result
		}

		ctx.Services.Db.
			Where(`"channelId" = ? AND name = ?`, ctx.ChannelId, *ctx.Text).
			Delete(&model.ChannelsCommands{})

		result.Result = append(result.Result, "✅ Command removed.")
		return result
	},
}
