package stats

import (
	"context"
	"github.com/satont/tsuwari/libs/grpc/generated/api/stats"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Stats struct {
}

func (c *Stats) Stats(ctx context.Context, empty *emptypb.Empty) (*stats.Response, error) {
	return &stats.Response{
		Users:    0,
		Channels: 0,
		Commands: 0,
		Messages: 0,
	}, nil
}
