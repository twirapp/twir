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
		Visible:     true,
		Module:      lo.ToPtr("MANAGE"),
	},
	Handler: func(ctx variables_cache.ExecutionContext) []string {
		if ctx.Text == nil {
			return []string{incorrectUsage}
		}

		args := strings.Split(*ctx.Text, " ")

		if len(args) < 2 {
			return []string{incorrectUsage}
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
			return []string{"Command not found."}
		}

		if cmd.Default {
			return []string{"Cannot delete default command."}
		}

		if cmd.Responses != nil && len(cmd.Responses) > 1 {
			return []string{"Cannot update response because you have more then 1 responses in command. Please use UI."}
		}

		err = ctx.Services.Db.
			Model(&model.ChannelsCommandsResponses{}).
			Where(`"commandId" = ?`, cmd.ID).
			Update(`text`, text).Error

		if err != nil {
			log.Fatalln(err)
			return []string{"Cannot update command response. This is internal bug, please report it."}
		}

		bytes, err := CreateRedisBytes(cmd, text, nil)

		if err != nil {
			log.Fatalln(err)
			return []string{"Cannot update command response. This is internal bug, please report it."}
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

		return []string{"✅ Command edited."}
	},
}