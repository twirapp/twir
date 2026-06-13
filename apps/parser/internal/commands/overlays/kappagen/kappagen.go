package kappagen

import (
	"context"
	"fmt"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/libs/bus-core/api"
	model "github.com/twirapp/twir/libs/gomodels"
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
		emotes := BuildTriggerEmotes(parseCtx.Emotes)

		param := "!" + parseCtx.RawText

		err := parseCtx.Services.Bus.Api.TriggerKappagen.Publish(
			ctx,
			api.TriggerKappagenMessage{
				ChannelId: parseCtx.Channel.DBChannelID,
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

func BuildTriggerEmotes(parseEmotes []*types.ParseContextEmote) []api.TriggerKappagenEmote {
	emotes := make([]api.TriggerKappagenEmote, 0, len(parseEmotes))
	emoteIndexes := make(map[string]int, len(parseEmotes))
	for _, parseEmote := range parseEmotes {
		key := parseEmote.ID + "\x00" + parseEmote.URL
		index, ok := emoteIndexes[key]
		if !ok {
			index = len(emotes)
			emoteIndexes[key] = index
			emotes = append(emotes, api.TriggerKappagenEmote{
				Id:        parseEmote.ID,
				Url:       parseEmote.URL,
				Positions: make([]string, 0, len(parseEmote.Positions)),
			})
		}

		for _, position := range parseEmote.Positions {
			emotes[index].Positions = append(emotes[index].Positions, fmt.Sprintf("%v-%v", position.Start, position.End))
		}
	}

	return emotes
}
