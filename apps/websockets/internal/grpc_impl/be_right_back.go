package grpc_impl

import (
	"context"

	"github.com/satont/twir/libs/grpc/generated/websockets"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (c *GrpcImpl) RefreshBrbSettings(
	ctx context.Context,
	req *websockets.RefreshBrbSettingsRequest,
) (*emptypb.Empty, error) {
	panic("not implemented")
}
