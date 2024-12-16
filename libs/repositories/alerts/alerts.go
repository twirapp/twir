package alerts

import (
	"context"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/repositories/alerts/model"
)

type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (model.Alert, error)
	GetManyByChannelID(ctx context.Context, channelID string) ([]model.Alert, error)
	Create(ctx context.Context, input CreateInput) (model.Alert, error)
	Update(ctx context.Context, id uuid.UUID, input UpdateInput) (model.Alert, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type CreateInput struct {
	Name         string
	ChannelID    string
	AudioID      *string
	AudioVolume  int
	CommandIDS   []string
	RewardIDS    []string
	GreetingsIDS []string
	KeywordsIDS  []string
}

type UpdateInput struct {
	Name         *string
	AudioID      *string
	AudioVolume  *int
	CommandIDS   []string
	RewardIDS    []string
	GreetingsIDS []string
	KeywordsIDS  []string
}
