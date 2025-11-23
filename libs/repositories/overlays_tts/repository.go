package overlays_tts

import (
	"context"
	"fmt"

	"github.com/twirapp/twir/libs/repositories/overlays_tts/model"
)

type Repository interface {
	GetByChannelID(ctx context.Context, channelID string) (model.TTSOverlay, error)
	Create(ctx context.Context, input CreateInput) (model.TTSOverlay, error)
	Update(ctx context.Context, channelID string, input UpdateInput) (model.TTSOverlay, error)
	GetOrCreate(ctx context.Context, channelID string) (model.TTSOverlay, error)

	// User settings methods
	GetUserSettings(ctx context.Context, channelID, userID string) (model.TTSUserSettings, error)
	GetAllUserSettings(ctx context.Context, channelID string) ([]model.TTSUserSettings, error)
	CreateUserSettings(ctx context.Context, input CreateUserSettingsInput) (
		model.TTSUserSettings,
		error,
	)
	UpdateUserSettings(
		ctx context.Context,
		channelID, userID string,
		input UpdateUserSettingsInput,
	) (model.TTSUserSettings, error)
	GetOrCreateUserSettings(
		ctx context.Context,
		channelID, userID string,
		defaults CreateUserSettingsInput,
	) (model.TTSUserSettings, error)
	DeleteUserSettings(ctx context.Context, channelID, userID string) error
	DeleteMultipleUserSettings(ctx context.Context, channelID string, userIDs []string) error
}

type CreateInput struct {
	ChannelID string
	Settings  model.TTSOverlaySettings
}

type UpdateInput struct {
	Settings model.TTSOverlaySettings
}

type CreateUserSettingsInput struct {
	ChannelID string
	UserID    string
	Voice     string
	Rate      int32
	Pitch     int32
}

type UpdateUserSettingsInput struct {
	Voice *string
	Rate  *int32
	Pitch *int32
}

var ErrNotFound = fmt.Errorf("not found")
