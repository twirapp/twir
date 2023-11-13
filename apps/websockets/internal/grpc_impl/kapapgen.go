package grpc_impl

import (
	"context"

	"github.com/satont/twir/libs/grpc/generated/websockets"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (c *GrpcImpl) TriggerKappagen(
	_ context.Context,
	msg *websockets.TriggerKappagenRequest,
) (*emptypb.Empty, error) {
	if err := c.kappagenServer.SendEvent(msg.ChannelId, "kappagen", msg); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
