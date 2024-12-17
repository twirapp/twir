package commands_response

import (
	"context"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/repositories/commands_response/model"
)

type Repository interface {
	// GetManyByIDs GetManyByChannelID returns groups in same order as requested
	GetManyByIDs(ctx context.Context, commandsIDs []uuid.UUID) ([]model.Response, error)
}
