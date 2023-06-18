package timers

import (
	"context"
	"github.com/satont/tsuwari/apps/api-twirp/internal/impl_deps"
	"github.com/satont/tsuwari/libs/grpc/generated/api/timers"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Timers struct {
	*impl_deps.Deps
}

func (c *Timers) TimersGet(ctx context.Context, empty *emptypb.Empty) (*timers.GetResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Timers) TimersUpdate(ctx context.Context, request *timers.UpdateRequest) (*timers.Timer, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Timers) TimersDelete(ctx context.Context, request *timers.DeleteRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Timers) TimersCreate(ctx context.Context, request *timers.CreateRequest) (*timers.Timer, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Timers) TimersEnableOrDisable(ctx context.Context, request *timers.PatchRequest) (*timers.Timer, error) {
	//TODO implement me
	panic("implement me")
}
