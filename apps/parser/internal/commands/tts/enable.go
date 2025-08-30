package tts

import (
	"context"

	"github.com/lib/pq"

	"github.com/guregu/null"
	"github.com/twirapp/twir/apps/parser/internal/types"

	model "github.com/twirapp/twir/libs/gomodels"
)

var EnableCommand = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "tts enable",
		Description: null.StringFrom("Enable tts."),
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
		err := switchEnableState(ctx, parseCtx.Services.TTSRepository, parseCtx.Channel.ID, true)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: "failed to enable tts",
				Err:     err,
			}
		}

		result.Result = append(result.Result, "TTS enabled")

		parseCtx.Services.TTSCache.Invalidate(ctx, parseCtx.Channel.ID)

		return result, nil
	},
}
