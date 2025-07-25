package kappagen

import (
	"context"
	"fmt"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/twirapp/twir/apps/parser/internal/types"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/bus-core/api"
)

var Kappagen = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "kappagen",
		Description: null.StringFrom("Send smiles to kappagen overlay."),
		Module:      "OVERLAYS",
		IsReply:     true,
		Visible:     true,
		RolesIDS:    pq.StringArray{},
	},
	SkipToxicityCheck: true,
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		var emotes []api.TriggerKappagenEmote
		for _, e := range parseCtx.Emotes {
			emote := api.TriggerKappagenEmote{
				Id:        e.ID,
				Positions: []string{},
			}

			for _, pos := range e.Positions {
				emote.Positions = append(emote.Positions, fmt.Sprintf("%v-%v", pos.Start, pos.End))
			}

			emotes = append(emotes, emote)
		}

		param := "!" + parseCtx.RawText

		err := parseCtx.Services.Bus.Api.TriggerKappagen.Publish(
			ctx,
			api.TriggerKappagenMessage{
				ChannelId: parseCtx.Channel.ID,
				Text:      param,
				Emotes:    emotes,
			},
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: "error while triggering kappagen",
				Err:     err,
			}
		}

		return &types.CommandsHandlerResult{}, nil
	},
}
