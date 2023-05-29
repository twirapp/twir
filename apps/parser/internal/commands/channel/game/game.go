package channel_game

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/satont/tsuwari/apps/parser/internal/types"
	"gorm.io/gorm"

	"github.com/guregu/null"
	"github.com/lib/pq"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/twitch"

	"github.com/nicklaw5/helix/v2"
)

var SetCommand = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "game",
		Description: null.StringFrom("Change category of channel"),
		Module:      "MODERATION",
		IsReply:     true,
		Visible:     false,
		RolesIDS:    pq.StringArray{model.ChannelRoleTypeModerator.String()},
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) *types.CommandsHandlerResult {
		result := &types.CommandsHandlerResult{
			Result: make([]string, 0),
		}

		twitchClient, err := twitch.NewUserClientWithContext(
			ctx,
			parseCtx.Channel.ID,
			*parseCtx.Services.Config,
			parseCtx.Services.GrpcClients.Tokens,
		)
		if err != nil {
			return nil
		}

		if parseCtx.Text == nil || *parseCtx.Text == "" {
			return result
		}
		categoryFromReq := *parseCtx.Text
		categoryFromAlias := &model.ChannelCategoryAlias{}

		err = parseCtx.Services.Gorm.Table("channels_categories_aliases").
			Where(`"channelId" = ? AND "alias" = ?`, parseCtx.Channel.ID, categoryFromReq).Find(categoryFromAlias).Error
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				parseCtx.Services.Logger.Error(err.Error())
				return nil
			}
		}
		if categoryFromAlias.Category != "" {
			req, err := twitchClient.EditChannelInformation(&helix.EditChannelInformationParams{
				BroadcasterID: parseCtx.Channel.ID,
				GameID:        categoryFromAlias.CategoryId,
			})

			if err != nil || req.StatusCode != 204 {
				fmt.Println(err)
				fmt.Println(req.ErrorMessage)
				result.Result = append(result.Result, "❌ internal error")
				return result
			}

			result.Result = append(result.Result, "✅ "+categoryFromAlias.Category)
			return result
		}

		gameReq, err := twitchClient.GetGames(&helix.GamesParams{
			Names: []string{categoryFromReq},
		})
		if err != nil {
			return nil
		}

		categoryId := ""
		categoryName := ""

		if len(gameReq.Data.Games) > 0 {
			categoryId = gameReq.Data.Games[0].ID
			categoryName = gameReq.Data.Games[0].Name
		} else {
			games, err := twitchClient.SearchCategories(&helix.SearchCategoriesParams{
				Query: categoryFromReq,
			})
			if err != nil {
				return nil
			}

			if len(games.Data.Categories) > 0 {
				categoryId = games.Data.Categories[0].ID
				categoryName = games.Data.Categories[0].Name

				for _, category := range games.Data.Categories {
					if strings.Index(strings.ToLower(category.Name), strings.ToLower(*parseCtx.Text)) == 0 {
						categoryId = category.ID
						categoryName = category.Name
						break
					}
				}
			}
		}

		if categoryId == "" || categoryName == "" {
			result.Result = append(result.Result, "❌ game not found on twitch")
			return result
		}

		req, err := twitchClient.EditChannelInformation(&helix.EditChannelInformationParams{
			BroadcasterID: parseCtx.Channel.ID,
			GameID:        categoryId,
		})

		if err != nil || req.StatusCode != 204 {
			fmt.Println(err)
			fmt.Println(req.ErrorMessage)
			result.Result = append(result.Result, "❌ internal error")
			return result
		}

		result.Result = append(result.Result, "✅ "+categoryName)
		return result
	},
}
