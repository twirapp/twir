package impl_unprotected

import (
	"context"
	"github.com/redis/go-redis/v9"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type UnProtected struct{}

func (u *UnProtected) Greet(ctx context.Context, empty *emptypb.Empty) (*emptypb.Empty, error) {
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
