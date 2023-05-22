package tts

import (
	"context"
	"fmt"
	"strconv"

	"github.com/guregu/null"
	model "github.com/satont/tsuwari/libs/gomodels"

	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/parser/internal/types"
	"go.uber.org/zap"
)

var RateCommand = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "tts rate",
		Description: null.StringFrom("Change tts rate"),
		Module:      "TTS",
		IsReply:     true,
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) *types.CommandsHandlerResult {
		result := &types.CommandsHandlerResult{}
		channelSettings, channelModele := getSettings(ctx, parseCtx.Services.Gorm, parseCtx.Channel.ID, "")

		if channelSettings == nil {
			result.Result = append(result.Result, "TTS is not configured.")
			return result
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
					lo.IfF(userSettings != nil, func() int {
						return userSettings.Rate
					}).Else(channelSettings.Rate),
				))
			return result
		}

		rate, err := strconv.Atoi(*parseCtx.Text)
		if err != nil {
			result.Result = append(result.Result, "Rate must be a number")
			return result
		}

		if rate < 0 || rate > 100 {
			result.Result = append(result.Result, "Rate must be between 0 and 100")
			return result
		}

		if parseCtx.Channel.ID == parseCtx.Sender.ID {
			channelSettings.Rate = rate
			err := updateSettings(ctx, parseCtx.Services.Gorm, channelModele, channelSettings)
			if err != nil {
				zap.S().Error(err)
				result.Result = append(result.Result, "Error while updating settings")
				return result
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
					zap.S().Error(err)
					result.Result = append(result.Result, "Error while creating settings")
					return result
				}
			} else {
				userSettings.Rate = rate
				err := updateSettings(ctx, parseCtx.Services.Gorm, currentUserModel, userSettings)
				if err != nil {
					zap.S().Error(err)
					result.Result = append(result.Result, "Error while updating settings")
					return result
				}
			}
		}

		result.Result = append(result.Result, fmt.Sprintf("Rate changed to %v", rate))

		return result
	},
}
