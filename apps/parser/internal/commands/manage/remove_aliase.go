package manage

import (
	"log"
	"strings"

	"github.com/satont/tsuwari/apps/parser/internal/types"

	model "github.com/satont/tsuwari/libs/gomodels"

	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"

	"github.com/samber/lo"
)

var RemoveAliaseCommand = types.DefaultCommand{
	Command: types.Command{
		Name:        "commands aliases remove",
		Description: lo.ToPtr("Remove aliase from command"),
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

		args := strings.Split(*ctx.Text, " ")

		if len(args) < 2 {
			result.Result = append(result.Result, incorrectUsage)
			return result
		}

		commandName := strings.ReplaceAll(args[0], "!", "")
		aliase := strings.ReplaceAll(strings.Join(args[1:], " "), "!", "")

		cmd := model.ChannelsCommands{}
		err := ctx.Services.Db.
			Where(`"channelId" = ? AND name = ?`, ctx.ChannelId, commandName).
			First(&cmd).Error

		if err != nil || cmd.ID == "" {
			result.Result = append(result.Result, "Command not found.")
			return result
		}

		if !lo.Contains(cmd.Aliases, aliase) {
			result.Result = append(result.Result, "That aliase not in the command")
			return result
		}

		index := lo.IndexOf(cmd.Aliases, aliase)
		cmd.Aliases = append(cmd.Aliases[:index], cmd.Aliases[index+1:]...)

		err = ctx.Services.Db.
			Save(&cmd).Error

		if err != nil {
			log.Fatalln(err)
			result.Result = append(
				result.Result,
				"Cannot update command aliases. This is internal bug, please report it.",
			)
			return result
		}

		result.Result = append(result.Result, "✅ Aliase removed.")
		return result
	},
}
