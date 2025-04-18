package giveaways

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/guregu/null"
	"github.com/oklog/ulid/v2"
	"github.com/samber/lo"
	"github.com/satont/twir/libs/logger"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	buscore "github.com/twirapp/twir/libs/bus-core"
	giveawaysbus "github.com/twirapp/twir/libs/bus-core/giveaways"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	"github.com/twirapp/twir/libs/repositories/giveaways"
	"github.com/twirapp/twir/libs/repositories/giveaways/model"
	giveawaysmodel "github.com/twirapp/twir/libs/repositories/giveaways/model"
	"github.com/twirapp/twir/libs/repositories/giveaways_participants"
	giveawaysparticipantsmodel "github.com/twirapp/twir/libs/repositories/giveaways_participants/model"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	GiveawaysRepository             giveaways.Repository
	GiveawaysParticipantsRepository giveaways_participants.Repository
	GiveawaysCacher                 *generic_cacher.GenericCacher[[]model.ChannelGiveaway]
	TwirBus                         *buscore.Bus
	Logger                          logger.Logger
}

func New(opts Opts) *Service {
	return &Service{
		giveawaysRepository:             opts.GiveawaysRepository,
		giveawaysParticipantsRepository: opts.GiveawaysParticipantsRepository,
		giveawaysCacher:                 opts.GiveawaysCacher,
		twirBus:                         opts.TwirBus,
		logger:                          opts.Logger,
	}
}

type Service struct {
	giveawaysRepository             giveaways.Repository
	giveawaysParticipantsRepository giveaways_participants.Repository
	giveawaysCacher                 *generic_cacher.GenericCacher[[]model.ChannelGiveaway]
	twirBus                         *buscore.Bus
	logger                          logger.Logger
}

type CreateInput struct {
	ChannelID       string
	Keyword         string
	CreatedByUserID string
}

func (c *Service) Start(
	ctx context.Context,
	giveawayID ulid.ULID,
	channelID string,
) (entity.ChannelGiveaway, error) {
	dbGiveaway, err := c.GiveawayGet(ctx, giveawayID, channelID)
	if err != nil {
		return entity.ChannelGiveawayNil, err
	}

	if dbGiveaway == entity.ChannelGiveawayNil {
		return entity.ChannelGiveawayNil, fmt.Errorf("Giveaway doesnt exists")
	}

	if dbGiveaway.ArchivedAt != nil {
		return entity.ChannelGiveawayNil, fmt.Errorf("Cannot run archived giveaway")
	}

	if dbGiveaway.EndedAt != nil {
		return entity.ChannelGiveawayNil, fmt.Errorf("Cannot run ended giveaway")
	}

	if dbGiveaway.StoppedAt != nil {
		dbGiveaway, err = c.GiveawayUpdateStatus(
			ctx,
			giveawayID,
			channelID,
			UpdateStatusInput{
				StartedAt: null.NewTime(time.Now(), true),
				StoppedAt: null.Time{},
			},
		)
		if err != nil {
			return entity.ChannelGiveawayNil, err
		}

		return dbGiveaway, nil
	}

	if dbGiveaway.StartedAt != nil {
		return entity.ChannelGiveawayNil, fmt.Errorf("Giveaway already started and not stopped")
	}

	if dbGiveaway.StartedAt == nil {
		dbGiveaway, err = c.GiveawayUpdateStatus(
			ctx,
			giveawayID,
			channelID,
			UpdateStatusInput{
				StartedAt: null.NewTime(time.Now(), true),
			},
		)
		if err != nil {
			return entity.ChannelGiveawayNil, err
		}

		return dbGiveaway, nil
	}

	return dbGiveaway, nil
}

func (c *Service) Stop(
	ctx context.Context,
	giveawayID ulid.ULID,
	channelID string,
) (entity.ChannelGiveaway, error) {
	dbGiveaway, err := c.GiveawayGet(ctx, giveawayID, channelID)
	if err != nil {
		return entity.ChannelGiveawayNil, err
	}

	if dbGiveaway == entity.ChannelGiveawayNil {
		return entity.ChannelGiveawayNil, fmt.Errorf("Giveaway doesnt exists")
	}

	if dbGiveaway.ArchivedAt != nil {
		return entity.ChannelGiveawayNil, fmt.Errorf("Cannot stop archived giveaway")
	}

	if dbGiveaway.EndedAt != nil {
		return entity.ChannelGiveawayNil, fmt.Errorf("Cannot stop ended giveaway")
	}

	if dbGiveaway.StartedAt == nil {
		return entity.ChannelGiveawayNil, fmt.Errorf("Cannot stop not started giveaway")
	}

	if dbGiveaway.StoppedAt != nil {
		return entity.ChannelGiveawayNil, fmt.Errorf("Giveaway already stopped")
	}

	dbGiveaway, err = c.GiveawayUpdateStatus(
		ctx,
		giveawayID,
		channelID,
		UpdateStatusInput{
			StoppedAt: null.NewTime(time.Now(), true),
		},
	)
	if err != nil {
		return entity.ChannelGiveawayNil, err
	}

	return dbGiveaway, nil
}

