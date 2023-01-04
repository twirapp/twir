package manage

import (
	"fmt"
	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/parser/internal/di"
	"gorm.io/gorm"
	"log"
	"strings"

	"github.com/satont/tsuwari/apps/parser/internal/types"

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
		db := do.MustInvoke[gorm.DB](di.Provider)

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

		commandName := strings.ToLower(strings.ReplaceAll(args[0], "!", ""))
		aliase := strings.ToLower(strings.ReplaceAll(strings.Join(args[1:], " "), "!", ""))

		existedCommands := []model.ChannelsCommands{}
		err := db.Where(`"channelId" = ?`, ctx.ChannelId).Select(`"channelId"`, "name", "aliases").Find(&existedCommands).Error
		if err != nil {
			fmt.Println("cannot get count", err)
			result.Result = append(result.Result, "internal error")
			return result
		}

		existsError := fmt.Sprintf(`command with "%s" name or aliase already exists`, aliase)
		for _, c := range existedCommands {
			if c.Name == aliase {
				result.Result = append(result.Result, existsError)
				return result
			}

			if lo.Contains(c.Aliases, aliase) {
				result.Result = append(result.Result, existsError)
				return result
			}
		}

		cmd := model.ChannelsCommands{}
		err = db.
			Where(`"channelId" = ? AND name = ?`, ctx.ChannelId, commandName).
			Preload(`Responses`).
			First(&cmd).Error

		if err != nil || cmd.ID == "" {
			result.Result = append(result.Result, "Command not found.")
			return result
		}

		cmd.Aliases = append(cmd.Aliases, aliase)

		err = db.
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
