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
	ttsRateArgName = "rate"
)

var RateCommand = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "tts rate",
		Description: null.StringFrom("Change tts rate"),
		Module:      "TTS",
		IsReply:     true,
	},
	SkipToxicityCheck: true,
	Args: []command_arguments.Arg{
		command_arguments.Int{
			Name:     ttsRateArgName,
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

		rateArg := parseCtx.ArgsParser.Get(ttsRateArgName)

		if rateArg == nil {
			result.Result = append(
				result.Result,
				fmt.Sprintf(
					"Global rate: %v | Your rate: %v",
					channelSettings.Rate,
					lo.IfF(
						userSettings != nil, func() int {
							return userSettings.Rate
						},
					).Else(channelSettings.Rate),
				),
			)
			return result, nil
		}

		rate := rateArg.Int()

		if parseCtx.Channel.ID == parseCtx.Sender.ID {
			channelSettings.Rate = rate
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
					rate,
					50,
					channelSettings.Voice,
					parseCtx.Channel.ID,
					parseCtx.Sender.ID,
				)
				if err != nil {
					return nil, &types.CommandHandlerError{
						Message: "error while creating settings",
						Err:     err,
					}
				}
			} else {
				userSettings.Rate = rate
				err := updateSettings(ctx, parseCtx.Services.Gorm, currentUserModel, userSettings)
				if err != nil {
					return nil, &types.CommandHandlerError{
						Message: "error while updating settings",
						Err:     err,
					}
				}
			}
		}

		result.Result = append(result.Result, fmt.Sprintf("Rate changed to %v", rate))

		parseCtx.Services.TTSCache.Invalidate(ctx, parseCtx.Channel.ID)

		return result, nil
	},
}
