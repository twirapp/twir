package services

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/guregu/null"
	"github.com/samber/lo"
	"github.com/satont/twir/libs/logger"
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
	Logger                          logger.Logger
}

func New(opts Opts) *Service {
	return &Service{
		txManager:                       opts.TxManager,
		giveawaysRepository:             opts.GiveawaysRepository,
		giveawaysParticipantsRepository: opts.GiveawaysParticipantsRepository,
		logger:                          opts.Logger,
	}
}

type Service struct {
	txManager                       trm.Manager
	giveawaysRepository             giveaways.Repository
	giveawaysParticipantsRepository giveaways_participants.Repository
	logger                          logger.Logger
}

func (c *Service) TryAddParticipant(
	ctx context.Context,
	userID string,
	displayName string,
	giveawayID string,
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
	giveawayID string,
) ([]entity.Winner, error) {
	participants, err := c.giveawaysParticipantsRepository.GetManyByGiveawayID(ctx, giveawayID)
	if err != nil {
		c.logger.Error("cannot get participants", slog.Any("err", err))
		return nil, err
	}

	if len(participants) == 0 {
		return nil, fmt.Errorf("Cannot do roll with empty participants")
	}

	winners := []model.ChannelGiveawayParticipant{}
	err = c.txManager.Do(ctx, func(ctx context.Context) error {
		err = c.giveawaysParticipantsRepository.ResetWinners(
			ctx,
			lo.Map(
				participants,
				func(participant model.ChannelGiveawayParticipant, _ int) string {
					return participant.ID.String()
				},
			)...,
		)
		if err != nil {
			c.logger.Error("reset winners error", slog.Any("err", err))
			return err
		}

		winnerInd := rand.Intn(len(participants))

		var winner model.ChannelGiveawayParticipant
		winner, err = c.giveawaysParticipantsRepository.Update(
			ctx,
			string(participants[winnerInd].ID.String()),
			giveaways_participants.UpdateInput{
				IsWinner: lo.ToPtr(true),
			},
		)
		if err != nil {
			c.logger.Error("update winner error", slog.Any("err", err))
			return err
		}
		winners = append(winners, winner)

		_, err = c.giveawaysRepository.UpdateStatuses(
			ctx,
			winner.GiveawayID,
			giveaways.UpdateStatusInput{
				EndedAt: null.NewTime(time.Now(), true),
			},
		)
		if err != nil {
			c.logger.Error("update error", slog.Any("err", err))
			return err
		}

		return nil
	})
	if err != nil {
		c.logger.Error("tx error", slog.Any("err", err))
		return nil, err
	}

	return lo.Map(winners, func(winner model.ChannelGiveawayParticipant, _ int) entity.Winner {
		return entity.Winner{
			UserID:      winner.UserID,
			DisplayName: winner.DisplayName,
		}
	}), nil
}
