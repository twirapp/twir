package users_with_channel

import (
	"context"

	"github.com/twirapp/twir/libs/repositories/users-with-channel/model"
)

type Repository interface {
	GetByID(ctx context.Context, id string) (model.UserWithChannel, error)
	GetManyByIDS(ctx context.Context, input GetManyInput) ([]model.UserWithChannel, error)
	GetManyCount(ctx context.Context, input GetManyInput) (int, error)
}

type GetManyInput struct {
	Page         int
	PerPage      int
	IDs          []string
	HasBadgesIDS []string

	ChannelEnabled    *bool
	ChannelIsBotAdmin *bool
	ChannelIsBanned   *bool
}
