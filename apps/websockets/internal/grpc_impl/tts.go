package grpc_impl

import (
	"context"
	"github.com/satont/tsuwari/libs/grpc/generated/websockets"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (c *grpcImpl) TextToSpeechSay(_ context.Context, msg *websockets.TTSMessage) (*emptypb.Empty, error) {
	if err := c.sockets.TTS.SendEvent(msg.ChannelId, "say", msg); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
func (c *grpcImpl) TextToSpeechSkip(_ context.Context, msg *websockets.TTSSkipMessage) (*emptypb.Empty, error) {
	if err := c.sockets.TTS.SendEvent(msg.ChannelId, "skip", msg); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
