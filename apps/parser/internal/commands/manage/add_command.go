package manage

import (
	"log"
	"strings"
	model "tsuwari/models"
	"tsuwari/parser/internal/types"
	"tsuwari/parser/pkg/helpers"

	variables_cache "tsuwari/parser/internal/variablescache"

	"github.com/guregu/null"
	"github.com/samber/lo"
	uuid "github.com/satori/go.uuid"
)

const (
	exampleUsage   = "!commands add name response"
	incorrectUsage = "Incorrect usage of command. Example: " + exampleUsage
	wentWrong      = "Something went wrong on creating command"
	alreadyExists  = "Command with that name or aliase already exists."
)

var AddCommand = types.DefaultCommand{
	Command: types.Command{
		Name:        "commands add",
		Description: lo.ToPtr("Add command"),
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

		name := args[0]
		text := strings.Join(args[1:], " ")

		if len(name) > 20 {
			result.Result = append(result.Result, "Command name cannot be greatest then 20.")
			return result
		}

		commands := []model.ChannelsCommands{}
		err := ctx.Services.Db.Model(&model.ChannelsCommands{}).
			Where(`"channelId" = ?`, ctx.ChannelId).
			Find(&commands).Error
		if err != nil {
			log.Fatalln(err)
			return nil
		}

		for _, c := range commands {
			if c.Name == name {
				result.Result = append(result.Result, alreadyExists)
				return result
			}

			if helpers.Contains(c.Aliases, name) {
				result.Result = append(result.Result, alreadyExists)
				return result
			}
		}

		commandID := uuid.NewV4().String()
		command := model.ChannelsCommands{
			ID:           commandID,
			Name:         name,
			CooldownType: "GLOBAL",
			Enabled:      true,
			Cooldown:     null.IntFrom(5),
			Aliases:      []string{},
			Description:  null.String{},
			DefaultName:  null.String{},
			Visible:      true,
			ChannelID:    ctx.ChannelId,
			Permission:   "VIEWER",
			Default:      false,
			Module:       "CUSTOM",
			Responses: []model.ChannelsCommandsResponses{
				{
					ID:        uuid.NewV4().String(),
					Text:      null.StringFrom(text),
					CommandID: commandID,
				},
			},
		}
		err = ctx.Services.Db.Create(&command).Error

		if err != nil {
			log.Fatalln(err)
			result.Result = append(result.Result, wentWrong)
			return result
		}

		result.Result = []string{"✅ Command added."}
		return result
	},
}
