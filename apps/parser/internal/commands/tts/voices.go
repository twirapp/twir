package tts

import (
	"context"
	"fmt"
	"strings"

	"github.com/guregu/null"
	model "github.com/satont/twir/libs/gomodels"

	"github.com/samber/lo"
	"github.com/satont/twir/apps/parser/internal/types"
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

		channelSettings, _ := getSettings(ctx, parseCtx.Services.Gorm, parseCtx.Channel.ID, "")
		if channelSettings == nil {
			return result, nil
		}

		voices := getVoices(ctx, parseCtx.Services.Config)
		if len(voices) == 0 {
			return result, nil
		}

		if channelSettings != nil && len(channelSettings.DisallowedVoices) > 0 {
			voices = lo.Filter(
				voices, func(item Voice, _ int) bool {
					return !lo.Contains(channelSettings.DisallowedVoices, item.Name)
				},
			)
		}

		mapped := lo.Map(
			voices, func(item Voice, _ int) string {
				return fmt.Sprintf("%s (%s)", item.Name, item.Country)
			},
		)

		result.Result = append(result.Result, strings.Join(mapped, " · "))

		return result, nil
	},
}
