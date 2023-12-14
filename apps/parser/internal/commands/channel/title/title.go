package channel_title

import (
	"context"
	"fmt"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/parser/internal/types"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/twitch"
)

var SetCommand = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "title",
		Description: null.StringFrom("Change category of channel."),
		Module:      "MODERATION",
		IsReply:     true,
		Visible:     false,
		RolesIDS:    pq.StringArray{model.ChannelRoleTypeModerator.String()},
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (*types.CommandsHandlerResult, error) {
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
			return nil, fmt.Errorf("cannot create broadcaster twitch api client: %w", err)
		}

		if parseCtx.Text == nil || *parseCtx.Text == "" {
			return result, nil
		}

		req, err := twitchClient.EditChannelInformation(
			&helix.EditChannelInformationParams{
				BroadcasterID: parseCtx.Channel.ID,
				Title:         *parseCtx.Text,
			},
		)

		if err != nil || req.StatusCode != 204 {
			result.Result = append(
				result.Result,
				lo.If(req.ErrorMessage != "", req.ErrorMessage).Else("❌ internal error"),
			)
			return result, nil
		}

		result.Result = append(result.Result, "✅ "+*parseCtx.Text)
		return result, nil
	},
}
