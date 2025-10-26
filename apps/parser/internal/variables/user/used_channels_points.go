package user

import (
	"context"
	"strconv"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"
	"github.com/twirapp/twir/libs/i18n"
)

var UsedChannelPoints = &types.Variable{
	Name:         "user.usedChannelPoints",
	Description:  lo.ToPtr("How many channel points user spent on channel"),
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

		var count int

		if targetUserId == parseCtx.Sender.ID {
			result.Result = strconv.Itoa(int(parseCtx.Sender.UserChannelStats.UsedChannelPoints))
		} else {
			dbUser := parseCtx.Cacher.GetGbUserStats(ctx, targetUserId)
			if dbUser != nil {
				count = int(dbUser.UsedChannelPoints)
			}
		}

		result.Result = i18n.GetCtx(
			ctx,
			locales.Translations.Variables.User.Info.Points.
				SetVars(locales.KeysVariablesUserInfoPointsVars{UserPoints: count}),
		)

		return &result, nil
	},
}
