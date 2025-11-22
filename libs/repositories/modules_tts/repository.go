package modules_tts

import (
	"context"
	"fmt"

	"github.com/twirapp/twir/libs/repositories/modules_tts/model"
)

type Repository interface {
	// Channel-level settings (userID will be null)
	GetByChannelID(ctx context.Context, channelID string) (model.TTS, error)
	CreateForChannel(ctx context.Context, input CreateInput) (model.TTS, error)
	UpdateForChannel(ctx context.Context, channelID string, input UpdateInput) (model.TTS, error)
	DeleteForChannel(ctx context.Context, channelID string) error

	// User-specific settings (userID will be set)
	GetByChannelIDAndUserID(ctx context.Context, channelID, userID string) (model.TTS, error)
	GetAllUsersByChannelID(ctx context.Context, channelID string) ([]model.TTS, error)
	CreateForUser(ctx context.Context, input CreateInput) (model.TTS, error)
	UpdateForUser(ctx context.Context, channelID, userID string, input UpdateInput) (model.TTS, error)
	DeleteForUser(ctx context.Context, channelID, userID string) error
	DeleteUsersForChannel(ctx context.Context, channelID string, userIDs []string) error
}

type CreateInput struct {
	ChannelID                          string
	UserID                             *string
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

type UpdateInput struct {
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

var ErrNotFound = fmt.Errorf("not found")
