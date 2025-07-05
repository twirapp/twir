package user

import (
	"context"

	"github.com/samber/lo"
	"github.com/satont/twir/apps/parser/internal/types"
)

var ID = &types.Variable{
	Name:        "user.id",
	Description: lo.ToPtr("User or mentioned user twitch id"),
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

		result.Result = targetUserId

		return &result, nil
	},
}
