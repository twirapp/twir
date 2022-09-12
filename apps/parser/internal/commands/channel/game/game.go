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
		Visible:     true,
		Module:      lo.ToPtr("CHANNEL"),
	},
	Handler: func(ctx variables_cache.ExecutionContext) []string {
		if ctx.Text == nil {
			return []string{}
		}

		twitchClient, err := ctx.Services.UsersAuth.Create(ctx.ChannelId)

		if err != nil || twitchClient == nil {
			return []string{}
		}

		games, err := twitchClient.SearchCategories(&helix.SearchCategoriesParams{
			Query: *ctx.Text,
		})

		if err != nil || len(games.Data.Categories) == 0 || games.StatusCode != 200 {
			return []string{"game not found on twitch"}
		}

		req, err := twitchClient.EditChannelInformation(&helix.EditChannelInformationParams{
			BroadcasterID: ctx.ChannelId,
			GameID:        games.Data.Categories[0].ID,
		})

		if err != nil || req.StatusCode != 204 {
			return []string{"❌"}
		}

		return []string{"✅ " + games.Data.Categories[0].Name}
	},
}
