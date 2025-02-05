package tts

import (
	"context"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/guregu/null"
	command_arguments "github.com/satont/twir/apps/parser/internal/command-arguments"
	model "github.com/satont/twir/libs/gomodels"

	"github.com/samber/lo"
	"github.com/satont/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/libs/grpc/websockets"
)

var emojiRx = regexp.MustCompile(`[\p{So}\p{Sk}\p{Sm}\p{Sc}]`)

const (
	ttsSayArgName = "text"
)

var SayCommand = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "tts",
		Description: null.StringFrom("Say text in tts. You can use !tts <voice> <text> to send tts with some voice."),
		Visible:     true,
		Module:      "TTS",
		IsReply:     true,
	},
	SkipToxicityCheck: true,
	Args: []command_arguments.Arg{
		command_arguments.VariadicString{
			Name: ttsSayArgName,
		},
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		result := &types.CommandsHandlerResult{}

		resultedText := parseCtx.ArgsParser.Get(ttsSayArgName).String()

		channelSettings, _ := getSettings(ctx, parseCtx.Services.Gorm, parseCtx.Channel.ID, "")
		if channelSettings == nil || !*channelSettings.Enabled {
			return result, nil
		}

		userSettings, _ := getSettings(
			ctx,
			parseCtx.Services.Gorm,
			parseCtx.Channel.ID,
			parseCtx.Sender.ID,
		)

		voice := lo.IfF(
			userSettings != nil, func() string {
				return userSettings.Voice
			},
		).
			Else(channelSettings.Voice)

		if channelSettings.AllowUsersChooseVoiceInMainCommand {
			voices := getVoices(ctx, parseCtx.Services.Config)
			splittedChatArgs := strings.Split(resultedText, " ")
			targetVoice, targetVoiceFound := lo.Find(
				voices, func(item Voice) bool {
					return strings.ToLower(item.Name) == strings.ToLower(splittedChatArgs[0])
				},
			)

			if targetVoiceFound {
				voice = targetVoice.Name

				_, isDisallowed := lo.Find(
					channelSettings.DisallowedVoices, func(item string) bool {
						return item == voice
					},
				)

				if isDisallowed {
					result.Result = append(
						result.Result,
						fmt.Sprintf("Voice %s is disallowed fopr usage", voice),
					)
					return result, nil
				}

				resultedText = strings.Join(splittedChatArgs[1:], " ")
			}
		}

		if channelSettings.MaxSymbols > 0 && utf8.RuneCountInString(resultedText) > channelSettings.MaxSymbols {
			return result, nil
		}

		rate := lo.IfF(
			userSettings != nil, func() int {
				return userSettings.Rate
			},
		).Else(channelSettings.Rate)
		pitch := lo.IfF(
			userSettings != nil, func() int {
				return userSettings.Pitch
			},
		).Else(channelSettings.Pitch)

		if channelSettings.DoNotReadEmoji {
			resultedText = emojiRx.ReplaceAllString(resultedText, ``)
		}

		if channelSettings.DoNotReadLinks {
			for _, part := range strings.Fields(resultedText) {
				isUrl := isValidUrl(part)
				if isUrl {
					resultedText = strings.ReplaceAll(resultedText, part, "")
				}
			}
		}

		if channelSettings.DoNotReadTwitchEmotes {
			for _, emote := range parseCtx.Emotes {
				resultedText = strings.Replace(resultedText, emote.Name, "", -1)
			}

			channelEmotesRedisKey := fmt.Sprintf("emotes:channel:%s:", parseCtx.Channel.ID)
			var channelEmotes []string
			channelEmotesIter := parseCtx.Services.Redis.Scan(ctx, 0, channelEmotesRedisKey, 0).Iterator()
			for channelEmotesIter.Next(ctx) {
				channelEmotes = append(channelEmotes, channelEmotesIter.Val())
			}
			if err := channelEmotesIter.Err(); err != nil {
				parseCtx.Services.Logger.Sugar().Error("error while getting channel emotes", err)
			}

			globalEmotesRedisKey := "emotes:global:"
			var globalEmotes []string
			iter := parseCtx.Services.Redis.Scan(ctx, 0, globalEmotesRedisKey, 0).Iterator()
			for iter.Next(ctx) {
				globalEmotes = append(globalEmotes, iter.Val())
			}
			if err := iter.Err(); err != nil {
				parseCtx.Services.Logger.Sugar().Error("error while getting global emotes", err)
			}

			for _, emotePattern := range channelEmotes {
				emote := strings.Split(emotePattern, channelEmotesRedisKey)[1]

				if slices.Contains(strings.Fields(resultedText), emote) {
					resultedText = strings.ReplaceAll(resultedText, emote, "")
				}
			}

			for _, emotePattern := range globalEmotes {
				emote := strings.Split(emotePattern, globalEmotesRedisKey)[1]

				if slices.Contains(strings.Fields(resultedText), emote) {
					resultedText = strings.ReplaceAll(resultedText, emote, "")
				}
			}
		}

		if len(resultedText) == 0 || resultedText == parseCtx.Sender.Name {
			return result, nil
		}

		_, err := parseCtx.Services.GrpcClients.WebSockets.TextToSpeechSay(
			ctx, &websockets.TTSMessage{
				ChannelId: parseCtx.Channel.ID,
				Text:      resultedText,
				Voice:     voice,
				Rate:      strconv.Itoa(rate),
				Pitch:     strconv.Itoa(pitch),
				Volume:    strconv.Itoa(channelSettings.Volume),
			},
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: "error while sending message to tts service",
				Err:     err,
			}
		}

		return result, nil
	},
}
