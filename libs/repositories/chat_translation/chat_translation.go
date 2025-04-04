package chat_translation

import (
	"context"
	"fmt"

	"github.com/oklog/ulid/v2"
	"github.com/twirapp/twir/libs/repositories/chat_translation/model"
)

var ErrSettingsNotFound = fmt.Errorf("channel settings not found")

type Repository interface {
	GetByChannelID(ctx context.Context, channelID string) (model.ChatTranslation, error)
	Create(ctx context.Context, input CreateInput) (model.ChatTranslation, error)
	Update(ctx context.Context, id ulid.ULID, input UpdateInput) (model.ChatTranslation, error)
	Delete(ctx context.Context, id ulid.ULID) error
}

type CreateInput struct {
	ChannelID         string
	Enabled           bool
	TargetLanguage    string
	ExcludedLanguages []string
	UseItalic         bool
	ExcludedUsersIDs  []string
}

type UpdateInput struct {
	Enabled           *bool
	TargetLanguage    *string
	ExcludedLanguages *[]string
	UseItalic         *bool
	ExcludedUsersIDs  *[]string
}
