package user

import (
	"context"
	"strconv"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
)

var Reputation = &types.Variable{
	Name:         "user.reputation",
	Description:  lo.ToPtr("User reputation"),
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
		if targetUserId == parseCtx.Sender.ID {
			result.Result = strconv.Itoa(int(parseCtx.Sender.UserChannelStats.Reputation))
		} else {
			dbUser := parseCtx.Cacher.GetGbUserStats(ctx, targetUserId)
			if dbUser != nil {
				result.Result = strconv.Itoa(int(dbUser.Reputation))
			} else {
				result.Result = "0"
			}
		}

		return &result, nil
	},
}
