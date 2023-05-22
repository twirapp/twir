package tts

import (
	"context"

	"github.com/guregu/null"
	"github.com/satont/tsuwari/apps/parser/internal/types"

	model "github.com/satont/tsuwari/libs/gomodels"
)

var EnableCommand = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "tts enable",
		Description: null.StringFrom("Enable tts."),
		Module:      "TTS",
		IsReply:     true,
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) *types.CommandsHandlerResult {
		result := &types.CommandsHandlerResult{}
		err := switchEnableState(ctx, parseCtx.Services.Gorm, parseCtx.Channel.ID, true)
		if err != nil {
			result.Result = append(result.Result, "Error while updating settings")
			return result
		}

		result.Result = append(result.Result, "TTS enabled")

		return result
	},
}
