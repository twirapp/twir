package channel_title

import (
	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/parser/internal/di"
	users_twitch_auth "github.com/satont/tsuwari/apps/parser/internal/twitch/user"
	"github.com/satont/tsuwari/apps/parser/internal/types"
	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"

	"github.com/samber/lo"
	"github.com/satont/go-helix/v2"
)

var Command = types.DefaultCommand{
	Command: types.Command{
		Name:        "title set",
		Description: lo.ToPtr("Changing title of the channel."),
		Permission:  "MODERATOR",
		Visible:     false,
		Module:      lo.ToPtr("MODERATION"),
		IsReply:     true,
	},
	Handler: func(ctx variables_cache.ExecutionContext) *types.CommandsHandlerResult {
		if ctx.Text == nil {
			return nil
		}
		users := do.MustInvoke[users_twitch_auth.UsersTokensService](di.Provider)
		twitchClient, err := users.Create(ctx.ChannelId)

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
