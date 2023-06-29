package modules

import (
	"context"
	"github.com/satont/twir/libs/grpc/generated/api/modules_obs_websocket"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (c *Modules) ModulesOBSWebsocketGet(ctx context.Context, empty *emptypb.Empty) (*modules_obs_websocket.GetResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Modules) ModulesOBSWebsocketUpdate(ctx context.Context, request *modules_obs_websocket.PostRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}
