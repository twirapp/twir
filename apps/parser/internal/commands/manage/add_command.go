package manage

import (
	"context"
	"strings"

	"github.com/lib/pq"
	"github.com/samber/lo"
	command_arguments "github.com/satont/twir/apps/parser/internal/command-arguments"
	"github.com/satont/twir/apps/parser/internal/types"

	model "github.com/satont/twir/libs/gomodels"

	"github.com/guregu/null"
	uuid "github.com/satori/go.uuid"
)

var AddCommand = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "commands add",
		Description: null.StringFrom("Add command"),
		RolesIDS:    pq.StringArray{model.ChannelRoleTypeModerator.String()},
		Module:      "MANAGE",
		IsReply:     true,
	},
	Args: []command_arguments.Arg{
		command_arguments.String{
			Name: commandNameArgName,
		},
		command_arguments.VariadicString{
			Name: commandTextArgName,
		},
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		result := &types.CommandsHandlerResult{
			Result: make([]string, 0),
		}

		name := strings.ToLower(
			strings.ReplaceAll(
				parseCtx.ArgsParser.Get(commandNameArgName).String(),
				"!",
				"",
			),
		)
		text := parseCtx.ArgsParser.Get(commandTextArgName).String()

		if len(name) > 20 {
			result.Result = append(result.Result, "Command name cannot be greatest then 20.")
			return result, nil
		}

		var commands []*model.ChannelsCommands
		err := parseCtx.Services.Gorm.
			WithContext(ctx).
			Model(&model.ChannelsCommands{}).
			Where(`"channelId" = ?`, parseCtx.Channel.ID).
			Find(&commands).Error
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: "cannot get existed commands",
				Err:     err,
			}
		}

		for _, c := range commands {
			if c.Name == name {
				result.Result = append(result.Result, alreadyExists)
				return result, nil
			}

			if lo.Contains(c.Aliases, name) {
				result.Result = append(result.Result, alreadyExists)
				return result, nil
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
			ChannelID:    parseCtx.Channel.ID,
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
		err = parseCtx.Services.Gorm.WithContext(ctx).Create(&command).Error

		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: "cannot create command",
				Err:     err,
			}
		}

		result.Result = []string{"✅ Command added."}
		return result, nil
	},
}