func (c *Service) ChooseWinners(
	ctx context.Context,
	giveawayID ulid.ULID,
	channelID string,
) ([]entity.ChannelGiveawayWinner, error) {
	dbGiveaway, err := c.GiveawayGet(ctx, giveawayID, channelID)
	if err != nil {
		return nil, err
	}

	if dbGiveaway == entity.ChannelGiveawayNil {
		return nil, fmt.Errorf("Giveaway doesnt exists")
	}

	if dbGiveaway.ArchivedAt != nil {
		return nil, fmt.Errorf("Cannot choose winners for archived giveaways")
	}

	winners, err := c.twirBus.Giveaways.ChooseWinner.Request(ctx, giveawaysbus.ChooseWinnerRequest{
		GiveawayID: giveawayID.String(),
	})
	if err != nil {
		c.logger.Error("Cannot choose winners", slog.Any("err", err))
		return nil, err
	}

	if len(winners.Data.Winners) == 0 {
		return nil, fmt.Errorf("Cannot process winners for now, please retry.")
	}

	var result []entity.ChannelGiveawayWinner
	for _, winner := range winners.Data.Winners {
		result = append(result, c.giveawayWinnerBusModelToEntity(winner))
	}

	err = c.updateGiveawaysCacheForChannel(ctx, channelID)
	if err != nil {
		c.logger.Error("Cannot update winners cache", slog.Any("err", err))
		return nil, err
	}

	return result, nil
}

func (c *Service) Create(ctx context.Context, input CreateInput) (entity.ChannelGiveaway, error) {
	dbGiveaway, err := c.giveawaysRepository.GetByChannelIDAndKeyword(
		ctx,
		input.ChannelID,
		input.Keyword,
	)
	if err != nil {
		if !errors.Is(err, giveaways.ErrNotFound) {
			return entity.ChannelGiveawayNil, err
		}
	}

	if dbGiveaway != model.ChannelGiveawayNil {
		return entity.ChannelGiveawayNil, fmt.Errorf(
			"Giveaways with same keyword already exists on this channel",
		)
	}

	giveaway, err := c.giveawaysRepository.Create(
		ctx,
		giveaways.CreateInput{
			ChannelID:       input.ChannelID,
			Keyword:         input.Keyword,
			CreatedByUserID: input.CreatedByUserID,
		},
	)
	if err != nil {
		return entity.ChannelGiveawayNil, err
	}

	err = c.updateGiveawaysCacheForChannel(ctx, input.ChannelID)
	if err != nil {
		return entity.ChannelGiveawayNil, err
	}

	return c.giveawayModelToEntity(giveaway), nil
}

func (c *Service) GetParticipantsForGiveaway(
	ctx context.Context,
	giveawayID ulid.ULID,
) ([]entity.ChannelGiveawayParticipant, error) {
	participants, err := c.giveawaysParticipantsRepository.GetManyByGiveawayID(
		ctx,
		giveawayID.String(),
	)
	if err != nil {
		return nil, err
	}

	return lo.Map(
		participants,
		func(item giveawaysparticipantsmodel.ChannelGiveawayParticipant, _ int) entity.ChannelGiveawayParticipant {
			return c.giveawayParticipantModelToEntity(item)
		},
	), nil
}

func (c *Service) GiveawayGet(
	ctx context.Context,
	giveawayID ulid.ULID,
	channelID string,
) (entity.ChannelGiveaway, error) {
	giveaway, err := c.giveawaysRepository.GetByID(ctx, giveawayID)
	if err != nil {
		if !errors.Is(err, giveaways.ErrNotFound) {
			return entity.ChannelGiveawayNil, err
		}
	}

	if giveaway == giveawaysmodel.ChannelGiveawayNil {
		return entity.ChannelGiveawayNil, nil
	}

	err = c.updateGiveawaysCacheForChannel(ctx, channelID)
	if err != nil {
		return entity.ChannelGiveawayNil, err
	}

	return c.giveawayModelToEntity(giveaway), nil
}

func (c *Service) GiveawaysGetMany(
	ctx context.Context,
	channelID string,
) ([]entity.ChannelGiveaway, error) {
	dbGiveaways, err := c.giveawaysRepository.GetManyByChannelID(ctx, channelID)
	if err != nil {
		return nil, err
	}

	return lo.Map(
		dbGiveaways,
		func(item giveawaysmodel.ChannelGiveaway, _ int) entity.ChannelGiveaway {
			return c.giveawayModelToEntity(item)
		},
	), nil
}

