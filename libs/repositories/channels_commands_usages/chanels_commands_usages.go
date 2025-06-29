package channelscommandsusages

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, input CreateInput) error
	Count(ctx context.Context, input CountInput) (uint64, error)
}

type CreateInput struct {
	ChannelID string
	UserID    string
	CommandID uuid.UUID
}

type CountInput struct {
	ChannelID *string
	UserID    *string
	CommandID *uuid.UUID
}
