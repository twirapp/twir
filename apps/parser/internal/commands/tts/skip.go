package tts

import (
	"context"
	"github.com/samber/do"
	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/parser/internal/di"
	"github.com/satont/tsuwari/apps/parser/internal/types"
	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"
	"github.com/satont/tsuwari/libs/grpc/generated/websockets"
	"go.uber.org/zap"
)

var SkipCommand = types.DefaultCommand{
	Command: types.Command{
		Name:        "tts skip",
		Description: lo.ToPtr("Skip current saying message in tts"),
		Permission:  "MODERATOR",
		Visible:     true,
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
			zap.S().Error(err)
		}

		return result
	},
}
