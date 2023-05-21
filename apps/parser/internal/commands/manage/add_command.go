package manage

import (
	"github.com/lib/pq"
	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/parser/internal/di"
	"gorm.io/gorm"
	"log"
	"strings"

	"github.com/satont/tsuwari/apps/parser/internal/types"
	"github.com/satont/tsuwari/apps/parser/pkg/helpers"

	model "github.com/satont/tsuwari/libs/gomodels"

	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"

	"github.com/guregu/null"
	uuid "github.com/satori/go.uuid"
)

const (
	exampleUsage   = "!commands add name response"
	incorrectUsage = "Incorrect usage of command. Example: " + exampleUsage
	wentWrong      = "Something went wrong on creating command"
	alreadyExists  = "Command with that name or aliase already exists."
)

var AddCommand = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "commands add",
		Description: null.StringFrom("Add command"),
		RolesIDS:    pq.StringArray{model.ChannelRoleTypeModerator.String()},
		Module:      "MANAGE",
		IsReply:     true,
	},
	Handler: func(ctx *variables_cache.ExecutionContext) *types.CommandsHandlerResult {
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

		name := strings.ToLower(strings.ReplaceAll(args[0], "!", ""))
		text := strings.Join(args[1:], " ")

		if len(name) > 20 {
			result.Result = append(result.Result, "Command name cannot be greatest then 20.")
			return result
		}

		commands := []model.ChannelsCommands{}
		err := db.Model(&model.ChannelsCommands{}).
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
			Default:      false,
			Module:       "CUSTOM",
			Responses: []*model.ChannelsCommandsResponses{
				{
					ID:        uuid.NewV4().String(),
					Text:      null.StringFrom(text),
					CommandID: commandID,
				},
			},
		}
		err = db.Create(&command).Error

		if err != nil {
			log.Fatalln(err)
			result.Result = append(result.Result, wentWrong)
			return result
		}

		result.Result = []string{"✅ Command added."}
		return result
	},
}
