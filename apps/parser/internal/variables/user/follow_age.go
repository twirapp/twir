package user

import (
	"context"
	"time"

	"github.com/samber/lo"
	"github.com/satont/twir/apps/parser/internal/types"
	"github.com/satont/twir/apps/parser/pkg/helpers"
)

var FollowAge = &types.Variable{
	Name:         "user.followage",
	Description:  lo.ToPtr(`User followage duration in "1y 3mo 22d" format`),
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
		user, err := parseCtx.Cacher.GetTwitchUserById(ctx, targetUserId)
		if err != nil {
			return nil, err
		}

		var followedAt *time.Time
		if user == nil {
			result.Result = "Cannot find user on twitch."
			return result, nil
		} else if parseCtx.Channel.ID == user.ID {
			followedAt = &user.CreatedAt.Time
		} else {
			follow := parseCtx.Cacher.GetTwitchUserFollow(ctx, user.ID)
			if follow != nil {
				followedAt = &follow.Followed.Time
			}
		}

		if followedAt == nil {
			result.Result = "not a follower"
		} else {
			result.Result = helpers.Duration(
				*followedAt,
				&helpers.DurationOpts{
					UseUtc: true,
					Hide: helpers.DurationOptsHide{
						Minutes: false,
						Seconds: true,
					},
				},
			)
		}

		return result, nil
	},
}
