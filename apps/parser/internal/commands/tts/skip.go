package tts

import (
	"context"
	"fmt"
	"github.com/guregu/null"
	model "github.com/satont/tsuwari/libs/gomodels"

	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/parser/internal/di"
	"github.com/satont/tsuwari/apps/parser/internal/types"
	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"
	"github.com/satont/tsuwari/libs/grpc/generated/websockets"
)

var SkipCommand = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "tts skip",
		Description: null.StringFrom("Skip current saying message in tts"),
		Module:      "TTS",
		IsReply:     true,
	},
	Handler: func(ctx *variables_cache.ExecutionContext) *types.CommandsHandlerResult {
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
