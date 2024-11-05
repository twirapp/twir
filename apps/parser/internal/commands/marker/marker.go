package marker

import (
	"context"
	"errors"
	"fmt"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/nicklaw5/helix/v2"
	command_arguments "github.com/satont/twir/apps/parser/internal/command-arguments"
	"github.com/satont/twir/apps/parser/internal/types"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/twitch"
)

var Marker = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "marker",
		Description: null.StringFrom("Create a stream marker"),
		RolesIDS:    pq.StringArray{model.ChannelRoleTypeModerator.String()},
		Module:      "MODERATION",
		Visible:     true,
		IsReply:     true,
		Aliases:     []string{},
		Enabled:     true,
	},
	Args: []command_arguments.Arg{
		command_arguments.String{
			Name:     "markerDescription",
			Optional: true,
		},
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		dbChannel := &model.Channels{}
		if err := parseCtx.Services.Gorm.First(dbChannel, parseCtx.Channel.ID).Error; err != nil {
			return nil, &types.CommandHandlerError{
				Message: "cannot find channel in db, please contact support",
				Err:     err,
			}
		}

		twitchClient, err := twitch.NewBotClientWithContext(
			ctx,
			dbChannel.BotID,
			*parseCtx.Services.Config,
			parseCtx.Services.GrpcClients.Tokens,
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: "cannot create twitch client, please contact support",
				Err:     err,
			}
		}

		var markerDescription string
		markerArg := parseCtx.ArgsParser.Get("markerArg")
		if markerArg == nil || markerArg.IsOptional() || markerArg.String() == "" {
			markerDescription = "Marker created by twirapp"
		} else {
			markerDescription = markerArg.String()
		}

		resp, err := twitchClient.CreateStreamMarker(
			&helix.CreateStreamMarkerParams{
				UserID:      dbChannel.ID,
				Description: "",
			},
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: "cannot create stream marker",
				Err:     err,
			}
		}
		if resp.StatusCode == 403 {
			return nil, &types.CommandHandlerError{
				Message: fmt.Sprintf(
					"cannot create stream marker, please make bot as editor, it can be done at https://dashboard.twitch.tv/u/%s/community/roles page",
					parseCtx.Channel.Name,
				),
				Err: errors.New("insufficient permissions"),
			}
		}

		if resp.ErrorMessage != "" {
			return nil, &types.CommandHandlerError{
				Message: fmt.Sprintf("cannot create stream marker: %s", resp.ErrorMessage),
				Err:     err,
			}
		}

		if len(resp.Data.CreateStreamMarkers) == 0 {
			return nil, &types.CommandHandlerError{
				Message: "cannot create stream marker",
				Err:     errors.New("empty stream marker response"),
			}
		}

		marker := resp.Data.CreateStreamMarkers[0]

		return &types.CommandsHandlerResult{
			Result: []string{
				fmt.Sprintf(
					`Marker created at %v second of stream with description: "%s"`,
					marker.PositionSeconds,
					markerDescription,
				),
			},
		}, nil
	},
}
