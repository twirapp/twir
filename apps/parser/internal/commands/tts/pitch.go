package tts

import (
	"context"
	"fmt"

	"github.com/guregu/null"
	command_arguments "github.com/satont/twir/apps/parser/internal/command-arguments"
	model "github.com/satont/twir/libs/gomodels"

	"github.com/samber/lo"
	"github.com/satont/twir/apps/parser/internal/types"
	"go.uber.org/zap"
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
		channelSettings, channelModele := getSettings(
			ctx,
			parseCtx.Services.Gorm,
			parseCtx.Channel.ID,
			"",
		)

		if channelSettings == nil {
			return result, nil
		}

		userSettings, currentUserModel := getSettings(
			ctx,
			parseCtx.Services.Gorm,
			parseCtx.Channel.ID,
			parseCtx.Sender.ID,
		)

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
			channelSettings.Pitch = pitch
			err := updateSettings(ctx, parseCtx.Services.Gorm, channelModele, channelSettings)
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
					parseCtx.Services.Gorm,
					pitch,
					50,
					channelSettings.Voice,
					parseCtx.Channel.ID,
					parseCtx.Sender.ID,
				)
				if err != nil {
					zap.S().Error(err)
					result.Result = append(result.Result, "error while updating settings")
					return result, nil
				}
			} else {
				userSettings.Pitch = pitch
				err := updateSettings(ctx, parseCtx.Services.Gorm, currentUserModel, userSettings)
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
