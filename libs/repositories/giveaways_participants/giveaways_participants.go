package giveaways_participants

import (
	"context"

	"github.com/oklog/ulid/v2"
	"github.com/twirapp/twir/libs/repositories/giveaways_participants/model"
)

type Repository interface {
	GetManyByGiveawayID(
		ctx context.Context,
		giveawayID ulid.ULID,
	) ([]model.ChannelGiveawayParticipant, error)
	GetByID(ctx context.Context, id ulid.ULID) (model.ChannelGiveawayParticipant, error)
	Create(ctx context.Context, input CreateInput) (model.ChannelGiveawayParticipant, error)
	Delete(ctx context.Context, id ulid.ULID) error
	Update(
		ctx context.Context,
		id ulid.ULID,
		input UpdateInput,
	) (model.ChannelGiveawayParticipant, error)
	GetWinnerForGiveaway(
		ctx context.Context,
		giveawayID ulid.ULID,
	) (model.ChannelGiveawayParticipant, error)
	ResetWinners(
		ctx context.Context, participantsIds ...ulid.ULID,
	) error
}

type CreateInput struct {
	GiveawayID  ulid.ULID
	DisplayName string
	UserID      string
}

type UpdateInput struct {
	IsWinner *bool
}
