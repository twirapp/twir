package tts

import (
	"context"
	"fmt"

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

		if channelSettings == nil {
			result.Result = []string{"TTS is not configured for this channel"}
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

		text := textArg.String()

		wantedVoice, err := parseCtx.Services.TTSService.ValidateVoice(ctx, parseCtx.Channel.ID, text)
		if err != nil {
			result.Result = append(result.Result, err.Error())
			return result, nil
		}

		if parseCtx.Channel.ID == parseCtx.Sender.ID {
			// Update channel settings
			channelSettings.Voice = wantedVoice.Name
			err := parseCtx.Services.TTSService.UpdateChannelSettings(
				ctx,
				parseCtx.Channel.ID,
				channelSettings,
			)
			if err != nil {
				return nil, &types.CommandHandlerError{
					Message: "error while updating settings",
					Err:     err,
				}
			}
		} else {
			// Update user settings
			if userSettings == nil {
				_, err := parseCtx.Services.TTSService.CreateUserSettings(
					ctx,
					parseCtx.Channel.ID,
					parseCtx.Sender.ID,
					50,
					50,
					wantedVoice.Name,
				)
				if err != nil {
					return nil, &types.CommandHandlerError{
						Message: "error while creating user settings",
						Err:     err,
					}
				}
			} else {
				userSettings.Voice = wantedVoice.Name
				err := parseCtx.Services.TTSService.UpdateUserSettings(
					ctx,
					parseCtx.Channel.ID,
					parseCtx.Sender.ID,
					userSettings,
				)
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
