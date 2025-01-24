package channels_commands_prefix

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	"github.com/twirapp/twir/libs/repositories/channels_commands_prefix"
	"github.com/twirapp/twir/libs/repositories/channels_commands_prefix/model"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	ChannelsCommandsPrefixRepository channels_commands_prefix.Repository
	Cacher                           *generic_cacher.GenericCacher[model.ChannelsCommandsPrefix]
}

func New(opts Opts) *Service {
	return &Service{
		channelsCommandsPrefixRepository: opts.ChannelsCommandsPrefixRepository,
		cacher:                           opts.Cacher,
	}
}

type Service struct {
	channelsCommandsPrefixRepository channels_commands_prefix.Repository
	cacher                           *generic_cacher.GenericCacher[model.ChannelsCommandsPrefix]
}

const DefaultPrefix = "!"

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
		if errors.Is(err, channels_commands_prefix.ErrNotFound) {
			return entity.ChannelsCommandsPrefixNil, nil
		}
		return entity.ChannelsCommandsPrefixNil, err
	}

	return c.modelToEntity(channelsCommandsPrefix), nil
}

func (c *Service) create(
	ctx context.Context,
	channelID, prefix string,
) (entity.ChannelsCommandsPrefix, error) {
	channelsCommandsPrefix, err := c.channelsCommandsPrefixRepository.Create(
		ctx,
		channels_commands_prefix.CreateInput{
			ChannelID: channelID,
			Prefix:    prefix,
		},
	)
	if err != nil {
		return entity.ChannelsCommandsPrefixNil, err
	}

	if err = c.cacher.Invalidate(ctx, channelID); err != nil {
		return entity.ChannelsCommandsPrefixNil, fmt.Errorf("cannot invalidate cache: %w", err)
	}

	return c.modelToEntity(channelsCommandsPrefix), nil
}

func (c *Service) update(
	ctx context.Context,
	id uuid.UUID,
	channelID, newPrefix string,
) (entity.ChannelsCommandsPrefix, error) {
	channelsCommandsPrefix, err := c.channelsCommandsPrefixRepository.Update(
		ctx,
		id,
		channels_commands_prefix.UpdateInput{
			Prefix: newPrefix,
		},
	)
	if err != nil {
		return entity.ChannelsCommandsPrefixNil, err
	}

	if err = c.cacher.Invalidate(ctx, channelID); err != nil {
		return entity.ChannelsCommandsPrefixNil, fmt.Errorf("cannot invalidate cache: %w", err)
	}

	return c.modelToEntity(channelsCommandsPrefix), nil
}

type UpdateInput struct {
	ChannelID string
	Prefix    string
}

func (c *Service) Update(
	ctx context.Context,
	input UpdateInput,
) (entity.ChannelsCommandsPrefix, error) {
	prefix, err := c.channelsCommandsPrefixRepository.GetByChannelID(ctx, input.ChannelID)
	if err != nil && !errors.Is(err, channels_commands_prefix.ErrNotFound) {
		return entity.ChannelsCommandsPrefixNil, err
	}

	if errors.Is(err, channels_commands_prefix.ErrNotFound) {
		return c.create(ctx, input.ChannelID, input.Prefix)
	} else {
		return c.update(ctx, prefix.ID, input.ChannelID, input.Prefix)
	}
}

func (c *Service) Delete(
	ctx context.Context,
	channelID string,
) error {
	prefix, err := c.channelsCommandsPrefixRepository.GetByChannelID(ctx, channelID)
	if err != nil {
		if errors.Is(err, channels_commands_prefix.ErrNotFound) {
			return nil
		}

		return err
	}

	if err = c.cacher.Invalidate(ctx, channelID); err != nil {
		return fmt.Errorf("cannot invalidate cache: %w", err)
	}

	return c.channelsCommandsPrefixRepository.Delete(ctx, prefix.ID)
}
