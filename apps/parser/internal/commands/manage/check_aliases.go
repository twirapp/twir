package manage

import (
	"fmt"
	"strings"

	"github.com/satont/tsuwari/apps/parser/internal/types"

	model "github.com/satont/tsuwari/libs/gomodels"

	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"

	"github.com/samber/lo"
)

var CheckAliasesCommand = types.DefaultCommand{
	Command: types.Command{
		Name:        "commands aliases",
		Description: lo.ToPtr("Check command aliases"),
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
			result.Result = append(result.Result, "type command name for check aliases.")
			return result
		}

		commandName := strings.ReplaceAll(strings.ToLower(*ctx.Text), "!", "")

		cmd := model.ChannelsCommands{}
		err := ctx.Services.Db.Where(`"channelId" = ? AND "name" = ?`, ctx.ChannelId, commandName).Find(&cmd).Error
		if err != nil {
			fmt.Println(err)
			result.Result = append(result.Result, "internal error")
			return result
		}

		if cmd.ID == "" {
			result.Result = append(result.Result, "command with that name not found.")
			return result
		}

		if len(cmd.Aliases) == 0 {
			result.Result = append(result.Result, "command have no aliases")
			return result
		}

		result.Result = append(result.Result, strings.Join(cmd.Aliases, ", "))
		return result
	},
}
