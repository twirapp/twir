package chat_wall

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/repositories/chat_wall/model"
)

var ErrSettingsNotFound = fmt.Errorf("channel settings not found")

type Repository interface {
	GetChannelSettings(ctx context.Context, channelID uuid.UUID) (model.ChatWallSettings, error)
	UpdateChannelSettings(ctx context.Context, input UpdateChannelSettingsInput) error
	GetByID(ctx context.Context, id uuid.UUID) (model.ChatWall, error)
	GetMany(ctx context.Context, input GetManyInput) ([]model.ChatWall, error)
	GetLogs(ctx context.Context, wallID uuid.UUID) ([]model.ChatWallLog, error)
	Create(ctx context.Context, input CreateInput) (model.ChatWall, error)
	CreateLog(ctx context.Context, input CreateLogInput) error
	CreateManyLogs(ctx context.Context, inputs []CreateLogInput) error
	Update(ctx context.Context, id uuid.UUID, input UpdateInput) (model.ChatWall, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type GetManyInput struct {
	ChannelID uuid.UUID
	Enabled   *bool
}

type CreateInput struct {
	ChannelID       uuid.UUID
	Phrase          string
	Enabled         bool
	Action          model.ChatWallAction
	Duration        time.Duration
	TimeoutDuration *time.Duration
}

type UpdateInput struct {
	Phrase          *string
	Enabled         *bool
	Action          *model.ChatWallAction
	Duration        *time.Duration
	TimeoutDuration *time.Duration
}

type CreateLogInput struct {
	WallID uuid.UUID
	UserID uuid.UUID
	Text   string
}

type UpdateChannelSettingsInput struct {
	ChannelID       uuid.UUID
	MuteSubscribers bool
	MuteVips        bool
}
