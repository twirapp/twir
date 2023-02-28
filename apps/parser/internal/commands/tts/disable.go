package tts

import (
	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/parser/internal/types"
	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"
)

var DisableCommand = types.DefaultCommand{
	Command: types.Command{
		Name:        "tts disable",
		Description: lo.ToPtr("Disable tts."),
		Permission:  "BROADCASTER",
		Visible:     false,
		Module:      lo.ToPtr("TTS"),
		IsReply:     true,
	},
	Handler: func(ctx variables_cache.ExecutionContext) *types.CommandsHandlerResult {
		result := &types.CommandsHandlerResult{}
		err := switchEnableState(ctx.ChannelId, false)
		if err != nil {
			result.Result = append(result.Result, "Error while updating settings")
			return result
		}

		result.Result = append(result.Result, "TTS disabled")

		return result
	},
}
