package roles

import (
	"context"
	"github.com/satont/tsuwari/apps/api-twirp/internal/impl_protected/deps"
	"github.com/satont/tsuwari/libs/grpc/generated/api/roles"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Roles struct {
	*deps.Deps
}

func (c *Roles) RolesGetAll(ctx context.Context, empty *emptypb.Empty) (*roles.GetAllResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Roles) RolesUpdate(ctx context.Context, request *roles.UpdateRequest) (*roles.Role, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Roles) RolesDelete(ctx context.Context, request *roles.DeleteRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Roles) RolesCreate(ctx context.Context, request *roles.CreateRequest) (*roles.Role, error) {
	//TODO implement me
	panic("implement me")
}
