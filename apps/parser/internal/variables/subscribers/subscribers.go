package subscribers

import (
	"context"
	"fmt"

	"github.com/kvizyx/twitchy/helix"
	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/internal/variables/shared"
	"github.com/twirapp/twir/apps/parser/locales"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	"github.com/twirapp/twir/libs/i18n"
	"github.com/twirapp/twir/libs/twitch"
	"go.uber.org/zap"
)

var LatestSubscriberUsername = &types.Variable{
	Name:                "subscribers.latest.userName",
	Description:         lo.ToPtr("Latest subscriber username"),
	CanBeUsedInRegistry: true,
	Handler: func(
		ctx context.Context,
		parseCtx *types.VariableParseContext,
		variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}

		if parseCtx.Platform != "twitch" {
			result.Result = "not supported on this platform"
			return result, nil
		}

		twitchClient, err := twitch.NewUserClientWithContext(
			ctx,
			parseCtx.Channel.TwitchUserID,
			*parseCtx.Services.Config,
			parseCtx.Services.Bus,
		)
		if err != nil {
			parseCtx.Services.Logger.Error(i18n.GetCtx(ctx, locales.Translations.Errors.Generic.CannotCreateTwitch), zap.Error(err))
			return result, nil
		}

		first := "1"
		subscribers, err := twitchClient.Subscriptions.GetBroadcasterSubscriptions(
			ctx,
			helix.GetBroadcasterSubscriptionsRequest{
				BroadcasterID: parseCtx.Channel.ID,
				First:         &first,
			},
		)
		if err != nil {
			parseCtx.Services.Logger.Error(i18n.GetCtx(ctx, locales.Translations.Variables.Subscribers.Errors.GetSubscribers), zap.Error(err))
			return result, nil
		}
		if len(subscribers.Data.Subscriptions) == 0 {
			return result, nil
		}

		result.Result = fmt.Sprint(subscribers.Data.Subscriptions[0].UserName)
		return result, nil
	},
}

var Count = &types.Variable{
	Name:                "subscribers.count",
	Description:         lo.ToPtr("Subscribers count"),
	CanBeUsedInRegistry: true,
	Handler: shared.HandlerByPlatform(map[platformentity.Platform]types.VariableHandler{
		shared.PlatformTwitch: func(
			ctx context.Context,
			parseCtx *types.VariableParseContext,
			variableData *types.VariableData,
		) (*types.VariableHandlerResult, error) {
			result := &types.VariableHandlerResult{}

			twitchClient, err := twitch.NewUserClientWithContext(
				ctx,
				parseCtx.Channel.TwitchUserID,
				*parseCtx.Services.Config,
				parseCtx.Services.Bus,
			)
			if err != nil {
				parseCtx.Services.Logger.Error(i18n.GetCtx(ctx, locales.Translations.Errors.Generic.CannotCreateTwitch), zap.Error(err))
				return result, nil
			}

			subscribers, err := twitchClient.Subscriptions.GetBroadcasterSubscriptions(
				ctx,
				helix.GetBroadcasterSubscriptionsRequest{
					BroadcasterID: parseCtx.Channel.ID,
				},
			)
			if err != nil {
				parseCtx.Services.Logger.Error(i18n.GetCtx(ctx, locales.Translations.Variables.Subscribers.Errors.GetSubscribers), zap.Error(err))
				return result, nil
			}
			result.Result = fmt.Sprint(lo.FromPtr(subscribers.Data.Total))
			return result, nil
		},
		shared.PlatformKick: func(
			ctx context.Context,
			parseCtx *types.VariableParseContext,
			variableData *types.VariableData,
		) (*types.VariableHandlerResult, error) {
			result := &types.VariableHandlerResult{}

			channelInfo, err := shared.GetKickChannel(ctx, parseCtx)
			if err != nil {
				parseCtx.Services.Logger.Error(i18n.GetCtx(ctx, locales.Translations.Variables.Subscribers.Errors.GetSubscribers), zap.Error(err))
				return result, nil
			}

			if channelInfo == nil {
				return result, nil
			}

			result.Result = fmt.Sprint(channelInfo.ActiveSubscribersCount)
			return result, nil
		},
	}),
}