func (c *Service) ArchiveGiveaway(
	ctx context.Context,
	giveawayID ulid.ULID,
	channelID string,
) (entity.ChannelGiveaway, error) {
	giveaway, err := c.giveawaysRepository.UpdateStatuses(
		ctx,
		giveawayID,
		giveaways.UpdateStatusInput{
			ArchivedAt: null.NewTime(time.Now(), true),
		},
	)
	if err != nil {
		return entity.ChannelGiveawayNil, err
	}

	err = c.updateGiveawaysCacheForChannel(ctx, channelID)
	if err != nil {
		return entity.ChannelGiveawayNil, err
	}

	return c.giveawayModelToEntity(giveaway), nil
}

func (c *Service) GiveawayRemove(
	ctx context.Context,
	giveawayID ulid.ULID,
	channelID string,
) error {
	err := c.giveawaysRepository.Delete(ctx, giveawayID)
	if err != nil {
		return nil
	}

	err = c.updateGiveawaysCacheForChannel(ctx, channelID)

	return err
}

type UpdateInput struct {
	StartedAt  *time.Time
	EndedAt    *time.Time
	Keyword    *string
	ArchivedAt *time.Time
	StoppedAt  *time.Time
}

func (c *Service) GiveawayUpdate(
	ctx context.Context,
	giveawayID ulid.ULID,
	channelID string,
	input UpdateInput,
) (entity.ChannelGiveaway, error) {
	dbGiveaway, err := c.giveawaysRepository.Update(ctx, giveawayID, giveaways.UpdateInput{
		StartedAt:  input.StartedAt,
		EndedAt:    input.EndedAt,
		Keyword:    input.Keyword,
		ArchivedAt: input.ArchivedAt,
		StoppedAt:  input.StoppedAt,
	})
	if err != nil {
		return entity.ChannelGiveawayNil, err
	}

	err = c.updateGiveawaysCacheForChannel(ctx, channelID)
	if err != nil {
		return entity.ChannelGiveawayNil, err
	}

	return c.giveawayModelToEntity(dbGiveaway), nil
}

type UpdateStatusInput struct {
	StartedAt  null.Time
	EndedAt    null.Time
	ArchivedAt null.Time
	StoppedAt  null.Time
}

func (c *Service) GiveawayUpdateStatus(
	ctx context.Context,
	giveawayID ulid.ULID,
	channelID string,
	input UpdateStatusInput,
) (entity.ChannelGiveaway, error) {
	dbGiveaway, err := c.giveawaysRepository.UpdateStatuses(
		ctx,
		giveawayID,
		giveaways.UpdateStatusInput{
			StartedAt:  input.StartedAt,
			EndedAt:    input.EndedAt,
			ArchivedAt: input.ArchivedAt,
			StoppedAt:  input.StoppedAt,
		},
	)
	if err != nil {
		return entity.ChannelGiveawayNil, err
	}

	err = c.updateGiveawaysCacheForChannel(ctx, channelID)
	if err != nil {
		return entity.ChannelGiveawayNil, err
	}

	return c.giveawayModelToEntity(dbGiveaway), nil
}

func (c *Service) giveawayModelToEntity(m giveawaysmodel.ChannelGiveaway) entity.ChannelGiveaway {
	return entity.ChannelGiveaway{
		ID:              m.ID,
		ChannelID:       m.ChannelID,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
		StartedAt:       m.StartedAt,
		EndedAt:         m.EndedAt,
		Keyword:         m.Keyword,
		CreatedByUserID: m.CreatedByUserID,
		ArchivedAt:      m.ArchivedAt,
		StoppedAt:       m.StoppedAt,
	}
}

func (c *Service) giveawayParticipantModelToEntity(
	m giveawaysparticipantsmodel.ChannelGiveawayParticipant,
) entity.ChannelGiveawayParticipant {
	return entity.ChannelGiveawayParticipant{
		DisplayName: m.DisplayName,
		UserID:      m.UserID,
		IsWinner:    m.IsWinner,
		ID:          m.ID,
		GiveawayID:  m.GiveawayID,
	}
}

func (c *Service) giveawayWinnerBusModelToEntity(
	m giveawaysbus.Winner,
) entity.ChannelGiveawayWinner {
	return entity.ChannelGiveawayWinner{
		DisplayName: m.DisplayName,
		UserID:      m.UserID,
	}
}

/*
We need to update value of array of channels giveaways in Redis,
cus we search for keyword for every message and don't wanna use database calls in message handlers,
so we are do probably some unnecessarily work here but provide better consistency.
Also, limits for max giveaways per channel is low, so it will be fast, I suppose.
*/
func (c *Service) updateGiveawaysCacheForChannel(ctx context.Context, channelID string) error {
	dbGiveaways, err := c.giveawaysRepository.GetManyActiveByChannelID(ctx, channelID)
	if err != nil {
		return err
	}

	err = c.giveawaysCacher.SetValue(ctx, channelID, dbGiveaways)
	if err != nil {
		return err
	}

	return nil
}
