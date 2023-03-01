package tts

import (
	"fmt"
	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/parser/internal/types"
	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"
	"go.uber.org/zap"
	"strconv"
)

var VolumeCommand = types.DefaultCommand{
	Command: types.Command{
		Name:        "tts volume",
		Description: lo.ToPtr("Change tts volume. This is not per user, it's global for the channel."),
		Permission:  "BROADCASTER",
		Visible:     true,
		Module:      lo.ToPtr("TTS"),
		IsReply:     true,
	},
	Handler: func(ctx variables_cache.ExecutionContext) *types.CommandsHandlerResult {
		result := &types.CommandsHandlerResult{}
		channelSettings, channelModele := getSettings(ctx.ChannelId, "")

		if channelSettings == nil {
			result.Result = append(result.Result, "TTS is not configured.")
			return result
		}

		if ctx.Text == nil {
			result.Result = append(
				result.Result,
				fmt.Sprintf("Global volume: %v", channelSettings.Volume))
			return result
		}

		volume, err := strconv.Atoi(*ctx.Text)
		if err != nil {
			result.Result = append(result.Result, "Volume must be a number")
			return result
		}

		if volume < 0 || volume > 100 {
			result.Result = append(result.Result, "Volume must be between 0 and 100")
			return result
		}

		channelSettings.Volume = volume
		err = updateSettings(channelModele, channelSettings)
		if err != nil {
			zap.S().Error(err)
			result.Result = append(result.Result, "Error while updating settings")
			return result
		}

		result.Result = append(result.Result, fmt.Sprintf("Volume changed to %v", volume))

		return result
	},
}
