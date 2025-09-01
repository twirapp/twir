package tts

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/guregu/null"
	command_arguments "github.com/twirapp/twir/apps/parser/internal/command-arguments"
	"github.com/twirapp/twir/apps/parser/internal/services/tts"
	emotes_cacher "github.com/twirapp/twir/libs/bus-core/emotes-cacher"
	model "github.com/twirapp/twir/libs/gomodels"
	"golang.org/x/sync/errgroup"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
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
		splittedResult := strings.Fields(resultedText)

		channelSettings, _, err := parseCtx.Services.TTSService.GetChannelSettings(
			ctx,
			parseCtx.Channel.ID,
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: "error while getting channel settings",
				Err:     err,
			}
		}

		if channelSettings == nil || !*channelSettings.Enabled {
			return result, nil
		}

		userSettings, _, err := parseCtx.Services.TTSService.GetUserSettings(
			ctx,
			parseCtx.Channel.ID,
			parseCtx.Sender.ID,
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: "error while getting user settings",
				Err:     err,
			}
		}

		voice := lo.IfF(
			userSettings != nil, func() string {
				return userSettings.Voice
			},
		).
			Else(channelSettings.Voice)

		if channelSettings.AllowUsersChooseVoiceInMainCommand {
			voices := parseCtx.Services.TTSService.GetAvailableVoices(ctx)
			splittedChatArgs := strings.Split(resultedText, " ")
			targetVoice, targetVoiceFound := lo.Find(
				voices, func(item tts.Voice) bool {
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
						fmt.Sprintf("Voice %s is disallowed for usage", voice),
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
				isUrl := parseCtx.Services.TTSService.IsValidURL(part)
				if isUrl {
					resultedText = strings.ReplaceAll(resultedText, part, "")
				}
			}
		}

		if channelSettings.DoNotReadTwitchEmotes {
			for _, emote := range parseCtx.Emotes {
				resultedText = strings.Replace(resultedText, emote.Name, "", -1)
			}

			var (
				wg            errgroup.Group
				globalEmotes  []emotes_cacher.Emote
				channelEmotes []emotes_cacher.Emote
			)

			wg.Go(
				func() error {
					e, err := parseCtx.Services.Bus.EmotesCacher.GetGlobalEmotes.Request(
						ctx,
						emotes_cacher.GetGlobalEmotesRequest{},
					)
					if err != nil {
						return err
					}

					globalEmotes = e.Data.Emotes

					return nil
				},
			)

			wg.Go(
				func() error {
					e, err := parseCtx.Services.Bus.EmotesCacher.GetChannelEmotes.Request(
						ctx,
						emotes_cacher.GetChannelEmotesRequest{
							ChannelID: parseCtx.Channel.ID,
						},
					)
					if err != nil {
						return err
					}

					channelEmotes = e.Data.Emotes

					return nil
				},
			)

			for _, part := range splittedResult {
				for _, emote := range globalEmotes {
					if strings.EqualFold(part, emote.Name) {
						resultedText = strings.Replace(resultedText, part, "", -1)
						continue
					}
				}

				for _, emote := range channelEmotes {
					if strings.EqualFold(part, emote.Name) {
						resultedText = strings.Replace(resultedText, part, "", -1)
						continue
					}
				}
			}
		}

		if len(resultedText) == 0 || resultedText == parseCtx.Sender.Name {
			return result, nil
		}

		_, err = parseCtx.Services.GrpcClients.WebSockets.TextToSpeechSay(
			ctx,
			&websockets.TTSMessage{
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
