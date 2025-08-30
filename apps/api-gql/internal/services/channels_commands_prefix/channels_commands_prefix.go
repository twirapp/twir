package channels_commands_prefix

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	buscore "github.com/twirapp/twir/libs/bus-core"
	botssettings "github.com/twirapp/twir/libs/bus-core/bots-settings"
	"github.com/twirapp/twir/libs/cache"
	"github.com/twirapp/twir/libs/logger"
	"github.com/twirapp/twir/libs/repositories/channels_commands_prefix"
	"github.com/twirapp/twir/libs/repositories/channels_commands_prefix/model"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	Bus                              *buscore.Bus
	Logger                           logger.Logger
	ChannelsCommandsPrefixRepository channels_commands_prefix.Repository
	ChannelsCommandsPrefixCache      cache.Cache[model.ChannelsCommandsPrefix]
}

func New(opts Opts) *Service {
	return &Service{
		bus:                              opts.Bus,
		logger:                           opts.Logger,
		channelsCommandsPrefixRepository: opts.ChannelsCommandsPrefixRepository,
		channelsCommandsPrefixCache:      opts.ChannelsCommandsPrefixCache,
	}
}

type Service struct {
	bus                              *buscore.Bus
	logger                           logger.Logger
	channelsCommandsPrefixRepository channels_commands_prefix.Repository
	channelsCommandsPrefixCache      cache.Cache[model.ChannelsCommandsPrefix]
}

const DefaultPrefix = "!"

func (c *Service) GetByChannelID(
	ctx context.Context,
	channelID string,
) (entity.ChannelsCommandsPrefix, error) {
	channelsCommandsPrefix, err := c.channelsCommandsPrefixCache.Get(ctx, channelID)
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
	input UpdateInput,
) (entity.ChannelsCommandsPrefix, error) {
	prefix, err := c.channelsCommandsPrefixRepository.GetByChannelID(ctx, input.ChannelID)
	if err != nil && !errors.Is(err, channels_commands_prefix.ErrNotFound) {
		return entity.ChannelsCommandsPrefixNil, err
	}

	var prefixEntity entity.ChannelsCommandsPrefix

	if errors.Is(err, channels_commands_prefix.ErrNotFound) {
		if prefixEntity, err = c.create(ctx, input.ChannelID, input.Prefix); err != nil {
			return entity.ChannelsCommandsPrefixNil, fmt.Errorf("create: %w", err)
		}
	} else {
		if prefixEntity, err = c.update(ctx, prefix.ID, input.Prefix); err != nil {
			return entity.ChannelsCommandsPrefixNil, fmt.Errorf("update: %w", err)
		}
	}

	return prefixEntity, nil
}

func (c *Service) Reset(
	ctx context.Context,
	channelID string,
) error {
	input := UpdateInput{
		ChannelID: channelID,
		Prefix:    DefaultPrefix,
	}

	if _, err := c.Update(ctx, input); err != nil {
		if errors.Is(err, channels_commands_prefix.ErrNotFound) {
			return nil
		}

		return err
	}

	return nil
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

	go func() {
		if err = c.bus.BotsSettings.UpdatePrefix.Publish(
			ctx, botssettings.UpdatePrefixRequest{
				ID:        channelsCommandsPrefix.ID,
				ChannelID: channelsCommandsPrefix.ChannelID,
				Prefix:    channelsCommandsPrefix.Prefix,
				CreatedAt: channelsCommandsPrefix.CreatedAt,
				UpdatedAt: channelsCommandsPrefix.UpdatedAt,
			},
		); err != nil {
			c.logger.Error(
				"failed to publish channel command prefix update",
				slog.String("channel_id", channelID),
				slog.Any("error", err),
			)
		}
	}()

	return c.modelToEntity(channelsCommandsPrefix), nil
}

func (c *Service) update(
	ctx context.Context,
	id uuid.UUID,
	newPrefix string,
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

	go func() {
		if err = c.bus.BotsSettings.UpdatePrefix.Publish(
			ctx, botssettings.UpdatePrefixRequest{
				ID:        channelsCommandsPrefix.ID,
				ChannelID: channelsCommandsPrefix.ChannelID,
				Prefix:    channelsCommandsPrefix.Prefix,
				CreatedAt: channelsCommandsPrefix.CreatedAt,
				UpdatedAt: channelsCommandsPrefix.UpdatedAt,
			},
		); err != nil {
			c.logger.Error(
				"failed to publish channel command prefix update",
				slog.String("channel_id", channelsCommandsPrefix.ChannelID),
				slog.Any("error", err),
			)
		}
	}()

	return c.modelToEntity(channelsCommandsPrefix), nil
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
