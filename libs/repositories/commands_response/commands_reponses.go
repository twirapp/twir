package commands_response

import (
	"context"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/repositories/commands_response/model"
)

type Repository interface {
	// GetManyByIDs GetManyByChannelID returns groups in same order as requested
	GetManyByIDs(ctx context.Context, commandsIDs []uuid.UUID) ([]model.Response, error)
	Create(ctx context.Context, input CreateInput) (model.Response, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, id uuid.UUID, input UpdateInput) (model.Response, error)
}

type CreateInput struct {
	CommandID         uuid.UUID
	Text              *string
	Order             int
	TwitchCategoryIDs []string
	OnlineOnly        bool
	OfflineOnly       bool
}

type UpdateInput struct {
	Text              *string
	Order             *int
	TwitchCategoryIDs []string
	OnlineOnly        *bool
	OfflineOnly       *bool
}
