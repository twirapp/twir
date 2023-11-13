package kappagen

import (
	"context"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/satont/twir/apps/parser/internal/types"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/generated/websockets"
)

var Kappagen = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "kappagen",
		Description: null.StringFrom("Magic ball will answer to all your questions!"),
		Module:      "GAMES",
		IsReply:     true,
		Visible:     true,
		RolesIDS:    pq.StringArray{},
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) *types.CommandsHandlerResult {
		if parseCtx.Text == nil || *parseCtx.Text == "" {
			return nil
		}

		parseCtx.Services.GrpcClients.WebSockets.TriggerKappagen(
			ctx, &websockets.TriggerKappagenRequest{
				ChannelId: parseCtx.Channel.ID,
				Text:      *parseCtx.Text,
			},
		)

		return nil
	},
}
