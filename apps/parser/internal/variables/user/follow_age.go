package user

import (
	"context"
	"strings"
	"time"

	"github.com/nicklaw5/helix/v2"
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

		var targetUser *helix.User
		if parseCtx.Text != nil {
			userName := strings.ReplaceAll(*parseCtx.Text, "@", "")

			user, err := parseCtx.Cacher.GetTwitchUserByName(ctx, userName)
			if err != nil {
				return nil, err
			}

			if user != nil {
				targetUser = user
			}
		} else {
			targetUser = parseCtx.Cacher.GetTwitchSenderUser(ctx)
		}

		var followedAt *time.Time
		if targetUser == nil {
			result.Result = "Cannot find user on twitch."
			return result, nil
		} else if parseCtx.Channel.ID == targetUser.ID {
			followedAt = &targetUser.CreatedAt.Time
		} else {
			follow := parseCtx.Cacher.GetTwitchUserFollow(ctx, targetUser.ID)
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
