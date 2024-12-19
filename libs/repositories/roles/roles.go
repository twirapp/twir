package roles

import (
	"context"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/repositories/roles/model"
)

type Repository interface {
	GetManyByIDS(ctx context.Context, ids []uuid.UUID) ([]model.Role, error)
	GetManyByChannelID(ctx context.Context, channelID string) ([]model.Role, error)
}
