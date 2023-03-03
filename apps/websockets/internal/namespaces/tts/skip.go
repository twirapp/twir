package tts

import (
	"context"
	"encoding/json"
	"github.com/olahol/melody"
	"github.com/satont/tsuwari/apps/websockets/types"
	"github.com/satont/tsuwari/libs/grpc/generated/websockets"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (c *TTS) Skip(ctx context.Context, msg *websockets.TTSSkipMessage) (*emptypb.Empty, error) {
	message := &types.WebSocketMessage{
		EventName: "skip",
		Data:      msg,
	}

	bytes, err := json.Marshal(message)
	if err != nil {
		c.services.Logger.Error(err)
		return nil, err
	}

	err = c.manager.BroadcastFilter(bytes, func(session *melody.Session) bool {
		userId, ok := session.Get("userId")
		return ok && userId.(string) == msg.ChannelId
	})

	if err != nil {
		c.services.Logger.Error(err)
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
