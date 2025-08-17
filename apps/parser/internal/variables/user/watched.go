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

		var watched int64 = 0

		if targetUserId == parseCtx.Sender.ID {
			watched = parseCtx.Sender.UserChannelStats.Watched
		} else {
			dbUser := parseCtx.Cacher.GetGbUserStats(ctx, targetUserId)
			if dbUser != nil {
				watched = dbUser.Watched
			} else {
				watched = 0
			}
		}

		watchedD := time.Duration(watched) * time.Millisecond

		result.Result = fmt.Sprintf("%.1fh", watchedD.Hours())

		return &result, nil
	},
}
