package tts

import (
	"fmt"
	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/parser/internal/types"
	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"
	"strings"
)

var VoicesCommand = types.DefaultCommand{
	Command: types.Command{
		Name:        "tts voices",
		Description: lo.ToPtr("List available voices"),
		Permission:  "VIEWER",
		Visible:     true,
		Module:      lo.ToPtr("TTS"),
		IsReply:     true,
	},
	Handler: func(ctx variables_cache.ExecutionContext) *types.CommandsHandlerResult {
		result := &types.CommandsHandlerResult{}

		voices := getVoices()
		if len(voices) == 0 {
			result.Result = append(result.Result, "No voices found")
			return result
		}

		mapped := lo.Map(voices, func(item Voice, _ int) string {
			return fmt.Sprintf("%s (%s)", item.Name, item.Country)
		})
		result.Result = append(result.Result, strings.Join(mapped, " · "))

		return result
	},
}
