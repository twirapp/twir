package chat_wall

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	chatwallrepository "github.com/twirapp/twir/libs/repositories/chat_wall"
	chatwallmodel "github.com/twirapp/twir/libs/repositories/chat_wall/model"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	Repository chatwallrepository.Repository
}

func New(opts Opts) *Service {
	return &Service{
		Repository: opts.Repository,
	}
}

type Service struct {
	Repository chatwallrepository.Repository
}

func (c *Service) mapModelToEntity(m chatwallmodel.ChatWall) entity.ChatWall {
	return entity.ChatWall{
		ID:                     m.ID.String(),
		ChannelID:              m.ChannelID.String(),
		CreatedAt:              m.CreatedAt,
		UpdatedAt:              m.UpdatedAt,
		Phrase:                 m.Phrase,
		Enabled:                m.Enabled,
		Action:                 entity.ChatWallAction(m.Action),
		DurationSeconds:        m.DurationSeconds,
		TimeoutDurationSeconds: m.TimeoutDurationSeconds,
		AffectedMessages:       m.AffectedMessages,
	}
}

func (c *Service) mapModelToEntityLog(m chatwallmodel.ChatWallLog) entity.ChatWallLog {
	return entity.ChatWallLog{
		ID:        m.ID.String(),
		CreatedAt: m.CreatedAt,
		Text:      m.Text,
		UserID:    m.UserID.String(),
	}
}

func (c *Service) GetChatWalls(ctx context.Context, channelID string) (
	[]entity.ChatWall,
	error,
) {
	parsedChannelID, err := uuid.Parse(channelID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse channel id: %w", err)
	}

	walls, err := c.Repository.GetMany(
		ctx,
		chatwallrepository.GetManyInput{
			ChannelID: parsedChannelID,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get chat walls: %w", err)
	}

	converted := make([]entity.ChatWall, len(walls))
	for idx, wall := range walls {
		converted[idx] = c.mapModelToEntity(wall)
	}

	return converted, nil
}

func (c *Service) GetLogs(ctx context.Context, channelId, wallID string) (
	[]entity.ChatWallLog,
	error,
) {
	parsedId, err := uuid.Parse(wallID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse wall id: %w", err)
	}

	wall, err := c.Repository.GetByID(ctx, parsedId)
	if err != nil {
		return nil, fmt.Errorf("failed to get chat wall: %w", err)
	}

	if wall.ChannelID.String() != channelId {
		return nil, fmt.Errorf("wall does not belong to channel")
	}

	logs, err := c.Repository.GetLogs(ctx, parsedId)
	if err != nil {
		return nil, fmt.Errorf("failed to get chat wall logs: %w", err)
	}

	converted := make([]entity.ChatWallLog, len(logs))
	for idx, log := range logs {
		converted[idx] = c.mapModelToEntityLog(log)
	}

	return converted, nil
}

var ErrSettingsNotFound = fmt.Errorf("channel settings not found")

func (c *Service) GetChannelSettings(ctx context.Context, channelID string) (
	entity.ChatWallSettings,
	error,
) {
	parsedChannelID, err := uuid.Parse(channelID)
	if err != nil {
		return entity.ChatWallSettings{}, fmt.Errorf("failed to parse channel id: %w", err)
	}

	settings, err := c.Repository.GetChannelSettings(ctx, parsedChannelID)
	if err != nil {
		if errors.Is(err, chatwallrepository.ErrSettingsNotFound) {
			return entity.ChatWallSettings{}, ErrSettingsNotFound
		}

		return entity.ChatWallSettings{}, fmt.Errorf("failed to get chat wall settings: %w", err)
	}

	return entity.ChatWallSettings{
		MuteSubscribers: settings.MuteSubscribers,
		MuteVips:        settings.MuteVips,
	}, nil
}

func (c *Service) UpdateChannelSettings(
	ctx context.Context,
	channelID string,
	muteSubscribers bool,
	muteVips bool,
) error {
	parsedChannelID, err := uuid.Parse(channelID)
	if err != nil {
		return fmt.Errorf("failed to parse channel id: %w", err)
	}

	return c.Repository.UpdateChannelSettings(
		ctx,
		chatwallrepository.UpdateChannelSettingsInput{
			ChannelID:       parsedChannelID,
			MuteSubscribers: muteSubscribers,
			MuteVips:        muteVips,
		},
	)
}
