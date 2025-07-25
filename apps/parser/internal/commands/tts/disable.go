package tts

import (
	"context"

	"github.com/lib/pq"

	"github.com/guregu/null"
	"github.com/twirapp/twir/apps/parser/internal/types"
	model "github.com/twirapp/twir/libs/gomodels"
)

var DisableCommand = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "tts disable",
		Description: null.StringFrom("Disable tts."),
		Module:      "TTS",
		IsReply:     true,
		RolesIDS:    pq.StringArray{model.ChannelRoleTypeBroadcaster.String()},
	},
	SkipToxicityCheck: true,
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		result := &types.CommandsHandlerResult{}
		err := switchEnableState(ctx, parseCtx.Services.Gorm, parseCtx.Channel.ID, false)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: "error while disabling tts",
				Err:     err,
			}
		}

		result.Result = append(result.Result, "TTS disabled")

		parseCtx.Services.TTSCache.Invalidate(ctx, parseCtx.Channel.ID)

		return result, nil
	},
}
