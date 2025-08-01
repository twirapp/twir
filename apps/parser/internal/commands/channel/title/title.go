package channel_title

import (
	"context"
	"fmt"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	command_arguments "github.com/twirapp/twir/apps/parser/internal/command-arguments"
	"github.com/twirapp/twir/apps/parser/internal/types"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/twitch"
)

const (
	titleArgName = "title"
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
	Args: []command_arguments.Arg{
		command_arguments.VariadicString{
			Name:     titleArgName,
			Optional: true,
		},
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		result := &types.CommandsHandlerResult{
			Result: make([]string, 0),
		}

		twitchClient, err := twitch.NewUserClientWithContext(
			ctx,
			parseCtx.Channel.ID,
			*parseCtx.Services.Config,
			parseCtx.Services.Bus,
		)
		if err != nil {
			return nil, fmt.Errorf("cannot create broadcaster twitch api client: %w", err)
		}

		if !parseCtx.ArgsParser.IsExists(titleArgName) {
			channelInfo, err := twitchClient.GetChannelInformation(
				&helix.GetChannelInformationParams{
					BroadcasterIDs: []string{parseCtx.Channel.ID},
				},
			)
			if err != nil {
				return nil, &types.CommandHandlerError{
					Message: "cannot get channel information",
					Err:     err,
				}
			}
			if len(channelInfo.Data.Channels) == 0 {
				return nil, &types.CommandHandlerError{
					Message: "channel not found",
					Err:     fmt.Errorf("channel not found"),
				}
			}

			result.Result = append(result.Result, channelInfo.Data.Channels[0].Title)
			return result, nil
		}

		title := parseCtx.ArgsParser.Get(titleArgName).String()

		req, err := twitchClient.EditChannelInformation(
			&helix.EditChannelInformationParams{
				BroadcasterID: parseCtx.Channel.ID,
				Title:         title,
			},
		)

		if err != nil || req.StatusCode != 204 {
			result.Result = append(
				result.Result,
				lo.If(req.ErrorMessage != "", req.ErrorMessage).Else("❌ internal error"),
			)
			return result, nil
		}

		result.Result = append(result.Result, "✅ "+title)
		return result, nil
	},
}
