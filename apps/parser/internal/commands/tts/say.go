package tts

import (
	"context"
	"github.com/samber/do"
	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/parser/internal/di"
	"github.com/satont/tsuwari/apps/parser/internal/types"
	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"
	"github.com/satont/tsuwari/libs/grpc/generated/websockets"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"unicode/utf8"
)

var SayCommand = types.DefaultCommand{
	Command: types.Command{
		Name:        "tts",
		Description: lo.ToPtr("Say text in tts. You can use !tts <voice> <text> to send tts with some voice."),
		Permission:  "VIEWER",
		Visible:     true,
		Module:      lo.ToPtr("TTS"),
		IsReply:     true,
	},
	Handler: func(ctx variables_cache.ExecutionContext) *types.CommandsHandlerResult {
		webSocketsGrpc := do.MustInvoke[websockets.WebsocketClient](di.Provider)

		result := &types.CommandsHandlerResult{}

		if ctx.Text == nil {
			return result
		}

		channelSettings, _ := getSettings(ctx.ChannelId, "")
		if channelSettings == nil || !*channelSettings.Enabled {
			return result
		}

		userSettings, _ := getSettings(ctx.ChannelId, ctx.SenderId)

		voice := lo.IfF(userSettings != nil, func() string {
			return userSettings.Voice
		}).
			Else(channelSettings.Voice)

		if channelSettings.AllowUsersChooseVoiceInMainCommand {
			voices := getVoices()
			splittedChatArgs := strings.Split(*ctx.Text, " ")
			targetVoice, targetVoiceFound := lo.Find(voices, func(item Voice) bool {
				return strings.ToLower(item.Name) == strings.ToLower(splittedChatArgs[0])
			})

			if targetVoiceFound {
				voice = targetVoice.Name
				*ctx.Text = strings.Join(splittedChatArgs[1:], " ")
			}
		}

		if channelSettings.MaxSymbols > 0 && utf8.RuneCountInString(*ctx.Text) > channelSettings.MaxSymbols {
			return result
		}

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
			zap.S().Error(err)
			return result
		}

		return result
	},
}
