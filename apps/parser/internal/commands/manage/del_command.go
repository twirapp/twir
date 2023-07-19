package manage

import (
	"context"
	"strings"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/satont/twir/apps/parser/internal/types"

	model "github.com/satont/twir/libs/gomodels"
)

var DelCommand = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "commands remove",
		Description: null.StringFrom("Remove command"),
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

		name := strings.ToLower(strings.ReplaceAll(*parseCtx.Text, "!", ""))

		var cmd *model.ChannelsCommands = nil
		err := parseCtx.Services.Gorm.
			WithContext(ctx).
			Where(`"channelId" = ? AND name = ?`, parseCtx.Channel.ID, name).
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

		parseCtx.Services.Gorm.
			WithContext(ctx).
			Where(`"channelId" = ? AND name = ?`, parseCtx.Channel.ID, name).
			Delete(&model.ChannelsCommands{})

		result.Result = append(result.Result, "✅ Command removed.")
		return result
	},
}
