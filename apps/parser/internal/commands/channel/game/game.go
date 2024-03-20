package channel_game

import (
	"context"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/parser/internal/types"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/twitch"

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
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		result := &types.CommandsHandlerResult{
			Result: make([]string, 0),
		}

		if parseCtx.Text == nil || *parseCtx.Text == "" {
			return result, nil
		}

		category, err := twitch.SearchCategory(ctx, *parseCtx.Text)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: "game not found on twitch",
				Err:     err,
			}
		}

		twitchClient, err := twitch.NewUserClientWithContext(
			ctx,
			parseCtx.Channel.ID,
			*parseCtx.Services.Config,
			parseCtx.Services.GrpcClients.Tokens,
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: "cannot create broadcaster twitch client",
				Err:     err,
			}
		}

		changeResponse, err := twitchClient.EditChannelInformation(
			&helix.EditChannelInformationParams{
				BroadcasterID: parseCtx.Channel.ID,
				GameID:        category.ID,
			},
		)

		if err != nil || changeResponse.StatusCode != 204 {
			result.Result = append(
				result.Result,
				lo.If(changeResponse.ErrorMessage != "", changeResponse.ErrorMessage).Else(
					"❌ internal error",
				),
			)
			return result, nil
		}

		result.Result = append(result.Result, "✅ "+category.Name)
		return result, nil
	},
}
