package tts

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v9"
	"github.com/guregu/null"
	model "github.com/satont/tsuwari/libs/gomodels"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/samber/do"
	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/parser/internal/di"
	"github.com/satont/tsuwari/apps/parser/internal/types"
	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"
	"github.com/satont/tsuwari/libs/grpc/generated/websockets"
	"go.uber.org/zap"
)

var emojiRx = regexp.MustCompile(`[\p{So}\p{Sk}\p{Sm}\p{Sc}]`)

var SayCommand = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "tts",
		Description: null.StringFrom("Say text in tts. You can use !tts <voice> <text> to send tts with some voice."),
		Visible:     true,
		Module:      "TTS",
		IsReply:     true,
	},
	Handler: func(ctx *variables_cache.ExecutionContext) *types.CommandsHandlerResult {
		redisClient := do.MustInvoke[redis.Client](di.Provider)
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

				_, isDisallowed := lo.Find(channelSettings.DisallowedVoices, func(item string) bool {
					return item == voice
				})

				if isDisallowed {
					result.Result = append(
						result.Result,
						fmt.Sprintf("Voice %s is disallowed fopr usage", voice),
					)
					return result
				}

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

		if channelSettings.DoNotReadEmoji {
			*ctx.Text = emojiRx.ReplaceAllString(*ctx.Text, ``)
		}

		if channelSettings.DoNotReadTwitchEmotes {
			for _, emote := range ctx.Emotes {
				*ctx.Text = strings.Replace(*ctx.Text, emote.Name, "", -1)
			}
			channelKey := fmt.Sprintf("emotes:channel:%s:", ctx.ChannelId)
			channelEmotes := redisClient.Keys(
				context.Background(),
				fmt.Sprintf("%s:%s:*", channelKey, ctx.ChannelId),
			).Val()

			for _, emote := range channelEmotes {
				*ctx.Text = strings.Replace(*ctx.Text, strings.Split(emote, channelKey)[1], "", -1)
			}

			globalKey := "emotes:global:"
			globalEmotes := redisClient.Keys(
				context.Background(),
				fmt.Sprintf("%s:*", globalKey),
			).Val()

			for _, emote := range globalEmotes {
				*ctx.Text = strings.Replace(*ctx.Text, strings.Split(emote, globalKey)[1], "", -1)
			}

		}

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
