package grpc_impl

import (
	"context"

	"github.com/satont/twir/libs/grpc/generated/websockets"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (c *GrpcImpl) YoutubeAddSongToQueue(
	ctx context.Context, msg *websockets.YoutubeAddSongToQueueRequest,
) (*emptypb.Empty, error) {
	return c.sockets.YouTube.AddSongToQueue(ctx, msg)
}
func (c *GrpcImpl) YoutubeRemoveSongToQueue(
	ctx context.Context, msg *websockets.YoutubeRemoveSongFromQueueRequest,
) (*emptypb.Empty, error) {
	return c.sockets.YouTube.RemoveSongFromQueue(ctx, msg)
}
