package tts

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/guregu/null"
	"github.com/satont/tsuwari/apps/bots/pkg/tlds"
	model "github.com/satont/tsuwari/libs/gomodels"

	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/parser/internal/types"
	"github.com/satont/tsuwari/libs/grpc/generated/websockets"
	"go.uber.org/zap"
)

var emojiRx = regexp.MustCompile(`[\p{So}\p{Sk}\p{Sm}\p{Sc}]`)

// [a-zA-Z0-9]+([a-zA-Z0-9-]+)?\\.(${tlds.join('|')})
var linksWithSpaces = regexp.MustCompile(
	fmt.Sprintf(
		`(https?:\/\/)?(www)? ??\.? ?[a-zA-Z0-9]+([a-zA-Z0-9-]+) ??\. ?(%s)\b`,
		strings.Join(tlds.TLDS, "|"),
	),
)

var SayCommand = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "tts",
		Description: null.StringFrom("Say text in tts. You can use !tts <voice> <text> to send tts with some voice."),
		Visible:     true,
		Module:      "TTS",
		IsReply:     true,
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) *types.CommandsHandlerResult {
		result := &types.CommandsHandlerResult{}

		if parseCtx.Text == nil {
			return result
		}

		channelSettings, _ := getSettings(ctx, parseCtx.Services.Gorm, parseCtx.Channel.ID, "")
		if channelSettings == nil || !*channelSettings.Enabled {
			return result
		}

		userSettings, _ := getSettings(
			ctx,
			parseCtx.Services.Gorm,
			parseCtx.Channel.ID,
			parseCtx.Sender.ID,
		)

		voice := lo.IfF(userSettings != nil, func() string {
			return userSettings.Voice
		}).
			Else(channelSettings.Voice)

		if channelSettings.AllowUsersChooseVoiceInMainCommand {
			voices := getVoices(ctx, parseCtx.Services.Config)
			splittedChatArgs := strings.Split(*parseCtx.Text, " ")
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

				*parseCtx.Text = strings.Join(splittedChatArgs[1:], " ")
			}
		}

		if channelSettings.MaxSymbols > 0 && utf8.RuneCountInString(*parseCtx.Text) > channelSettings.MaxSymbols {
			return result
		}

		rate := lo.IfF(userSettings != nil, func() int {
			return userSettings.Rate
		}).Else(channelSettings.Rate)
		pitch := lo.IfF(userSettings != nil, func() int {
			return userSettings.Pitch
		}).Else(channelSettings.Pitch)

		if channelSettings.DoNotReadEmoji {
			*parseCtx.Text = emojiRx.ReplaceAllString(*parseCtx.Text, ``)
		}

		if channelSettings.DoNotReadLinks {
			*parseCtx.Text = strings.TrimSpace(linksWithSpaces.ReplaceAllString(*parseCtx.Text, ``))
		}

		if channelSettings.DoNotReadTwitchEmotes {
			for _, emote := range parseCtx.Emotes {
				*parseCtx.Text = strings.Replace(*parseCtx.Text, emote.Name, "", -1)
			}
			channelKey := fmt.Sprintf("emotes:channel:%s:", parseCtx.Channel.ID)
			channelEmotes := parseCtx.Services.Redis.Keys(
				context.Background(),
				fmt.Sprintf("%s*", channelKey),
			).Val()

			for _, emote := range channelEmotes {
				*parseCtx.Text = strings.Replace(*parseCtx.Text, strings.Split(emote, channelKey)[1], "", -1)
			}

			globalKey := "emotes:global:"
			globalEmotes := parseCtx.Services.Redis.Keys(
				context.Background(),
				fmt.Sprintf("%s:*", globalKey),
			).Val()

			for _, emote := range globalEmotes {
				*parseCtx.Text = strings.Replace(*parseCtx.Text, strings.Split(emote, globalKey)[1], "", -1)
			}

		}

		_, err := parseCtx.Services.GrpcClients.WebSockets.TextToSpeechSay(context.Background(), &websockets.TTSMessage{
			ChannelId: parseCtx.Channel.ID,
			Text:      *parseCtx.Text,
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
