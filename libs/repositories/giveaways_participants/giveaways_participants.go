package giveaways_participants

import (
	"context"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/repositories/giveaways_participants/model"
)

type Repository interface {
	GetManyByGiveawayID(
		ctx context.Context,
		giveawayID uuid.UUID,
		input GetManyInput,
	) ([]model.ChannelGiveawayParticipant, error)
	GetByID(ctx context.Context, id uuid.UUID) (model.ChannelGiveawayParticipant, error)
	Create(ctx context.Context, input CreateInput) (model.ChannelGiveawayParticipant, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Update(
		ctx context.Context,
		id uuid.UUID,
		input UpdateInput,
	) (model.ChannelGiveawayParticipant, error)
	GetWinnerForGiveaway(
		ctx context.Context,
		giveawayID uuid.UUID,
	) (model.ChannelGiveawayParticipant, error)
	ResetWinners(
		ctx context.Context, participantsIds ...uuid.UUID,
	) error
}

type CreateInput struct {
	GiveawayID      uuid.UUID
	UserID          uuid.UUID
	UserLogin       string
	UserDisplayName string
	IsWinner        bool
}

type GetManyInput struct {
	OnlyWinners   bool
	IgnoreWinners bool
}

type UpdateInput struct {
	IsWinner *bool
}
