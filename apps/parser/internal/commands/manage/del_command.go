package manage

import (
	"context"
	"fmt"
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
		Visible:     true,
		Module:      lo.ToPtr("MANAGE"),
	},
	Handler: func(ctx variables_cache.ExecutionContext) []string {
		if ctx.Text == nil {
			return []string{incorrectUsage}
		}

	
		var cmd *model.ChannelsCommands = nil
		err := ctx.Services.Db.Where(`"channelId" = ? AND name = ?`, ctx.ChannelId, *ctx.Text).First(&cmd).Error
		
		if err != nil || cmd == nil {
			return []string{"Command not found."}
		}

		if cmd.Default {
			return []string{"Cannot delete default command."}
		}

		ctx.Services.Db.
			Where(`"channelId" = ? AND name = ?`, ctx.ChannelId, *ctx.Text).
			Delete(&model.ChannelsCommands{})

		ctx.Services.Redis.Del(
			context.TODO(), 
			fmt.Sprintf("nest:cache:v1/channels/%s/commands", ctx.ChannelId), 
		)

		ctx.Services.Redis.Del(context.TODO(), fmt.Sprintf("commands:%s:%s", ctx.ChannelId, *ctx.Text))

		return []string{"✅ Command removed."}
	},
}