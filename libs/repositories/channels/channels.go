package channels

import (
	"context"

	"github.com/twirapp/twir/libs/repositories/channels/model"
)

type Repository interface {
	GetMany(ctx context.Context, input GetManyInput) ([]model.Channel, error)
	GetByID(ctx context.Context, channelID string) (model.Channel, error)
	GetCount(ctx context.Context, input GetCountInput) (int, error)
}

type GetManyInput struct {
	Enabled *bool
	PerPage int
	Page    int
}

type GetCountInput struct {
	OnlyEnabled bool
}
