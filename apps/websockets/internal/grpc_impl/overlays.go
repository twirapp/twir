package grpc_impl

import (
	"context"

	"github.com/satont/twir/libs/grpc/generated/websockets"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (c *GrpcImpl) RefreshOverlays(ctx context.Context, req *websockets.RefreshOverlaysRequest) (
	*emptypb.Empty, error,
) {
	if err := c.sockets.OverlaysRegistry.SendEvent(
		req.ChannelId,
		"refreshOverlays",
		nil,
	); err != nil {
		c.services.Logger.Error(err)
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
