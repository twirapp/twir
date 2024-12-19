package commands_with_groups_and_responses

import (
	"context"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/repositories/commands_with_groups_and_responses/model"
)

type Repository interface {
	GetManyByChannelID(ctx context.Context, channelID string) (
		[]model.CommandWithGroupAndResponses,
		error,
	)
	GetByID(ctx context.Context, id uuid.UUID) (model.CommandWithGroupAndResponses, error)
}
