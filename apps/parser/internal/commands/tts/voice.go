package tts

import (
	"context"
	"fmt"
	"strings"

	"github.com/guregu/null"
	command_arguments "github.com/twirapp/twir/apps/parser/internal/command-arguments"
	model "github.com/twirapp/twir/libs/gomodels"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
)

const (
	ttsVoiceArgName = "voice"
)

var VoiceCommand = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "tts voice",
		Description: null.StringFrom("Change tts voice"),
		Module:      "TTS",
		IsReply:     true,
	},
	SkipToxicityCheck: true,
	Args: []command_arguments.Arg{
		command_arguments.String{
			Name:     ttsVoiceArgName,
			Optional: true,
		},
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		result := &types.CommandsHandlerResult{}
		channelSettings, channelModel, err := getSettings(
			ctx,
			parseCtx.Services.TTSRepository,
			parseCtx.Channel.ID,
			"",
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: "error while getting channel settings",
				Err:     err,
			}
		}

		if channelSettings == nil {
			result.Result = []string{"TTS is not configured for this channel"}
			return result, nil
		}

		userSettings, currentUserModel, err := getSettings(
			ctx,
			parseCtx.Services.TTSRepository,
			parseCtx.Channel.ID,
			parseCtx.Sender.ID,
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: "error while getting user settings",
				Err:     err,
			}
		}

		textArg := parseCtx.ArgsParser.Get(ttsVoiceArgName)

		if textArg == nil {
			result.Result = append(
				result.Result,
				fmt.Sprintf(
					"Global voice: %s | Your voice: %s",
					channelSettings.Voice,
					lo.IfF(
						userSettings != nil, func() string {
							return userSettings.Voice
						},
					).Else("not setted"),
				),
			)
			return result, nil
		}

		voices := getVoices(ctx, parseCtx.Services.Config)
		if len(voices) == 0 {
			result.Result = append(result.Result, "No voices found")
			return result, nil
		}

		text := textArg.String()

		wantedVoice, ok := lo.Find(
			voices, func(item Voice) bool {
				return item.Name == strings.ToLower(text)
			},
		)

		if !ok {
			result.Result = append(result.Result, fmt.Sprintf("Voice %s not found", text))
			return result, nil
		}

		_, isDisallowed := lo.Find(
			channelSettings.DisallowedVoices, func(item string) bool {
				return item == wantedVoice.Name
			},
		)

		if isDisallowed {
			result.Result = append(
				result.Result,
				fmt.Sprintf("Voice %s is disallowed for usage", wantedVoice.Name),
			)
			return result, nil
		}

		if parseCtx.Channel.ID == parseCtx.Sender.ID {
			channelSettings.Voice = wantedVoice.Name
			err := updateSettings(ctx, parseCtx.Services.TTSRepository, channelModel, channelSettings)
			if err != nil {
				return nil, &types.CommandHandlerError{
					Message: "error while updating settings",
					Err:     err,
				}
			}
		} else {
			if userSettings == nil {
				_, _, err := createUserSettings(
					ctx,
					parseCtx.Services.TTSRepository,
					50,
					50,
					wantedVoice.Name,
					parseCtx.Channel.ID,
					parseCtx.Sender.ID,
				)
				if err != nil {
					return nil, &types.CommandHandlerError{
						Message: "error while creating user settings",
						Err:     err,
					}
				}
			} else {
				userSettings.Voice = wantedVoice.Name
				err := updateSettings(ctx, parseCtx.Services.TTSRepository, currentUserModel, userSettings)
				if err != nil {
					return nil, &types.CommandHandlerError{
						Message: "error while updating user settings",
						Err:     err,
					}
				}
			}
		}

		result.Result = append(result.Result, fmt.Sprintf("Voice changed to %s", wantedVoice.Name))

		parseCtx.Services.TTSCache.Invalidate(ctx, parseCtx.Channel.ID)

		return result, nil
	},
}
