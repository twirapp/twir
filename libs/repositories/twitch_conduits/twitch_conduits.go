package twitchconduits

import (
	"context"

	"github.com/twirapp/twir/libs/repositories/twitch_conduits/model"
)

type Repository interface {
	GetOne(ctx context.Context) (model.Conduit, error)
	Create(ctx context.Context, input CreateInput) (model.Conduit, error)
	Update(ctx context.Context, id string, input UpdateInput) (model.Conduit, error)
	Delete(ctx context.Context, id string) error
	DeleteAll(ctx context.Context) error
}

type CreateInput struct {
	ID         string
	ShardCount int8
}

type UpdateInput struct {
	ShardCount int8
}
