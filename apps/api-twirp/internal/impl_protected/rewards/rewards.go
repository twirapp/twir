package rewards

import (
	"context"
	"github.com/satont/tsuwari/apps/api-twirp/internal/impl_protected/deps"
	"github.com/satont/tsuwari/libs/grpc/generated/api/rewards"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Rewards struct {
	*deps.Deps
}

func (c *Rewards) RewardsGet(ctx context.Context, req *emptypb.Empty) (*rewards.GetResponse, error) {
	return nil, nil
}
