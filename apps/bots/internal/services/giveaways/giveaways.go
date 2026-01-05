package giveaways

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	buscore "github.com/twirapp/twir/libs/bus-core"
	giveawaysbusmodel "github.com/twirapp/twir/libs/bus-core/giveaways"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	"github.com/twirapp/twir/libs/logger"
	"github.com/twirapp/twir/libs/repositories/giveaways"
	giveawaymodel "github.com/twirapp/twir/libs/repositories/giveaways/model"
	"github.com/twirapp/twir/libs/repositories/giveaways_participants"
	"github.com/twirapp/twir/libs/repositories/giveaways_participants/model"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In
	LC fx.Lifecycle

	TxManager                       trm.Manager
	GiveawaysRepository             giveaways.Repository
	GiveawaysParticipantsRepository giveaways_participants.Repository
	GiveawaysCacher                 *generic_cacher.GenericCacher[[]giveawaymodel.ChannelGiveaway]
	Logger                          *slog.Logger
	Redis                           *redis.Client
	TwirBus                         *buscore.Bus
}

func New(opts Opts) *Service {
	s := &Service{
		txManager:                       opts.TxManager,
		giveawaysRepository:             opts.GiveawaysRepository,
		giveawaysParticipantsRepository: opts.GiveawaysParticipantsRepository,
		giveawaysCacher:                 opts.GiveawaysCacher,
		logger:                          opts.Logger,
		redis:                           opts.Redis,
		twirBus:                         opts.TwirBus,
	}

	opts.LC.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				return s.twirBus.Giveaways.ChooseWinner.SubscribeGroup(
					"giveaways",
					s.chooseWinner,
				)
			},
			OnStop: func(ctx context.Context) error {
				s.twirBus.Giveaways.ChooseWinner.Unsubscribe()

				return nil
			},
		},
	)

	return s
}

type Service struct {
	txManager                       trm.Manager
	giveawaysRepository             giveaways.Repository
	giveawaysParticipantsRepository giveaways_participants.Repository
	giveawaysCacher                 *generic_cacher.GenericCacher[[]giveawaymodel.ChannelGiveaway]
	logger                          *slog.Logger
	redis                           *redis.Client
	twirBus                         *buscore.Bus
}

const redisParticipantKey = "giveaways:%s:participants:%s"

func (c *Service) TryAddParticipant(
	ctx context.Context,
	userID string,
	userLogin string,
	userDisplayName string,
	giveawayID string,
) error {
	cacheKey := fmt.Sprintf(redisParticipantKey, giveawayID, userID)
	exists, _ := c.redis.Exists(ctx, cacheKey).Result()
	if exists >= 1 {
		return nil
	}

	_, err := c.giveawaysParticipantsRepository.Create(
		ctx,
		giveaways_participants.CreateInput{
			GiveawayID:      giveawayID,
			UserID:          userID,
			UserLogin:       userLogin,
			UserDisplayName: userDisplayName,
		},
	)
	if err != nil {
		return err
	}

	if err := c.twirBus.Giveaways.NewParticipants.Publish(
		ctx,
		giveawaysbusmodel.NewParticipant{
			GiveawayID:      giveawayID,
			UserID:          userID,
			UserLogin:       userLogin,
			UserDisplayName: userDisplayName,
		},
	); err != nil {
		c.logger.Error("cannot publish new participant", logger.Error(err))
	}

	if err := c.redis.Set(ctx, cacheKey, 1, 1*time.Hour).Err(); err != nil {
		c.logger.Error("cannot set giveaway participant to redis cache", logger.Error(err))
	}

	return nil
}

func (c *Service) chooseWinner(
	ctx context.Context,
	req giveawaysbusmodel.ChooseWinnerRequest,
) (giveawaysbusmodel.ChooseWinnerResponse, error) {
	parsedGiveawayId, err := uuid.Parse(req.GiveawayID)
	if err != nil {
		return giveawaysbusmodel.ChooseWinnerResponse{}, err
	}

	giveaway, err := c.giveawaysRepository.GetByID(ctx, parsedGiveawayId)
	if err != nil {
		return giveawaysbusmodel.ChooseWinnerResponse{}, err
	}
	if giveaway == giveawaymodel.ChannelGiveawayNil {
		return giveawaysbusmodel.ChooseWinnerResponse{}, fmt.Errorf("giveaway not found")
	}

	participants, err := c.giveawaysParticipantsRepository.GetManyByGiveawayID(
		ctx,
		req.GiveawayID,
		giveaways_participants.GetManyInput{
			IgnoreWinners: true,
		},
	)
	if err != nil {
		c.logger.Error("cannot get participants", logger.Error(err))
		return giveawaysbusmodel.ChooseWinnerResponse{}, err
	}

	if len(participants) == 0 {
		return giveawaysbusmodel.ChooseWinnerResponse{}, fmt.Errorf("Cannot do roll with empty participants")
	}

	winners := make([]model.ChannelGiveawayParticipant, 0, 1)
	err = c.txManager.Do(
		ctx,
		func(txCtx context.Context) error {
			// TODO: what if we wanna choose multiple winners?
			// err = c.giveawaysParticipantsRepository.ResetWinners(
			// 	txCtx,
			// 	lo.Map(
			// 		participants,
			// 		func(participant model.ChannelGiveawayParticipant, _ int) string {
			// 			return participant.ID.String()
			// 		},
			// 	)...,
			// )
			// if err != nil {
			// 	c.logger.Error("reset winners error", logger.Error(err))
			// 	return err
			// }

			winnerInd := rand.Intn(len(participants))

			var winner model.ChannelGiveawayParticipant
			winner, err = c.giveawaysParticipantsRepository.Update(
				txCtx,
				participants[winnerInd].ID.String(),
				giveaways_participants.UpdateInput{
					IsWinner: lo.ToPtr(true),
				},
			)
			if err != nil {
				c.logger.Error("update winner error", logger.Error(err))
				return err
			}
			winners = append(winners, winner)

			_, err = c.giveawaysRepository.UpdateStatuses(
				txCtx,
				winner.GiveawayID,
				giveaways.UpdateStatusInput{
					StoppedAt: null.NewTime(time.Now(), true),
				},
			)
			if err != nil {
				c.logger.Error("update error", logger.Error(err))
				return err
			}

			return nil
		},
	)
	if err != nil {
		c.logger.Error("tx error", logger.Error(err))
		return giveawaysbusmodel.ChooseWinnerResponse{}, err
	}

	if err := c.giveawaysCacher.Invalidate(ctx, giveaway.ChannelID); err != nil {
		c.logger.Error("cannot invalidate giveaways cache", logger.Error(err))
	}

	mappedWinners := make([]giveawaysbusmodel.Winner, 0, len(winners))
	for _, winner := range winners {
		mappedWinners = append(
			mappedWinners,
			giveawaysbusmodel.Winner{
				UserID:          winner.UserID,
				UserLogin:       winner.UserLogin,
				UserDisplayName: winner.DisplayName,
			},
		)
	}

	return giveawaysbusmodel.ChooseWinnerResponse{
		Winners: mappedWinners,
	}, nil
}
