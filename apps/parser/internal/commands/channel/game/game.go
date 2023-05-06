package channel_game

import (
	"fmt"
	"strings"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/parser/internal/di"
	"github.com/satont/tsuwari/apps/parser/internal/types"
	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"
	config "github.com/satont/tsuwari/libs/config"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/grpc/generated/tokens"
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
	Handler: func(ctx *variables_cache.ExecutionContext) *types.CommandsHandlerResult {
		cfg := do.MustInvoke[config.Config](di.Provider)
		tokensGrpc := do.MustInvoke[tokens.TokensClient](di.Provider)

		result := &types.CommandsHandlerResult{
			Result: make([]string, 0),
		}

		twitchClient, err := twitch.NewUserClient(ctx.ChannelId, cfg, tokensGrpc)
		if err != nil {
			return nil
		}

		if ctx.Text == nil || *ctx.Text == "" {
			return result
		}

		gameReq, err := twitchClient.GetGames(&helix.GamesParams{
			Names: []string{*ctx.Text},
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
				Query: *ctx.Text,
			})
			if err != nil {
				return nil
			}

			if len(games.Data.Categories) > 0 {
				categoryId = games.Data.Categories[0].ID
				categoryName = games.Data.Categories[0].Name

				for _, category := range games.Data.Categories {
					if strings.Index(strings.ToLower(category.Name), strings.ToLower(*ctx.Text)) == 0 {
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
			BroadcasterID: ctx.ChannelId,
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
