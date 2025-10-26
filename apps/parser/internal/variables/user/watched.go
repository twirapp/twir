package user

import (
	"context"
	"fmt"
	"time"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
)

var Watched = &types.Variable{
	Name:         "user.watched",
	Description:  lo.ToPtr("User watched time"),
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

		count := parseCtx.Sender.UserChannelStats.Watched
		if targetUserId != parseCtx.Sender.ID {
			dbUser := parseCtx.Cacher.GetGbUserStats(ctx, targetUserId)
			if dbUser != nil {
				count = dbUser.Watched
			} else {
				count = 0
			}
		}

		watchedD := time.Duration(count) * time.Millisecond

		result.Result = fmt.Sprintf(
			"%.1f",
			watchedD.Hours(),
		)

		return &result, nil
	},
}
