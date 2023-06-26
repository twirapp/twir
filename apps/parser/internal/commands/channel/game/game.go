package channel_game

import (
	"context"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/parser/internal/types"
	"strings"

	"github.com/guregu/null"
	"github.com/lib/pq"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/twitch"

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

		gameReq, err := twitchClient.GetGames(
			&helix.GamesParams{
				Names: []string{*parseCtx.Text},
			},
		)
		if err != nil {
			return nil
		}

		categoryId := ""
		categoryName := ""

		if len(gameReq.Data.Games) > 0 {
			categoryId = gameReq.Data.Games[0].ID
			categoryName = gameReq.Data.Games[0].Name
		} else {
			games, err := twitchClient.SearchCategories(
				&helix.SearchCategoriesParams{
					Query: *parseCtx.Text,
				},
			)
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

		req, err := twitchClient.EditChannelInformation(
			&helix.EditChannelInformationParams{
				BroadcasterID: parseCtx.Channel.ID,
				GameID:        categoryId,
			},
		)

		if err != nil || req.StatusCode != 204 {
			result.Result = append(
				result.Result,
				lo.If(req.ErrorMessage != "", req.ErrorMessage).Else("❌ internal error"),
			)
			return result
		}

		result.Result = append(result.Result, "✅ "+categoryName)
		return result
	},
}
