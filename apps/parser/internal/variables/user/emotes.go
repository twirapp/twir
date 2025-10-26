package user

import (
	"context"
	"strconv"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"
	"github.com/twirapp/twir/libs/i18n"
)

var Emotes = &types.Variable{
	Name:         "user.emotes",
	Description:  lo.ToPtr("User used emotes count"),
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

		var emotes int

		if targetUserId == parseCtx.Sender.ID {
			result.Result = strconv.Itoa(parseCtx.Sender.UserChannelStats.Emotes)
		} else {
			dbUser := parseCtx.Cacher.GetGbUserStats(ctx, targetUserId)
			if dbUser != nil {
				emotes = dbUser.Emotes
			}
		}

		result.Result = i18n.GetCtx(
			ctx,
			locales.Translations.Variables.User.Info.Emotes.
				SetVars(locales.KeysVariablesUserInfoEmotesVars{UserEmotes: emotes}),
		)

		return result, nil
	},
}
