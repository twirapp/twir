package brb

import (
	"context"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/libs/bus-core/api"
	model "github.com/twirapp/twir/libs/gomodels"
)

var Stop = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "brbstop",
		Description: null.StringFrom("Be right back overlay stop command"),
		Module:      "OVERLAYS",
		IsReply:     true,
		Visible:     false,
		RolesIDS: pq.StringArray{
			model.ChannelRoleTypeBroadcaster.String(),
		},
	},
	SkipToxicityCheck: true,
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		result := types.CommandsHandlerResult{}

		err := parseCtx.Services.Bus.Api.TriggerBrbStop.Publish(
			ctx, api.TriggerBrbStop{
				ChannelId: parseCtx.Channel.ID,
			},
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: "cannot trigger stop brb",
				Err:     err,
			}
		}

		return &result, nil
	},
}
