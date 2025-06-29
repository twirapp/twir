package user

import (
	"context"
	"fmt"

	"github.com/samber/lo"
	"github.com/satont/twir/apps/parser/internal/types"
	channelscommandsusages "github.com/twirapp/twir/libs/repositories/channels_commands_usages"
)

var Commands = &types.Variable{
	Name:         "user.commands",
	Description:  lo.ToPtr("User used commands count"),
	CommandsOnly: true,
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}

		targetUserId := lo.
			IfF(
				len(parseCtx.Mentions) > 0, func() string {
					return parseCtx.Mentions[0].UserId
				},
			).
			Else(parseCtx.Sender.ID)

		count, err := parseCtx.Services.ChannelsCommandsUsagesRepo.Count(
			ctx, channelscommandsusages.CountInput{
				ChannelID: &parseCtx.Channel.ID,
				UserID:    &targetUserId,
			},
		)
		if err != nil {
			parseCtx.Services.Logger.Sugar().Error(err)

			result.Result = "internal error"
			return result, nil
		}

		result.Result = fmt.Sprint(count)

		return result, nil
	},
}
