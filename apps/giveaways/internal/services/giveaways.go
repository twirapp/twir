package services

import (
	"context"
	"math/rand"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/oklog/ulid/v2"
	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/giveaways/internal/entity"
	"github.com/twirapp/twir/libs/repositories/giveaways"
	"github.com/twirapp/twir/libs/repositories/giveaways_participants"
	"github.com/twirapp/twir/libs/repositories/giveaways_participants/model"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	TxManager                       trm.Manager
	GiveawaysRepository             giveaways.Repository
	GiveawaysParticipantsRepository giveaways_participants.Repository
}

func New(opts Opts) *Service {
	return &Service{
		txManager:                       opts.TxManager,
		giveawaysRepository:             opts.GiveawaysRepository,
		giveawaysParticipantsRepository: opts.GiveawaysParticipantsRepository,
	}
}

type Service struct {
	txManager                       trm.Manager
	giveawaysRepository             giveaways.Repository
	giveawaysParticipantsRepository giveaways_participants.Repository
}

func (c *Service) TryAddParticipant(
	ctx context.Context,
	userID string,
	displayName string,
	giveawayID ulid.ULID,
) error {
	_, err := c.giveawaysParticipantsRepository.Create(ctx, giveaways_participants.CreateInput{
		GiveawayID:  giveawayID,
		DisplayName: displayName,
		UserID:      userID,
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *Service) ChooseWinner(
	ctx context.Context,
	giveawayID ulid.ULID,
) ([]entity.Winner, error) {
	participants, err := c.giveawaysParticipantsRepository.GetManyByGiveawayID(ctx, giveawayID)
	if err != nil {
		return nil, err
	}

	if len(participants) == 0 {
		return nil, nil
	}

	winners := []model.ChannelGiveawayParticipant{}
	err = c.txManager.Do(ctx, func(ctx context.Context) error {
		err = c.giveawaysParticipantsRepository.ResetWinners(
			ctx,
			lo.Map(
				participants,
				func(participant model.ChannelGiveawayParticipant, _ int) ulid.ULID {
					return participant.ID
				},
			)...,
		)
		if err != nil {
			return err
		}

		winnerInd := rand.Intn(len(participants))

		var winner model.ChannelGiveawayParticipant
		winner, err = c.giveawaysParticipantsRepository.Update(
			ctx,
			participants[winnerInd].ID,
			giveaways_participants.UpdateInput{
				IsWinner: lo.ToPtr(true),
			},
		)
		if err != nil {
			return err
		}
		winners = append(winners, winner)

		return nil
	})
	if err != nil {
		return nil, err
	}

	return lo.Map(winners, func(winner model.ChannelGiveawayParticipant, _ int) entity.Winner {
		return entity.Winner{
			UserID:      winner.UserID,
			DisplayName: winner.DisplayName,
		}
	}), nil
}
