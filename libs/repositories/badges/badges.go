package badges

import (
	"context"

	"github.com/twirapp/twir/libs/repositories/badges/model"
)

type Repository interface {
	GetMany(ctx context.Context, input GetManyInput) ([]model.Badge, error)
}

type GetManyInput struct {
	Enabled bool
}
