package channel_game

import (
	"tsuwari/parser/internal/types"
	variables_cache "tsuwari/parser/internal/variablescache"

	"github.com/samber/lo"
	"github.com/satont/go-helix"
)

var Command = types.DefaultCommand{
	Command: types.Command{
		Name:        "game set",
		Description: lo.ToPtr("Changing game of the channel."),
		Permission:  "MODERATOR",
		Visible:     false,
		Module:      lo.ToPtr("CHANNEL"),
		IsReply:     true,
	},
	Handler: func(ctx variables_cache.ExecutionContext) *types.CommandsHandlerResult {
		result := &types.CommandsHandlerResult{
			Result: make([]string, 0),
		}

		if ctx.Text == nil {
			return nil
		}

		twitchClient, err := ctx.Services.UsersAuth.Create(ctx.ChannelId)

		if err != nil || twitchClient == nil {
			return nil
		}

		games, err := twitchClient.SearchCategories(&helix.SearchCategoriesParams{
			Query: *ctx.Text,
		})

		if err != nil || len(games.Data.Categories) == 0 || games.StatusCode != 200 {
			result.Result = append(result.Result, "game not found on twitch")
			return result
		}

		req, err := twitchClient.EditChannelInformation(&helix.EditChannelInformationParams{
			BroadcasterID: ctx.ChannelId,
			GameID:        games.Data.Categories[0].ID,
		})

		if err != nil || req.StatusCode != 204 {
			result.Result = append(result.Result, "❌")
			return result
		}

		result.Result = append(result.Result, "✅ "+games.Data.Categories[0].Name)
		return result
	},
}
