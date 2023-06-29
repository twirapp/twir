package modules

import (
	"context"
	"github.com/satont/twir/libs/grpc/generated/api/modules_sr"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (c *Modules) ModulesSRGet(ctx context.Context, empty *emptypb.Empty) (*modules_sr.GetResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Modules) ModulesSRSearchVideosOrChannels(ctx context.Context, request *modules_sr.GetSearchRequest) (*modules_sr.GetSearchResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Modules) ModulesSRUpdate(ctx context.Context, request *modules_sr.PostRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}
