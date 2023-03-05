package channel_title

import (
	"github.com/guregu/null"
	"github.com/samber/do"
	"github.com/samber/lo"
	"github.com/satont/go-helix/v2"
	"github.com/satont/tsuwari/apps/parser/internal/di"
	"github.com/satont/tsuwari/apps/parser/internal/types"
	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"
	config "github.com/satont/tsuwari/libs/config"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/grpc/generated/tokens"
	"github.com/satont/tsuwari/libs/twitch"
)

var SetCommand = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "title",
		Description: null.StringFrom("Print or change title of channel."),
		Module:      "MODERATION",
		IsReply:     true,
		Visible:     true,
	},
	Handler: func(ctx *variables_cache.ExecutionContext) *types.CommandsHandlerResult {
		cfg := do.MustInvoke[config.Config](di.Provider)
		tokensGrpc := do.MustInvoke[tokens.TokensClient](di.Provider)

		result := &types.CommandsHandlerResult{
			Result: make([]string, 0),
		}

		_, isHavePermToChange := lo.Find(ctx.SenderBadges, func(item string) bool {
			return item == "BROADCASTER" || item == "MODERATOR"
		})

		twitchClient, err := twitch.NewUserClient(ctx.ChannelId, cfg, tokensGrpc)
		if err != nil {
			return nil
		}

		if ctx.Text == nil || *ctx.Text == "" || !isHavePermToChange {
			channelInfo, err := twitchClient.GetChannelInformation(&helix.GetChannelInformationParams{
				BroadcasterIDs: []string{ctx.ChannelId},
			})

			if err != nil || channelInfo.ErrorMessage != "" || len(channelInfo.Data.Channels) == 0 {
				return nil
			}

			result.Result = append(result.Result, channelInfo.Data.Channels[0].Title)
			return result
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
