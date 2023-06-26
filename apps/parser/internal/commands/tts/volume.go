package tts

import (
	"context"
	"fmt"
	"strconv"

	"github.com/guregu/null"
	"github.com/satont/twir/apps/parser/internal/types"
	model "github.com/satont/twir/libs/gomodels"
	"go.uber.org/zap"
)

var VolumeCommand = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "tts volume",
		Description: null.StringFrom("Change tts volume. This is not per user, it's global for the channel."),
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

		if parseCtx.Text == nil {
			result.Result = append(
				result.Result,
				fmt.Sprintf("Global volume: %v", channelSettings.Volume),
			)
			return result
		}

		volume, err := strconv.Atoi(*parseCtx.Text)
		if err != nil {
			result.Result = append(result.Result, "Volume must be a number")
			return result
		}

		if volume < 0 || volume > 100 {
			result.Result = append(result.Result, "Volume must be between 0 and 100")
			return result
		}

		channelSettings.Volume = volume
		err = updateSettings(ctx, parseCtx.Services.Gorm, channelModele, channelSettings)
		if err != nil {
			zap.S().Error(err)
			result.Result = append(result.Result, "Error while updating settings")
			return result
		}

		result.Result = append(result.Result, fmt.Sprintf("Volume changed to %v", volume))

		return result
	},
}
