package roles

import (
	"context"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/repositories/roles/model"
)

type Repository interface {
	GetManyByIDS(ctx context.Context, ids []uuid.UUID) ([]model.Role, error)
	GetManyByChannelID(ctx context.Context, channelID string) ([]model.Role, error)
	Create(ctx context.Context, input CreateInput) (model.Role, error)
	Update(ctx context.Context, id uuid.UUID, input UpdateInput) (model.Role, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (model.Role, error)
}

type CreateInput struct {
	ChannelID                 string
	Name                      string
	Type                      model.ChannelRoleEnum
	Permissions               []string
	RequiredWatchTime         int64
	RequiredMessages          int32
	RequiredUsedChannelPoints int64
}

type UpdateInput struct {
	Name                      *string
	Type                      *model.ChannelRoleEnum
	Permissions               []string
	RequiredWatchTime         *int64
	RequiredMessages          *int32
	RequiredUsedChannelPoints *int64
}
