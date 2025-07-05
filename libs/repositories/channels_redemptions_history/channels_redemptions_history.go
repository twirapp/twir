package channelsredemptionshistory

import (
	"context"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/repositories/channels_redemptions_history/model"
)

type Repository interface {
	Create(ctx context.Context, input CreateInput) error
	CreateMany(ctx context.Context, input []CreateInput) error
	GetMany(ctx context.Context, input GetManyInput) (GetManyPayload, error)
	Count(ctx context.Context, input CountInput) (uint64, error)
}

type CreateInput struct {
	ChannelID    string
	UserID       string
	RewardID     uuid.UUID
	RewardPrompt *string
	RewardTitle  string
	RewardCost   int
}

type GetManyInput struct {
	ChannelID  string
	Page       int
	PerPage    int
	UserIDs    []string
	RewardsIDs []string
}

type GetManyPayload struct {
	Items []model.ChannelsRedemptionHistoryItem
	Total uint64
}

type CountInput struct {
	ChannelID  *string
	RewardsIDs []string
}
