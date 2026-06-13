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

var Age = &types.Variable{
	Name:         "user.age",
	Description:  lo.ToPtr("User account age"),
	CommandsOnly: true,
	Handler: shared.HandlerByPlatform(map[platformentity.Platform]types.VariableHandler{
		shared.PlatformTwitch: func(
			ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
		) (*types.VariableHandlerResult, error) {
			result := types.VariableHandlerResult{}

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

			if user == nil {
				result.Result = i18n.GetCtx(ctx, locales.Translations.Variables.User.Errors.FindUserOnTwitch)
			} else {
				result.Result = helpers.Duration(
					user.CreatedAt.Time,
					&helpers.DurationOpts{
						UseUtc: true,
						Hide: helpers.DurationOptsHide{
							Seconds: true,
						},
					},
				)
			}

			return &result, nil
		},
		shared.PlatformKick: func(
			ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
		) (*types.VariableHandlerResult, error) {
			result := types.VariableHandlerResult{}

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
				result.Result = i18n.GetCtx(ctx, locales.Translations.Variables.User.Errors.FindUserOnTwitch)
				return &result, nil
			}

			createdAt, err := time.Parse(time.RFC3339Nano, channelUser.CreatedAt)
			if err != nil {
				return nil, err
			}

			result.Result = helpers.Duration(
				createdAt,
				&helpers.DurationOpts{
					UseUtc: true,
					Hide: helpers.DurationOptsHide{
						Seconds: true,
					},
				},
			)

			return &result, nil
		},
	}),
}
