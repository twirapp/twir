package obs

import (
	"context"

	"github.com/satont/twir/libs/grpc/generated/websockets"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (c *OBS) SetScene(_ context.Context, msg *websockets.ObsSetSceneMessage) (*emptypb.Empty, error) {
	return nil, nil
}
