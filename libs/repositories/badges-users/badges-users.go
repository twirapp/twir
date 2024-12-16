package badges_users

import (
	"context"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/repositories/badges-users/model"
)

type Repository interface {
	GetMany(ctx context.Context, input GetManyInput) ([]model.BadgeUser, error)
	Create(ctx context.Context, input CreateInput) (model.BadgeUser, error)
	Delete(ctx context.Context, input DeleteInput) error
}

type GetManyInput struct {
	BadgeID uuid.UUID
}

type CreateInput struct {
	BadgeID uuid.UUID
	UserID  string
}

type DeleteInput struct {
	BadgeID uuid.UUID
	UserID  string
}
