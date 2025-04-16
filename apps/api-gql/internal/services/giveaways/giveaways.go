package giveaways

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
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
}

func New(opts Opts) *Service {
	return &Service{
		giveawaysRepository:             opts.GiveawaysRepository,
		giveawaysParticipantsRepository: opts.GiveawaysParticipantsRepository,
		giveawaysCacher:                 opts.GiveawaysCacher,
	}
}

type Service struct {
	giveawaysRepository             giveaways.Repository
	giveawaysParticipantsRepository giveaways_participants.Repository
	giveawaysCacher                 *generic_cacher.GenericCacher[[]model.ChannelGiveaway]
}

type CreateInput struct {
	ChannelID       string
	Keyword         string
	CreatedByUserID string
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

	// TODO: need to check for unqiue only for non archived giveaways only
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
	participants, err := c.giveawaysParticipantsRepository.GetManyByGiveawayID(ctx, giveawayID)
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

func (c *Service) GiveawayRemove(ctx context.Context, giveawayID ulid.ULID) error {
	return c.giveawaysRepository.Delete(ctx, giveawayID)
}

type UpdateInput struct {
	StartedAt  *time.Time
	EndedAt    *time.Time
	Keyword    *string
	ArchivedAt *time.Time
	StoppedAt  *time.Time
	IsArchived *bool
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
