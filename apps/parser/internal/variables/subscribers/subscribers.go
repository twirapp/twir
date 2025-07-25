package subscribers

import (
	"context"
	"fmt"

	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
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

		twitchClient, err := twitch.NewUserClientWithContext(
			ctx,
			parseCtx.Channel.ID,
			*parseCtx.Services.Config,
			parseCtx.Services.Bus,
		)
		if err != nil {
			parseCtx.Services.Logger.Error("cannot create twitch client", zap.Error(err))
			return result, nil
		}

		subscribers, err := twitchClient.GetSubscriptions(
			&helix.SubscriptionsParams{
				BroadcasterID: parseCtx.Channel.ID,
				First:         1,
			},
		)
		if err != nil {
			parseCtx.Services.Logger.Error("cannot get subscribers", zap.Error(err))
			return result, nil
		}
		if subscribers.ErrorMessage != "" {
			parseCtx.Services.Logger.Error("cannot get subscribers", zap.Error(err))
			result.Result = subscribers.ErrorMessage
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
	Handler: func(
		ctx context.Context,
		parseCtx *types.VariableParseContext,
		variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}

		twitchClient, err := twitch.NewUserClientWithContext(
			ctx,
			parseCtx.Channel.ID,
			*parseCtx.Services.Config,
			parseCtx.Services.Bus,
		)
		if err != nil {
			parseCtx.Services.Logger.Error("cannot create twitch client", zap.Error(err))
			return result, nil
		}

		subscribers, err := twitchClient.GetSubscriptions(
			&helix.SubscriptionsParams{
				BroadcasterID: parseCtx.Channel.ID,
			},
		)
		if err != nil {
			parseCtx.Services.Logger.Error("cannot get subscribers", zap.Error(err))
			return result, nil
		}
		if subscribers.ErrorMessage != "" {
			parseCtx.Services.Logger.Error("cannot get subscribers", zap.Error(err))
			result.Result = subscribers.ErrorMessage
			return result, nil
		}

		result.Result = fmt.Sprint(subscribers.Data.Total)
		return result, nil
	},
}
