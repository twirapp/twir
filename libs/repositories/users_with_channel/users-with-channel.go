package users_with_channel

import (
	"context"

	"github.com/google/uuid"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	"github.com/twirapp/twir/libs/repositories/users_with_channel/model"
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
	SearchQuery  string
	Platforms    []platformentity.Platform
	HasBadgesIDS []uuid.UUID

	ChannelEnabled    *bool
	ChannelIsBotAdmin *bool
	IsBanned          *bool
}
