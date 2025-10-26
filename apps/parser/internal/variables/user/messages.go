package user

import (
	"context"
	"strconv"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"
	"github.com/twirapp/twir/libs/i18n"
)

var Messages = &types.Variable{
	Name:         "user.messages",
	Description:  lo.ToPtr("User messages"),
	CommandsOnly: true,
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := types.VariableHandlerResult{}

		targetUserId := lo.
			IfF(
				len(parseCtx.Mentions) > 0, func() string {
					return parseCtx.Mentions[0].UserId
				},
			).
			Else(parseCtx.Sender.ID)

		var msgs int

		if targetUserId == parseCtx.Sender.ID {
			result.Result = strconv.Itoa(int(parseCtx.Sender.UserChannelStats.Messages))
		} else {
			dbUser := parseCtx.Cacher.GetGbUserStats(ctx, targetUserId)
			if dbUser != nil {
				msgs = int(dbUser.Messages)
			}
		}

		result.Result = i18n.GetCtx(
			ctx,
			locales.Translations.Variables.User.Info.Messages.
				SetVars(locales.KeysVariablesUserInfoMessagesVars{UserMessages: msgs}),
		)

		return &result, nil
	},
}
