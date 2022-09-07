package channel_title

import (
	"tsuwari/parser/internal/types"
	variables_cache "tsuwari/parser/internal/variablescache"

	"github.com/samber/lo"
	"github.com/satont/go-helix"
)

var Command = types.DefaultCommand{
	Command: types.Command{
		Name:        "title set",
		Description: lo.ToPtr("Changing title of the channel."),
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

		req, err := twitchClient.EditChannelInformation(&helix.EditChannelInformationParams{
			BroadcasterID: ctx.ChannelId,
			Title:         *ctx.Text,
		})

		if err != nil || req.StatusCode != 204 {
			return []string{"❌"}
		}

		return []string{"✅ " + *ctx.Text}
	},
}
