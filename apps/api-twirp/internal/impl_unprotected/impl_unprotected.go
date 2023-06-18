package impl_unprotected

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/satont/tsuwari/libs/grpc/generated/api/stats"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type UnProtected struct{}

func (c *UnProtected) Stats(ctx context.Context, empty *emptypb.Empty) (*stats.Response, error) {
	return &stats.Response{
		Users:    0,
		Channels: 0,
		Commands: 0,
		Messages: 0,
	}, nil
}

func (c *UnProtected) Greet(ctx context.Context, empty *emptypb.Empty) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

type Opts struct {
	Redis *redis.Client
	DB    *gorm.DB
}

func New(opts Opts) *UnProtected {
	return &UnProtected{}
}
