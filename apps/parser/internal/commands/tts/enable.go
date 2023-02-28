package tts

import (
	"github.com/guregu/null"
	"github.com/satont/tsuwari/apps/parser/internal/types"
	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"
	model "github.com/satont/tsuwari/libs/gomodels"
)

var EnableCommand = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "tts enable",
		Description: null.StringFrom("Enable tts."),
		Module:      "TTS",
		IsReply:     true,
	},
	Handler: func(ctx *variables_cache.ExecutionContext) *types.CommandsHandlerResult {
		result := &types.CommandsHandlerResult{}
		err := switchEnableState(ctx.ChannelId, true)
		if err != nil {
			result.Result = append(result.Result, "Error while updating settings")
			return result
		}

		result.Result = append(result.Result, "TTS enabled")

		return result
	},
}
