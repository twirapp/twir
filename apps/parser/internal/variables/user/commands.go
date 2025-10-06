package user

import (
	"context"
	"fmt"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"
	"github.com/twirapp/twir/libs/i18n"
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

			result.Result = i18n.GetCtx(ctx, locales.Translations.Errors.Generic.Internal)
			return result, nil
		}

		result.Result = fmt.Sprint(count)

		return result, nil
	},
}
