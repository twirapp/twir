package scheduled_vips

import (
	"context"
	"fmt"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	scheduledvipsrepository "github.com/twirapp/twir/libs/repositories/scheduled_vips"
	scheduledvipmodel "github.com/twirapp/twir/libs/repositories/scheduled_vips/model"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	ScheduledVipsRepository scheduledvipsrepository.Repository
}

func New(opts Opts) *Service {
	return &Service{
		repo: opts.ScheduledVipsRepository,
	}
}

type Service struct {
	repo scheduledvipsrepository.Repository
}

func (c *Service) modelToEntity(m scheduledvipmodel.ScheduledVip) entity.ScheduledVip {
	return entity.ScheduledVip{
		ID:        m.ID,
		UserID:    m.UserID,
		ChannelID: m.ChannelID,
		CreatedAt: m.CreatedAt,
		RemoveAt:  m.RemoveAt,
	}
}

func (c *Service) GetScheduledVips(ctx context.Context, channelID string) (
	[]entity.ScheduledVip,
	error,
) {
	scheduledVips, err := c.repo.GetManyByChannelID(ctx, channelID)
	if err != nil {
		return nil, err
	}

	entities := make([]entity.ScheduledVip, 0, len(scheduledVips))
	for _, vip := range scheduledVips {
		entities = append(entities, c.modelToEntity(vip))
	}

	return entities, nil
}

type CreateInput struct {
	UserID    string
	ChannelID string
	RemoveAt  *time.Time
}

func (c *Service) Create(ctx context.Context, input CreateInput) error {
	return c.repo.Create(
		ctx,
		scheduledvipsrepository.CreateInput{
			ChannelID: input.ChannelID,
			UserID:    input.UserID,
			RemoveAt:  input.RemoveAt,
		},
	)
}

type RemoveInput struct {
	ID        string
	ChannelID string
	KeepVip   *bool
}

func (c *Service) Remove(ctx context.Context, input RemoveInput) error {
	id, err := ulid.Parse(input.ID)
	if err != nil {
		return err
	}

	vip, err := c.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if vip.ChannelID != input.ChannelID {
		return fmt.Errorf("vip does not belong to the channel")
	}

	return c.repo.Delete(ctx, id)
}

func (c *Service) Update(ctx context.Context, id, channelID string, removeAt *time.Time) error {
	vipID, err := ulid.Parse(id)
	if err != nil {
		return err
	}

	vip, err := c.repo.GetByID(ctx, vipID)
	if err != nil {
		return err
	}
	if vip.ChannelID != channelID {
		return fmt.Errorf("vip does not belong to the channel")
	}

	return c.repo.Update(
		ctx,
		vipID,
		scheduledvipsrepository.UpdateInput{
			RemoveAt: removeAt,
		},
	)
}
