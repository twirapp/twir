package twitch

import (
	"context"
	"github.com/satont/tsuwari/apps/api-twirp/internal/impl_protected/deps"
	"github.com/satont/tsuwari/libs/grpc/generated/api/twitch"
)

type Twitch struct {
	*deps.Deps
}

func (c *Twitch) TwitchSearchUsers(ctx context.Context, req *twitch.TwitchSearchUsersRequest) (*twitch.TwitchSearchUsersResponse, error) {
	return nil, nil
}
