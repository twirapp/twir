package giveaways

import (
	"context"

	"github.com/oklog/ulid/v2"
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

	newGiveaways, err := c.giveawaysRepository.GetManyByChannelID(ctx, input.ChannelID)
	if err != nil {
		return entity.ChannelGiveawayNil, err
	}

	err = c.giveawaysCacher.SetValue(ctx, input.ChannelID, newGiveaways)
	if err != nil {
		return entity.ChannelGiveawayNil, err
	}

	return c.giveawayModelToEntity(giveaway), nil
}

func (c *Service) GetParticipantsForGiveaway(
	ctx context.Context,
	giveawayID ulid.ULID,
) ([]entity.ChannelGiveawayParticipant, error) {
	// participants, err := c.giveawaysParticipantsRepository.Create()
	return nil, nil
}

func (c *Service) giveawayModelToEntity(m giveawaysmodel.ChannelGiveaway) entity.ChannelGiveaway {
	return entity.ChannelGiveaway{
		ID:              m.ID,
		ChannelID:       m.ChannelID,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
		StartedAt:       m.StartedAt,
		EndedAt:         m.EndedAt,
		IsRunning:       m.IsRunning,
		IsStopped:       m.IsStopped,
		IsFinished:      m.IsFinished,
		Keyword:         m.Keyword,
		CreatedByUserID: m.CreatedByUserID,
		ArchivedAt:      m.ArchivedAt,
		IsArchived:      m.IsArchived,
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
