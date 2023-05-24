package top

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/parser/internal/types"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/twitch"
)

var EmotesUsers = &types.Variable{
	Name:        "top.emotes.users",
	Description: lo.ToPtr("Top users by emotes. Prints counts of used emotes"),
	Example:     lo.ToPtr("top.emotes.users|10"),
	Handler: func(ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}

		twitchClient, err := twitch.NewAppClientWithContext(
			ctx,
			*parseCtx.Services.Config,
			parseCtx.Services.GrpcClients.Tokens,
		)
		if err != nil {
			parseCtx.Services.Logger.Sugar().Error(err)
			return nil, err
		}

		limit := 10
		if variableData.Params != nil {
			newLimit, err := strconv.Atoi(*variableData.Params)
			if err == nil {
				limit = newLimit
			}
		}

		if limit > 50 {
			limit = 10
		}

		var usages []*model.ChannelEmoteUsageWithCount
		err = parseCtx.Services.Gorm.
			Model(&model.ChannelEmoteUsageWithCount{}).
			Select(`"userId", COUNT(*) as count`).
			Where(`"channelId" = ?`, parseCtx.Channel.ID).
			Group(`"userId"`).
			Order("count DESC").
			Limit(10).
			Scan(&usages).
			WithContext(ctx).
			Error

		if err != nil {
			return nil, err
		}

		twitchUsers, err := twitchClient.GetUsers(&helix.UsersParams{
			IDs: lo.Map(usages, func(item *model.ChannelEmoteUsageWithCount, _ int) string {
				return item.UserID
			}),
		})
		if err != nil {
			parseCtx.Services.Logger.Sugar().Error(err)
			return nil, err
		}

		mappedTop := []string{}

		for _, usage := range usages {
			user, ok := lo.Find(twitchUsers.Data.Users, func(item helix.User) bool {
				return item.ID == usage.UserID
			})

			if !ok {
				continue
			}

			mappedTop = append(mappedTop, fmt.Sprintf(
				"%s × %v",
				user.Login,
				usage.Count,
			))
		}

		result.Result = strings.Join(mappedTop, " · ")
		return result, nil
	},
}
