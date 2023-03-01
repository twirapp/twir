package tts

import (
	"fmt"
	"strconv"

	"github.com/guregu/null"
	"github.com/satont/tsuwari/apps/parser/internal/types"
	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"
	model "github.com/satont/tsuwari/libs/gomodels"
)

var VolumeCommand = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "tts volume",
		Description: null.StringFrom("Change tts volume. This is not per user, it's global for the channel."),
		Module:      "TTS",
		IsReply:     true,
	},
	Handler: func(ctx *variables_cache.ExecutionContext) *types.CommandsHandlerResult {
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
			fmt.Println(err)
			result.Result = append(result.Result, "Error while updating settings")
			return result
		}

		result.Result = append(result.Result, fmt.Sprintf("Volume changed to %v", volume))

		return result
	},
}
