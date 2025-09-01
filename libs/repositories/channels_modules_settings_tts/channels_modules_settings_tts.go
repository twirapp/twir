package channels_modules_settings_tts

import (
	"context"

	"github.com/twirapp/twir/libs/repositories/channels_modules_settings_tts/model"
)

type Repository interface {
	// Channel-level settings (userId will be null)

	GetByChannelID(ctx context.Context, channelID string) (model.ChannelModulesSettingsTTS, error)
	CreateForChannel(ctx context.Context, input CreateOrUpdateInput) (
		model.ChannelModulesSettingsTTS,
		error,
	)
	UpdateForChannel(
		ctx context.Context,
		channelID string,
		input CreateOrUpdateInput,
	) (model.ChannelModulesSettingsTTS, error)
	DeleteForChannel(ctx context.Context, channelID string) error

	// User-specific settings (userId will be set)

	GetByChannelIDAndUserID(
		ctx context.Context,
		channelID, userID string,
	) (model.ChannelModulesSettingsTTS, error)
	CreateForUser(ctx context.Context, input CreateOrUpdateInput) (
		model.ChannelModulesSettingsTTS,
		error,
	)
	UpdateForUser(
		ctx context.Context,
		channelID, userID string,
		input CreateOrUpdateInput,
	) (model.ChannelModulesSettingsTTS, error)
	DeleteForUser(ctx context.Context, channelID, userID string) error

	// Get all settings for a channel (both channel-level and user-specific)
	GetAllByChannelID(ctx context.Context, channelID string) (
		[]model.ChannelModulesSettingsTTS,
		error,
	)
}

type CreateOrUpdateInput struct {
	ChannelID                          string
	UserID                             *string // null for channel-level settings, set for user-specific
	Enabled                            *bool
	Rate                               int
	Volume                             int
	Pitch                              int
	Voice                              string
	AllowUsersChooseVoiceInMainCommand bool
	MaxSymbols                         int
	DisallowedVoices                   []string
	DoNotReadEmoji                     bool
	DoNotReadTwitchEmotes              bool
	DoNotReadLinks                     bool
	ReadChatMessages                   bool
	ReadChatMessagesNicknames          bool
}
