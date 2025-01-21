package channels_commands_prefix

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	"github.com/twirapp/twir/libs/repositories/channels_commands_prefix"
	"github.com/twirapp/twir/libs/repositories/channels_commands_prefix/model"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	ChannelsCommandsPrefixRepository channels_commands_prefix.Repository
}

func New(opts Opts) *Service {
	return &Service{
		channelsCommandsPrefixRepository: opts.ChannelsCommandsPrefixRepository,
	}
}

type Service struct {
	channelsCommandsPrefixRepository channels_commands_prefix.Repository
}

func (c *Service) modelToEntity(m model.ChannelsCommandsPrefix) entity.ChannelsCommandsPrefix {
	return entity.ChannelsCommandsPrefix{
		ID:        m.ID,
		ChannelID: m.ChannelID,
		Prefix:    m.Prefix,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func (c *Service) GetByChannelID(
	ctx context.Context,
	channelID string,
) (entity.ChannelsCommandsPrefix, error) {
	channelsCommandsPrefix, err := c.channelsCommandsPrefixRepository.GetByChannelID(ctx, channelID)
	if err != nil {
		return entity.ChannelsCommandsPrefixNil, err
	}

	return c.modelToEntity(channelsCommandsPrefix), nil
}

type CreateInput struct {
	ChannelID string
	Prefix    string
}

func (c *Service) Create(
	ctx context.Context,
	input CreateInput,
) (entity.ChannelsCommandsPrefix, error) {
	channelsCommandsPrefix, err := c.channelsCommandsPrefixRepository.Create(
		ctx,
		channels_commands_prefix.CreateInput{
			ChannelID: input.ChannelID,
			Prefix:    input.Prefix,
		},
	)
	if err != nil {
		return entity.ChannelsCommandsPrefixNil, err
	}

	return c.modelToEntity(channelsCommandsPrefix), nil
}

type UpdateInput struct {
	ChannelID string
	Prefix    string
}

func (c *Service) Update(
	ctx context.Context,
	id uuid.UUID,
	input UpdateInput,
) (entity.ChannelsCommandsPrefix, error) {
	prefix, err := c.channelsCommandsPrefixRepository.GetByChannelID(ctx, input.ChannelID)
	if err != nil {
		return entity.ChannelsCommandsPrefixNil, err
	}

	if prefix.ChannelID != input.ChannelID {
		return entity.ChannelsCommandsPrefixNil, fmt.Errorf("channel ids not match")
	}

	channelsCommandsPrefix, err := c.channelsCommandsPrefixRepository.Update(
		ctx,
		id,
		channels_commands_prefix.UpdateInput{
			Prefix: input.Prefix,
		},
	)

	if err != nil {
		return entity.ChannelsCommandsPrefixNil, err
	}

	return c.modelToEntity(channelsCommandsPrefix), nil
}

func (c *Service) Delete(
	ctx context.Context,
	id uuid.UUID,
) error {
	return c.channelsCommandsPrefixRepository.Delete(ctx, id)
}
