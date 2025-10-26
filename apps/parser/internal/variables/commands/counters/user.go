package command_counters

import (
	"context"
	"fmt"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"
	"github.com/twirapp/twir/libs/i18n"
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

		count, err := parseCtx.Services.ChannelsCommandsUsagesRepo.Count(
			ctx, channelscommandsusages.CountInput{
				ChannelID: &parseCtx.Channel.ID,
				CommandID: &parseCtx.Command.ID,
				UserID:    &parseCtx.Sender.ID,
			},
		)
		if err != nil {
			parseCtx.Services.Logger.Sugar().Error(err)
			result.Result = i18n.GetCtx(ctx, locales.Translations.Variables.Commands.Info.GetCount)
			return result, nil
		}

		result.Result = fmt.Sprint(count)

		return result, nil
	},
}
