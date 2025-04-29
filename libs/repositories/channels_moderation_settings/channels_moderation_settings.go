package channels_moderation_settings

import (
	"context"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/repositories/channels_moderation_settings/model"
)

type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (model.ChannelModerationSettings, error)
	GetByChannelID(ctx context.Context, channelID string) ([]model.ChannelModerationSettings, error)
	Create(ctx context.Context, input CreateOrUpdateInput) (model.ChannelModerationSettings, error)
	Update(ctx context.Context, id uuid.UUID, input CreateOrUpdateInput) (
		model.ChannelModerationSettings,
		error,
	)
	Delete(ctx context.Context, id uuid.UUID) error
}

type CreateOrUpdateInput struct {
	Name                        *string
	ChannelID                   string
	Type                        model.ModerationSettingsType
	Enabled                     bool
	BanTime                     int32
	BanMessage                  string
	WarningMessage              string
	CheckClips                  bool
	TriggerLength               int
	MaxPercentage               int
	DenyList                    []string
	DenyListRegexpEnabled       bool
	DenyListWordBoundaryEnabled bool
	DenyListSensitivityEnabled  bool
	DeniedChatLanguages         []string
	ExcludedRoles               []string
	MaxWarnings                 int
}
