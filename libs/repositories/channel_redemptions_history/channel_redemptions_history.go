package channelredemptionshistory

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	CreateMany(ctx context.Context, inputs []CreateInput) error
}

type CreateInput struct {
	ChannelID    string
	UserID       string
	RewardID     uuid.UUID
	RewardTitle  string
	RewardPrompt *string
	RewardCost   int
}
