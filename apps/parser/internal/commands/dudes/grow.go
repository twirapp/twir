package dudes

import (
	"context"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/satont/twir/apps/parser/internal/types"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/bus-core/websockets"
)

var Grow = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "dudes grow",
		Description: null.StringFrom("Increase the size of user in the dudes overlay"),
		Module:      "GAMES",
		IsReply:     true,
		Visible:     true,
		RolesIDS:    pq.StringArray{},
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		result := types.CommandsHandlerResult{}

		err := parseCtx.Services.Bus.Websocket.DudesGrow.Publish(
			websockets.DudesGrowRequest{
				ChannelID: parseCtx.Channel.ID,
				UserID:    parseCtx.Sender.ID,
			},
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: "cannot trigger dudes grow",
				Err:     err,
			}
		}

		return &result, nil
	},
}
