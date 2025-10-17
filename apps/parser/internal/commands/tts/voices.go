package tts

import (
	"context"
	"fmt"
	"strings"

	"github.com/guregu/null"
	"github.com/twirapp/twir/apps/parser/internal/services/tts"
	"github.com/twirapp/twir/apps/parser/locales"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/i18n"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
)

var VoicesCommand = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "tts voices",
		Description: null.StringFrom("List available voices"),
		Module:      "TTS",
		IsReply:     true,
	},
	SkipToxicityCheck: true,
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		result := &types.CommandsHandlerResult{}

		voices, err := parseCtx.Services.TTSService.GetFilteredVoices(ctx, parseCtx.Channel.ID)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(ctx, locales.Translations.Commands.Tts.Errors.WhileGettingVoices),
				Err:     err,
			}
		}

		if len(voices) == 0 {
			result.Result = []string{i18n.GetCtx(ctx, locales.Translations.Commands.Tts.Info.NoVoices)}
			return result, nil
		}

		mapped := lo.Map(
			voices, func(item tts.Voice, _ int) string {
				return fmt.Sprintf("%s (%s)", item.Name, item.Country)
			},
		)

		result.Result = append(result.Result, strings.Join(mapped, " · "))

		return result, nil
	},
}
