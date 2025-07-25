package tts

import (
	"context"
	"encoding/json"
	"time"

	"github.com/olahol/melody"
	"github.com/twirapp/twir/apps/websockets/types"
	"github.com/twirapp/twir/libs/grpc/websockets"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (c *TTS) Skip(ctx context.Context, msg *websockets.TTSSkipMessage) (*emptypb.Empty, error) {
	message := &types.WebSocketMessage{
		EventName: "skip",
		Data:      msg,
		CreatedAt: time.Now().UTC().String(),
	}

	bytes, err := json.Marshal(message)
	if err != nil {
		c.logger.Error(err.Error())
		return nil, err
	}

	err = c.manager.BroadcastFilter(
		bytes, func(session *melody.Session) bool {
			userId, ok := session.Get("userId")
			return ok && userId.(string) == msg.ChannelId
		},
	)

	if err != nil {
		c.logger.Error(err.Error())
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
