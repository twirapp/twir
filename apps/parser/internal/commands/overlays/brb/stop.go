package brb

import (
	"context"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/satont/twir/apps/parser/internal/types"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/generated/websockets"
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
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		result := types.CommandsHandlerResult{}

		if _, err := parseCtx.Services.GrpcClients.WebSockets.TriggerHideBrb(
			ctx,
			&websockets.TriggerHideBrbRequest{
				ChannelId: parseCtx.Channel.ID,
			},
		); err != nil {
			return &result, &types.CommandHandlerError{
				Message: "cannot trigger hide brb",
				Err:     err,
			}
		}

		return &result, nil
	},
}
