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
	ttsPitchArgName = "pitch"
)

var PitchCommand = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "tts pitch",
		Description: null.StringFrom("Change tts pitch"),
		Module:      "TTS",
		IsReply:     true,
	},
	SkipToxicityCheck: true,
	Args: []command_arguments.Arg{
		command_arguments.Int{
			Name:     ttsPitchArgName,
			Min:      lo.ToPtr(1),
			Max:      lo.ToPtr(100),
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

		pitchArg := parseCtx.ArgsParser.Get(ttsPitchArgName)

		if pitchArg == nil {
			result.Result = append(
				result.Result,
				fmt.Sprintf(
					"Global pitch: %v | Your pitch: %v",
					channelSettings.Pitch,
					lo.IfF(
						userSettings != nil, func() int {
							return userSettings.Pitch
						},
					).Else(channelSettings.Pitch),
				),
			)
			return result, nil
		}

		pitch := pitchArg.Int()

		if parseCtx.Channel.ID == parseCtx.Sender.ID {
			// Update channel settings
			channelSettings.Pitch = pitch
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
					pitch,
					channelSettings.Voice,
				)
				if err != nil {
					return nil, &types.CommandHandlerError{
						Message: "error while creating settings",
						Err:     err,
					}
				}
			} else {
				userSettings.Pitch = pitch
				err := parseCtx.Services.TTSService.UpdateUserSettings(
					ctx,
					parseCtx.Channel.ID,
					parseCtx.Sender.ID,
					userSettings,
				)
				if err != nil {
					return nil, &types.CommandHandlerError{
						Message: "error while updating settings",
						Err:     err,
					}
				}
			}
		}

		result.Result = append(result.Result, fmt.Sprintf("Pitch changed to %v", pitch))

		parseCtx.Services.TTSCache.Invalidate(ctx, parseCtx.Channel.ID)

		return result, nil
	},
}
