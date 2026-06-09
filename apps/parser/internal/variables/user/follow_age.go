package user

import (
	"context"
	"time"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/internal/variables/shared"
	"github.com/twirapp/twir/apps/parser/locales"
	"github.com/twirapp/twir/apps/parser/pkg/helpers"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	"github.com/twirapp/twir/libs/i18n"
)

var FollowAge = &types.Variable{
	Name:         "user.followage",
	Description:  lo.ToPtr(`User followage duration in "1y 3mo 22d" format`),
	CommandsOnly: true,
	Handler: shared.HandlerByPlatform(map[platformentity.Platform]types.VariableHandler{
		shared.PlatformTwitch: func(
			ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
		) (*types.VariableHandlerResult, error) {
			result := &types.VariableHandlerResult{}

			targetUserId := lo.
				IfF(
					len(parseCtx.Mentions) > 0, func() string {
						return parseCtx.Mentions[0].UserID
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
				result.Result = i18n.GetCtx(ctx, locales.Translations.Errors.Generic.NotAFollower)
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
		shared.PlatformKick: func(
			ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
		) (*types.VariableHandlerResult, error) {
			result := &types.VariableHandlerResult{}

			channelSlug := parseCtx.Channel.Name
			userSlug := lo.
				IfF(
					len(parseCtx.Mentions) > 0, func() string {
						return parseCtx.Mentions[0].UserLogin
					},
				).
				Else(parseCtx.Sender.Name)

			channelUser, err := shared.GetKickChannelUser(ctx, parseCtx, channelSlug, userSlug)
			if err != nil {
				parseCtx.Services.Logger.Sugar().Error(err)
				result.Result = i18n.GetCtx(ctx, locales.Translations.Errors.Generic.CannotFindUserTwitch)
				return result, nil
			}

			if channelUser.FollowingSince == "" {
				result.Result = i18n.GetCtx(ctx, locales.Translations.Errors.Generic.NotAFollower)
				return result, nil
			}

			followedAt, err := time.Parse(time.RFC3339Nano, channelUser.FollowingSince)
			if err != nil {
				return nil, err
			}

			result.Result = helpers.Duration(
				followedAt,
				&helpers.DurationOpts{
					UseUtc: true,
					Hide: helpers.DurationOptsHide{
						Minutes: false,
						Seconds: true,
					},
				},
			)

			return result, nil
		},
	}),
}
