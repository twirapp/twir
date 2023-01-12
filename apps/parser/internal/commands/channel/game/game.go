package channel_game

import (
	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/parser/internal/di"
	"github.com/satont/tsuwari/apps/parser/internal/types"
	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"
	config "github.com/satont/tsuwari/libs/config"
	"github.com/satont/tsuwari/libs/grpc/generated/tokens"
	"github.com/satont/tsuwari/libs/twitch"

	"github.com/samber/lo"
	"github.com/satont/go-helix/v2"
)

var Command = types.DefaultCommand{
	Command: types.Command{
		Name:        "game set",
		Description: lo.ToPtr("Changing game of the channel."),
		Permission:  "MODERATOR",
		Visible:     false,
		Module:      lo.ToPtr("MODERATION"),
		IsReply:     true,
	},
	Handler: func(ctx variables_cache.ExecutionContext) *types.CommandsHandlerResult {
		result := &types.CommandsHandlerResult{
			Result: make([]string, 0),
		}

		if ctx.Text == nil {
			return nil
		}

		cfg := do.MustInvoke[config.Config](di.Provider)
		tokensGrpc := do.MustInvoke[tokens.TokensClient](di.Provider)

		twitchClient, err := twitch.NewUserClient(ctx.ChannelId, cfg, tokensGrpc)
		if err != nil {
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
