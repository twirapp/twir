package commandslist

import (
	"strings"
	model "tsuwari/models"
	"tsuwari/parser/internal/types"
	variables_cache "tsuwari/parser/internal/variablescache"

	"github.com/samber/lo"
)

var Variable = types.Variable{
	Name:        "commands.list",
	Description: lo.ToPtr("Command list"),
	Handler: func(ctx *variables_cache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
		cmds := []model.ChannelsCommands{}
		err := ctx.Services.Db.
			Model(&model.ChannelsCommands{}).
			Select("enabled", "visible", "name").
			Find(&cmds).Error
		if err != nil {
			return nil, err
		}

		commandNames := make([]string, len(cmds))
		for _, c := range cmds {
			if c.Enabled && c.Visible {
				commandNames = append(commandNames, c.Name)
			}
		}

		r := types.VariableHandlerResult{
			Result: strings.Join(commandNames, ", "),
		}

		return &r, nil
	},
}
