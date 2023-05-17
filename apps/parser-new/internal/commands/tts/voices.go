package tts

import (
	"fmt"
	"github.com/guregu/null"
	model "github.com/satont/tsuwari/libs/gomodels"
	"strings"

	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/parser/internal/types"
	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"
)

var VoicesCommand = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "tts voices",
		Description: null.StringFrom("List available voices"),
		Module:      "TTS",
		IsReply:     true,
	},
	Handler: func(ctx *variables_cache.ExecutionContext) *types.CommandsHandlerResult {
		result := &types.CommandsHandlerResult{}

		channelSettings, _ := getSettings(ctx.ChannelId, "")

		voices := getVoices()
		if len(voices) == 0 {
			result.Result = append(result.Result, "No voices found")
			return result
		}

		if channelSettings != nil && len(channelSettings.DisallowedVoices) > 0 {
			voices = lo.Filter(voices, func(item Voice, _ int) bool {
				return !lo.Contains(channelSettings.DisallowedVoices, item.Name)
			})
		}

		mapped := lo.Map(voices, func(item Voice, _ int) string {
			return fmt.Sprintf("%s (%s)", item.Name, item.Country)
		})

		result.Result = append(result.Result, strings.Join(mapped, " · "))

		return result
	},
}
