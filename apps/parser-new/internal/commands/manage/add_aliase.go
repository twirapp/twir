package manage

import (
	"context"
	"fmt"
	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/satont/tsuwari/apps/parser-new/internal/types"
	"log"
	"strings"

	model "github.com/satont/tsuwari/libs/gomodels"

	"github.com/samber/lo"
)

var AddAliaseCommand = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "commands aliases add",
		Description: null.StringFrom("Add aliase to command"),
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

		commandName := strings.ToLower(strings.ReplaceAll(args[0], "!", ""))
		aliase := strings.ToLower(strings.ReplaceAll(strings.Join(args[1:], " "), "!", ""))

		var existedCommands []*model.ChannelsCommands
		err := parseCtx.Services.Gorm.
			WithContext(ctx).
			Where(`"channelId" = ?`, parseCtx.Channel.ID).
			Select(`"channelId"`, "name", "aliases").
			Find(&existedCommands).
			Error
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
		err = parseCtx.Services.Gorm.
			WithContext(ctx).
			Where(`"channelId" = ? AND name = ?`, parseCtx.Channel.ID, commandName).
			Preload(`Responses`).
			First(&cmd).Error

		if err != nil || cmd.ID == "" {
			result.Result = append(result.Result, "Command not found.")
			return result
		}

		cmd.Aliases = append(cmd.Aliases, aliase)

		err = parseCtx.Services.Gorm.
			WithContext(ctx).
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
