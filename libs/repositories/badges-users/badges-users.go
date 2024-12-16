package badges_users

import (
	"context"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/repositories/badges-users/model"
)

type Repository interface {
	GetMany(ctx context.Context, input GetManyInput) ([]model.BadgeUser, error)
}

type GetManyInput struct {
	BadgeID uuid.UUID
}
