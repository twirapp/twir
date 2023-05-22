package command_counters

import (
	"context"
	"fmt"

	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/parser/internal/types"
	model "github.com/satont/tsuwari/libs/gomodels"
)

var CommandFromOtherCounter = &types.Variable{
	Name:        "command.counter.fromother",
	Description: lo.ToPtr("Counter saying how many times OTHER command was used"),
	Example:     lo.ToPtr("command.counter.fromother|commandName"),
	Handler: func(ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}

		if variableData.Params == nil {
			result.Result = "Have not passed params to variable. "
			return result, nil
		}

		cmd := model.ChannelsCommands{}
		err := parseCtx.Services.Gorm.
			WithContext(ctx).
			Where(`"channelId" = ? AND "name" = ?`, parseCtx.Channel.ID, *variableData.Params).
			First(&cmd).Error

		if err != nil || cmd.ID == "" {
			result.Result = fmt.Sprintf(`Command with name "%s" not found`, *variableData.Params)
			return result, nil
		}

		count, err := getCount(parseCtx.Services.Gorm, cmd.ID, nil)
		if err != nil {
			parseCtx.Services.Logger.Sugar().Error(err)

			result.Result = "cannot get count"
			return result, nil
		}
		result.Result = count

		return result, nil
	},
}
