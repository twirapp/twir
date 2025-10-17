package user

import (
	"context"
	"fmt"
	"time"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"
	"github.com/twirapp/twir/libs/i18n"
)

var FollowSince = &types.Variable{
	Name:         "user.followsince",
	Description:  lo.ToPtr(`User follow since in "16 January 2023 (22 days)" format.`),
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
			result.Result = i18n.GetCtx(ctx, locales.Translations.Errors.Generic.CannotFindUserTwitch)
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
			result.Result = i18n.GetCtx(ctx, locales.Translations.Errors.Generic.CannotFindUserTwitch)
		} else {
			result.Result = fmt.Sprintf(
				"%s (%.0f days)",
				followedAt.UTC().Format("2 January 2006"),
				time.Now().UTC().Sub(followedAt.UTC()).Hours()/24,
			)
		}

		return result, nil
	},
}
