package grpc_impl

import (
	"context"

	"github.com/twirapp/twir/libs/grpc/websockets"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (c *GrpcImpl) TextToSpeechSay(_ context.Context, msg *websockets.TTSMessage) (
	*emptypb.Empty,
	error,
) {
	if err := c.ttsServer.SendEvent(msg.ChannelId, "say", msg); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
func (c *GrpcImpl) TextToSpeechSkip(
	_ context.Context,
	msg *websockets.TTSSkipMessage,
) (*emptypb.Empty, error) {
	if err := c.ttsServer.SendEvent(msg.ChannelId, "skip", msg); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
