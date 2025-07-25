package user

import (
	"context"
	"strconv"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
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
		dbUser := parseCtx.Cacher.GetGbUserStats(ctx, targetUserId)
		if dbUser != nil {
			result.Result = strconv.Itoa(int(dbUser.Messages))
		} else {
			result.Result = "0"
		}

		return &result, nil
	},
}
