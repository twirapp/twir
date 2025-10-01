package tts

import (
	"context"

	"github.com/lib/pq"

	"github.com/guregu/null"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"

	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/i18n"
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
		err := parseCtx.Services.TTSService.ToggleChannelEnabled(ctx, parseCtx.Channel.ID, true)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(ctx, locales.Translations.Commands.Tts.Errors.WhileEnable),
				Err:     err,
			}
		}

		result.Result = append(result.Result, i18n.GetCtx(ctx, locales.Translations.Commands.Tts.Info.Enabled))

		parseCtx.Services.TTSCache.Invalidate(ctx, parseCtx.Channel.ID)

		return result, nil
	},
}
