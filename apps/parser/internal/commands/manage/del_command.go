package manage

import (
	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/parser/internal/di"
	"gorm.io/gorm"
	"strings"

	"github.com/satont/tsuwari/apps/parser/internal/types"

	model "github.com/satont/tsuwari/libs/gomodels"

	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"

	"github.com/samber/lo"
)

var DelCommand = types.DefaultCommand{
	Command: types.Command{
		Name:        "commands remove",
		Description: lo.ToPtr("Remove command"),
		RolesNames:  []model.ChannelRoleEnum{model.ChannelRoleTypeModerator},
		Visible:     false,
		Module:      lo.ToPtr("MANAGE"),
		IsReply:     true,
	},
	Handler: func(ctx variables_cache.ExecutionContext) *types.CommandsHandlerResult {
		db := do.MustInvoke[gorm.DB](di.Provider)

		result := &types.CommandsHandlerResult{
			Result: make([]string, 0),
		}

		if ctx.Text == nil {
			result.Result = append(result.Result, incorrectUsage)
			return result
		}

		name := strings.ToLower(strings.ReplaceAll(*ctx.Text, "!", ""))

		var cmd *model.ChannelsCommands = nil
		err := db.Where(`"channelId" = ? AND name = ?`, ctx.ChannelId, name).
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

		db.
			Where(`"channelId" = ? AND name = ?`, ctx.ChannelId, name).
			Delete(&model.ChannelsCommands{})

		result.Result = append(result.Result, "✅ Command removed.")
		return result
	},
}
