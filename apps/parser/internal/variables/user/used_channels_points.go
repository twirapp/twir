package user

import (
	"context"
	"strconv"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
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

		formattedCount := strconv.Itoa(count)
		result.Result = formattedCount

		return &result, nil
	},
}
