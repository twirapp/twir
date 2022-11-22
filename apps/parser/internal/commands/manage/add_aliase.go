package manage

import (
	"fmt"
	"log"
	"strings"

	"github.com/satont/tsuwari/apps/parser/internal/types"
	"gorm.io/gorm"

	model "github.com/satont/tsuwari/libs/gomodels"

	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"

	"github.com/samber/lo"
)

var AddAliaseCommand = types.DefaultCommand{
	Command: types.Command{
		Name:        "commands aliases add",
		Description: lo.ToPtr("Add aliase to command"),
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
		aliase := strings.ReplaceAll(args[1], "!", "")

		var existedCommandCount int64
		fmt.Println(ctx.Services.Db.ToSQL(func(tx *gorm.DB) *gorm.DB {
			return tx.
				Model(&model.ChannelsCommands{}).
				Where(`"channelId" = ? AND "aliases" @> ?`, ctx.ChannelId, []string{aliase}).
				Count(&existedCommandCount)
		}))
		err := ctx.Services.Db.
			Model(&model.ChannelsCommands{}).
			Where(`"channelId" = ? AND "aliases" @> ?`, ctx.ChannelId, []string{aliase}).
			Count(&existedCommandCount).Error
		if err != nil {
			fmt.Println("cannot get count", err)
			result.Result = append(result.Result, "internal error")
			return result
		}

		if existedCommandCount > 0 {
			result.Result = append(result.Result, "command with that aliase already exists")
			return result
		}

		cmd := model.ChannelsCommands{}
		err = ctx.Services.Db.
			Where(`"channelId" = ? AND name = ?`, ctx.ChannelId, commandName).
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

		cmd.Aliases = append(cmd.Aliases, aliase)

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

		result.Result = append(result.Result, "✅ Aliase added.")
		return result
	},
}
