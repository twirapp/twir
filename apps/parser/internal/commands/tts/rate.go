package tts

import (
	"context"
	"fmt"
	"strconv"

	"github.com/guregu/null"
	model "github.com/satont/twir/libs/gomodels"

	"github.com/samber/lo"
	"github.com/satont/twir/apps/parser/internal/types"
)

var RateCommand = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "tts rate",
		Description: null.StringFrom("Change tts rate"),
		Module:      "TTS",
		IsReply:     true,
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

		if parseCtx.Text == nil {
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

		rate, err := strconv.Atoi(*parseCtx.Text)
		if err != nil {
			result.Result = append(result.Result, "rate must be a number")
			return result, nil
		}

		if rate < 0 || rate > 100 {
			result.Result = append(result.Result, "rate must be between 0 and 100")
			return result, nil
		}

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

		return result, nil
	},
}
