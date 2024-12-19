package roles_users

import (
	"context"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/repositories/roles_users/model"
)

type Repository interface {
	GetManyByRoleID(ctx context.Context, roleID uuid.UUID) ([]model.RoleUser, error)
	Create(ctx context.Context, input CreateInput) (model.RoleUser, error)
	CreateMany(ctx context.Context, inputs []CreateInput) ([]model.RoleUser, error)
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteManyByRoleID(ctx context.Context, roleID uuid.UUID) error
}

type CreateInput struct {
	UserID string
	RoleID uuid.UUID
}
