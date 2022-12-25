package commandslist

import (
	"strings"

	"github.com/satont/tsuwari/apps/parser/internal/types"
	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"

	model "github.com/satont/tsuwari/libs/gomodels"

	"github.com/samber/lo"
)

var Variable = types.Variable{
	Name:        "commands.list",
	Description: lo.ToPtr("Command list"),
	Handler: func(ctx *variables_cache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
		cmds := []model.ChannelsCommands{}
		err := ctx.Services.Db.
			Model(&model.ChannelsCommands{}).
			Where(`"channelId" = ?`, ctx.ChannelId).
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
		commandNames = lo.Filter(commandNames, func(n string, _ int) bool {
			return n != ""
		})

		r := types.VariableHandlerResult{
			Result: strings.Join(commandNames, ", "),
		}

		return &r, nil
	},
}
