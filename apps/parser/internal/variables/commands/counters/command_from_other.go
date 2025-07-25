package command_counters

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
	model "github.com/twirapp/twir/libs/gomodels"
	channelscommandsusages "github.com/twirapp/twir/libs/repositories/channels_commands_usages"
)

var CommandFromOtherCounter = &types.Variable{
	Name:        "command.counter.fromother",
	Description: lo.ToPtr("Counter saying how many times OTHER command was used"),
	Example:     lo.ToPtr("command.counter.fromother|commandName"),
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
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

		commandUUID, err := uuid.Parse(cmd.ID)
		if err != nil {
			parseCtx.Services.Logger.Sugar().Error(err)
			result.Result = "cannot get count"
			return result, nil
		}

		count, err := parseCtx.Services.ChannelsCommandsUsagesRepo.Count(
			ctx, channelscommandsusages.CountInput{
				ChannelID: &parseCtx.Channel.ID,
				CommandID: &commandUUID,
			},
		)
		if err != nil {
			parseCtx.Services.Logger.Sugar().Error(err)
			result.Result = "cannot get count"
			return result, nil
		}

		result.Result = fmt.Sprint(count)

		return result, nil
	},
}
