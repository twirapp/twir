package channels

import (
	"context"

	"github.com/twirapp/twir/libs/repositories/channels/model"
)

type Repository interface {
	GetMany(ctx context.Context, input GetManyInput) ([]model.Channel, error)
}

type GetManyInput struct {
	Enabled *bool
	Limit   int
	Page    int
}
