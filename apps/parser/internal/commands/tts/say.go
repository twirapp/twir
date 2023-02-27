package tts

import (
	"context"
	"fmt"
	"github.com/samber/do"
	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/parser/internal/di"
	"github.com/satont/tsuwari/apps/parser/internal/types"
	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"
	"github.com/satont/tsuwari/libs/grpc/generated/websockets"
	"strconv"
)

var SayVariable = types.DefaultCommand{
	Command: types.Command{
		Name:        "tts",
		Description: lo.ToPtr("Say text in tts"),
		Permission:  "VIEWER",
		Visible:     false,
		Module:      lo.ToPtr("TTS"),
		IsReply:     true,
	},
	Handler: func(ctx variables_cache.ExecutionContext) *types.CommandsHandlerResult {
		webSocketsGrpc := do.MustInvoke[websockets.WebsocketClient](di.Provider)

		result := &types.CommandsHandlerResult{}
		if ctx.Text == nil {
			return result
		}

		channelSettings := getSettings(ctx.ChannelId, "")

		if channelSettings == nil || !*channelSettings.Enabled {
			return result
		}

		userSettings := getSettings(ctx.ChannelId, ctx.SenderId)

		voice := lo.IfF(userSettings != nil, func() string {
			return userSettings.Voice
		}).Else(channelSettings.Voice)
		rate := lo.IfF(userSettings != nil, func() int {
			return userSettings.Rate
		}).Else(channelSettings.Rate)
		pitch := lo.IfF(userSettings != nil, func() int {
			return userSettings.Pitch
		}).Else(channelSettings.Pitch)

		_, err := webSocketsGrpc.TextToSpeechSay(context.Background(), &websockets.TTSMessage{
			ChannelId: ctx.ChannelId,
			Text:      *ctx.Text,
			Voice:     voice,
			Rate:      strconv.Itoa(rate),
			Pitch:     strconv.Itoa(pitch),
			Volume:    strconv.Itoa(channelSettings.Volume),
		})

		if err != nil {
			fmt.Println(err)
			return result
		}

		return result
	},
}
