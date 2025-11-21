package grpc_impl

import (
	"context"

	"github.com/twirapp/twir/libs/grpc/websockets"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (c *GrpcImpl) TriggerShowBrb(
	_ context.Context,
	req *websockets.TriggerShowBrbRequest,
) (*emptypb.Empty, error) {
	// err := c.beRightBackServer.SendEvent(
	// 	req.GetChannelId(),
	// 	"start", map[string]any{
	// 		"minutes": req.GetMinutes(),
	// 		"text":    req.Text,
	// 	},
	// )

	return nil, nil
}

func (c *GrpcImpl) TriggerHideBrb(
	_ context.Context,
	req *websockets.TriggerHideBrbRequest,
) (
	*emptypb.Empty, error,
) {
	return nil, nil
	// err := c.beRightBackServer.SendEvent(req.GetChannelId(), "stop", nil)
	//
	// return &emptypb.Empty{}, err
}
