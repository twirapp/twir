package commands_list

import (
	"context"
	"strings"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
	model "github.com/twirapp/twir/libs/gomodels"
)

var Variable = &types.Variable{
	Name:                "commands.list",
	Description:         lo.ToPtr("Command list"),
	CanBeUsedInRegistry: true,
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		cmds := []*model.ChannelsCommands{}
		err := parseCtx.Services.Gorm.
			WithContext(ctx).
			Model(&model.ChannelsCommands{}).
			Where(`"channelId" = ?`, parseCtx.Channel.ID).
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
		commandNames = lo.Filter(
			commandNames, func(n string, _ int) bool {
				return n != ""
			},
		)

		r := types.VariableHandlerResult{
			Result: strings.Join(commandNames, ", "),
		}

		return &r, nil
	},
}
