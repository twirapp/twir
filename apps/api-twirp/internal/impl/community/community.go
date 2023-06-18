package community

import (
	"context"
	"github.com/satont/tsuwari/apps/api-twirp/internal/impl/deps"
	"github.com/satont/tsuwari/libs/grpc/generated/api/community"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Community struct {
	*deps.Deps
}

func (c *Community) CommunityGetUsers(ctx context.Context, request *community.GetUsersRequest) (*community.GetUsersResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Community) CommunityResetStats(ctx context.Context, request *community.ResetStatsRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}
