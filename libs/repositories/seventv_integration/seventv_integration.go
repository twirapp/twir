package seventv_integration

import (
	"context"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/repositories/seventv_integration/model"
)

type Repository interface {
	GetByChannelID(ctx context.Context, channelID string) (model.SevenTvIntegration, error)
	Create(ctx context.Context, input CreateInput) error
	Update(ctx context.Context, id uuid.UUID, input UpdateInput) error
}

type UpdateInput struct {
	RewardIdForAddEmote        *string
	RewardIdForRemoveEmote     *string
	DeleteEmotesOnlyAddedByApp *bool
}

type CreateInput struct {
	ChannelID                  string
	RewardIdForAddEmote        *string
	RewardIdForRemoveEmote     *string
	DeleteEmotesOnlyAddedByApp *bool
}
