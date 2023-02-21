package channel_game

import (
	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/parser/internal/di"
	"github.com/satont/tsuwari/apps/parser/internal/types"
	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"
	config "github.com/satont/tsuwari/libs/config"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/grpc/generated/tokens"
	"github.com/satont/tsuwari/libs/twitch"
	"gorm.io/gorm"

	"github.com/samber/lo"
	"github.com/satont/go-helix/v2"
)

var Command = types.DefaultCommand{
	Command: types.Command{
		Name:        "game",
		Description: lo.ToPtr("Print or change category of channel."),
		Permission:  "VIEWER",
		Visible:     false,
		Module:      lo.ToPtr("MODERATION"),
		IsReply:     true,
	},
	Handler: func(ctx variables_cache.ExecutionContext) *types.CommandsHandlerResult {
		db := do.MustInvoke[gorm.DB](di.Provider)
		cfg := do.MustInvoke[config.Config](di.Provider)
		tokensGrpc := do.MustInvoke[tokens.TokensClient](di.Provider)

		result := &types.CommandsHandlerResult{
			Result: make([]string, 0),
		}

		stream := &model.ChannelsStreams{}
		err := db.Where(`"userId" = ?`, ctx.ChannelId).Find(stream).Error
		if err != nil {
			result.Result = append(result.Result, "internal error")
			return result
		}

		if stream.ID == "" {
			result.Result = append(result.Result, "offline")
			return result
		}

		_, isHavePermToChange := lo.Find(ctx.SenderBadges, func(item string) bool {
			return item == "BROADCASTER" || item == "MODERATOR"
		})

		if !isHavePermToChange {
			result.Result = append(result.Result, stream.GameName)
			return result
		}

		if ctx.Text == nil || *ctx.Text == "" {
			result.Result = append(result.Result, stream.GameName)
			return result
		}

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
