package followers

import (
	"context"
	"fmt"

	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/twitch"
	"go.uber.org/zap"
)

var LatestFollowerUsername = &types.Variable{
	Name:                "followers.latest.userName",
	Description:         lo.ToPtr("Latest follower username"),
	CanBeUsedInRegistry: true,
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}

		entity := model.ChannelsEventsListItem{}
		if err := parseCtx.Services.Gorm.
			WithContext(ctx).
			Where(
				"channel_id = ? AND type = ?",
				parseCtx.Channel.ID,
				model.ChannelEventListItemTypeFollow,
			).
			Order(`"created_at" DESC`).First(&entity).Error; err != nil {
			return result, nil
		}

		result.Result = entity.Data.FollowUserName

		return result, nil
	},
}

var Count = &types.Variable{
	Name:                "followers.count",
	Description:         lo.ToPtr("Followers count"),
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

		followers, err := twitchClient.GetChannelFollows(
			&helix.GetChannelFollowsParams{
				BroadcasterID: parseCtx.Channel.ID,
			},
		)
		if err != nil {
			parseCtx.Services.Logger.Error("cannot get followers", zap.Error(err))
			return result, nil
		}
		if followers.ErrorMessage != "" {
			parseCtx.Services.Logger.Error("cannot get followers", zap.Error(err))
			result.Result = followers.ErrorMessage
			return result, nil
		}

		result.Result = fmt.Sprint(followers.Data.Total)
		return result, nil
	},
}
