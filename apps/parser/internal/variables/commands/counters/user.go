package command_counters

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
	channelscommandsusages "github.com/twirapp/twir/libs/repositories/channels_commands_usages"
)

var CommandUserCounter = &types.Variable{
	Name:         "command.counter.user",
	Description:  lo.ToPtr("Counter saying how many times command was used by sender user"),
	CommandsOnly: true,
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}

		commandUUID, err := uuid.Parse(parseCtx.Command.ID)
		if err != nil {
			result.Result = "cannot get count"
			return result, nil
		}

		count, err := parseCtx.Services.ChannelsCommandsUsagesRepo.Count(
			ctx, channelscommandsusages.CountInput{
				ChannelID: &parseCtx.Channel.ID,
				CommandID: &commandUUID,
				UserID:    &parseCtx.Sender.ID,
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
