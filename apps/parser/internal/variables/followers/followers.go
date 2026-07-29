package followers

import (
	"context"
	"fmt"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/i18n"
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
		channelID := parseCtx.Channel.DBChannelID
		if channelID == "" {
			channelID = parseCtx.Channel.ID
		}

		entity := model.ChannelsEventsListItem{}
		if err := parseCtx.Services.Gorm.
			WithContext(ctx).
			Where(
				"channel_id = ? AND type = ?",
				channelID,
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

		if parseCtx.Platform != "twitch" {
			result.Result = "not supported on this platform"
			return result, nil
		}

		followers, err := parseCtx.Services.CacheTwitchClient.GetChannelFollowersCountByChannelId(
			ctx,
			parseCtx.Channel.TwitchUserID,
			parseCtx.Channel.ID,
		)
		if err != nil {
			parseCtx.Services.Logger.Error(i18n.GetCtx(ctx, locales.Translations.Variables.Followers.Errors.GetFollowers), zap.Error(err))
			return result, nil
		}
		result.Result = fmt.Sprint(followers)
		return result, nil
	},
}
