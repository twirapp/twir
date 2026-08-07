package chat_messages

import (
	"context"
	"time"

	"github.com/twirapp/twir/libs/repositories/chat_messages/model"
)

type Repository interface {
	Create(ctx context.Context, input CreateInput) error
	CreateMany(ctx context.Context, input []CreateInput) error
	GetMany(ctx context.Context, input GetManyInput) ([]model.ChatMessage, error)
	GetLatestByUser(
		ctx context.Context,
		input GetLatestByUserInput,
	) (model.ChatMessage, error)
}

type PlatformChannelIdentity struct {
	Platform          string
	PlatformChannelID string
}

type CreateInput struct {
	ID                string
	Platform          string
	PlatformChannelID string
	UserID            string
	Text              string
	UserName          string
	UserDisplayName   string
	UserColor         string
}

type GetManyInput struct {
	Page    int
	PerPage int

	ChannelPairs      []PlatformChannelIdentity
	Platform          *string
	PlatformChannelID *string
	UserNameLike      *string
	TextLike          *string
	TextFuzzy         *TextFuzzyFilter
	UserIDs           []string

	TimeGte *time.Time
}

// TextFuzzyFilter matches message texts against a phrase using the chat wall
// fuzzy semantics: exact case-insensitive substring, or a token whose
// Levenshtein distance (whole token or token prefix of phrase length) to the
// phrase is within MaxDistance. Semantics must stay in sync with the live
// matcher in apps/bots/internal/chatwallmatcher.
type TextFuzzyFilter struct {
	Phrase      string
	Length      int
	MaxDistance int
}

type GetLatestByUserInput struct {
	Platform          string
	PlatformChannelID string
	UserName          string
}
