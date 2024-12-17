package commands_group

import (
	"context"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/repositories/commands_group/model"
)

type Repository interface {
	// GetManyByIDs GetManyByChannelID returns groups in same order as requested
	GetManyByIDs(ctx context.Context, ids []uuid.UUID) ([]*model.Group, error)
}
