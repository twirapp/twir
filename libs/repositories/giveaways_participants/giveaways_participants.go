package giveaways_participants

import (
	"context"

	"github.com/twirapp/twir/libs/repositories/giveaways_participants/model"
)

type Repository interface {
	GetManyByGiveawayID(
		ctx context.Context,
		giveawayID string,
		input GetManyInput,
	) ([]model.ChannelGiveawayParticipant, error)
	GetByID(ctx context.Context, id string) (model.ChannelGiveawayParticipant, error)
	Create(ctx context.Context, input CreateInput) (model.ChannelGiveawayParticipant, error)
	Delete(ctx context.Context, id string) error
	Update(
		ctx context.Context,
		id string,
		input UpdateInput,
	) (model.ChannelGiveawayParticipant, error)
	GetWinnerForGiveaway(
		ctx context.Context,
		giveawayID string,
	) (model.ChannelGiveawayParticipant, error)
	ResetWinners(
		ctx context.Context, participantsIds ...string,
	) error
}

type CreateInput struct {
	GiveawayID      string
	UserID          string
	UserLogin       string
	UserDisplayName string
}

type GetManyInput struct {
	OnlyWinners   bool
	IgnoreWinners bool
}

type UpdateInput struct {
	IsWinner *bool
}
