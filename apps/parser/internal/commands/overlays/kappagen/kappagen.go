package kappagen

import (
	"context"
	"fmt"

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

		var emotes []*websockets.TriggerKappagenRequest_Emote
		for _, e := range parseCtx.Emotes {
			emote := &websockets.TriggerKappagenRequest_Emote{
				Id:        e.ID,
				Positions: []string{},
			}

			for _, pos := range e.Positions {
				emote.Positions = append(emote.Positions, fmt.Sprintf("%v-%v", pos.Start, pos.End))
			}

			emotes = append(emotes, emote)
		}

		parseCtx.Services.GrpcClients.WebSockets.TriggerKappagen(
			ctx, &websockets.TriggerKappagenRequest{
				ChannelId: parseCtx.Channel.ID,
				Text:      "!" + parseCtx.RawText,
				Emotes:    emotes,
			},
		)

		return nil
	},
}
