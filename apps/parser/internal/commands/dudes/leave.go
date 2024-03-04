package dudes

import (
	"context"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/satont/twir/apps/parser/internal/types"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/bus-core/websockets"
)

var Leave = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "dudes leave",
		Description: null.StringFrom("Leave dude from overlay"),
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

		err := parseCtx.Services.Bus.Websocket.DudesLeave.Publish(
			websockets.DudesLeaveRequest{
				ChannelID:       parseCtx.Channel.ID,
				UserID:          parseCtx.Sender.ID,
				UserDisplayName: parseCtx.Sender.DisplayName,
				UserName:        parseCtx.Sender.Name,
			},
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: "cannot trigger dudes leave",
				Err:     err,
			}
		}

		return &result, nil
	},
}
