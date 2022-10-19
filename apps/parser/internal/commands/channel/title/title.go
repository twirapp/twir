package channel_title

import (
	"tsuwari/parser/internal/types"
	variables_cache "tsuwari/parser/internal/variablescache"

	"github.com/samber/lo"
	"github.com/satont/go-helix/v2"
)

var Command = types.DefaultCommand{
	Command: types.Command{
		Name:        "title set",
		Description: lo.ToPtr("Changing title of the channel."),
		Permission:  "MODERATOR",
		Visible:     false,
		Module:      lo.ToPtr("CHANNEL"),
		IsReply:     true,
	},
	Handler: func(ctx variables_cache.ExecutionContext) *types.CommandsHandlerResult {
		if ctx.Text == nil {
			return nil
		}

		twitchClient, err := ctx.Services.UsersAuth.Create(ctx.ChannelId)

		if err != nil || twitchClient == nil {
			return nil
		}

		result := &types.CommandsHandlerResult{
			Result: make([]string, 0),
		}

		req, err := twitchClient.EditChannelInformation(&helix.EditChannelInformationParams{
			BroadcasterID: ctx.ChannelId,
			Title:         *ctx.Text,
		})

		if err != nil || req.StatusCode != 204 {
			result.Result = append(result.Result, "❌")
			return result
		}

		result.Result = append(result.Result, "✅ "+*ctx.Text)
		return result
	},
}
