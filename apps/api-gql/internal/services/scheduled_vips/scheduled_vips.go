package scheduled_vips

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	scheduledvipsentity "github.com/twirapp/twir/libs/entities/scheduled_vips"
	scheduledvipsrepository "github.com/twirapp/twir/libs/repositories/scheduled_vips"
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

func (c *Service) GetScheduledVips(ctx context.Context, channelID string) (
	[]scheduledvipsentity.ScheduledVip,
	error,
) {
	scheduledVips, err := c.repo.GetMany(
		ctx,
		scheduledvipsrepository.GetManyInput{
			ChannelID: &channelID,
		},
	)
	if err != nil {
		return nil, err
	}

	return scheduledVips, nil
}

type CreateInput struct {
	UserID     string
	ChannelID  string
	RemoveAt   *time.Time
	RemoveType *scheduledvipsentity.RemoveType
}

func (c *Service) Create(ctx context.Context, input CreateInput) error {
	return c.repo.Create(
		ctx,
		scheduledvipsrepository.CreateInput{
			ChannelID:  input.ChannelID,
			UserID:     input.UserID,
			RemoveAt:   input.RemoveAt,
			RemoveType: input.RemoveType,
		},
	)
}

type RemoveInput struct {
	ID        string
	ChannelID string
	KeepVip   *bool
}

func (c *Service) Remove(ctx context.Context, input RemoveInput) error {
	id, err := uuid.Parse(input.ID)
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
	vipID, err := uuid.Parse(id)
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
