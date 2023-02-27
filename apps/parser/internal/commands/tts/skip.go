package tts

import (
	"context"
	"fmt"
	"github.com/samber/do"
	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/parser/internal/di"
	"github.com/satont/tsuwari/apps/parser/internal/types"
	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"
	"github.com/satont/tsuwari/libs/grpc/generated/websockets"
)

var SkipCommand = types.DefaultCommand{
	Command: types.Command{
		Name:        "tts skip",
		Description: lo.ToPtr("Skip current saying message in tts"),
		Permission:  "MODERATOR",
		Visible:     false,
		Module:      lo.ToPtr("TTS"),
		IsReply:     true,
	},
	Handler: func(ctx variables_cache.ExecutionContext) *types.CommandsHandlerResult {
		webSocketsGrpc := do.MustInvoke[websockets.WebsocketClient](di.Provider)

		result := &types.CommandsHandlerResult{}

		_, err := webSocketsGrpc.TextToSpeechSkip(context.Background(), &websockets.TTSSkipMessage{
			ChannelId: ctx.ChannelId,
		})

		if err != nil {
			fmt.Println(err)
		}

		return result
	},
}
