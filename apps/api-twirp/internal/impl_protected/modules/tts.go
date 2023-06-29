package modules

import (
	"context"
	"github.com/satont/twir/libs/grpc/generated/api/modules_tts"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (c *Modules) ModulesTTSGet(ctx context.Context, empty *emptypb.Empty) (*modules_tts.GetResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Modules) ModulesTTSUpdate(ctx context.Context, request *modules_tts.PostRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Modules) ModulesTTSGetInfo(ctx context.Context, empty *emptypb.Empty) (*modules_tts.GetInfoResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Modules) ModulesTTSSay(ctx context.Context, request *modules_tts.SayRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}
