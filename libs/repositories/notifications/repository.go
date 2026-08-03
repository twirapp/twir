package notifications

import (
	"context"
	"time"

	"github.com/twirapp/twir/libs/entities/notification"
)

type UpsertDiscordInput struct {
	DiscordMessageID      string
	DiscordChannelID      string
	Text                  *string
	EditorJSJSON          *string
	CreatedAt             time.Time
	UpdatedAt             *time.Time
	DiscordAttachmentKeys []string
}

type UpsertDiscordResult struct {
	Notification           notification.Notification
	PreviousAttachmentKeys []string
	Created                bool
}

type Repository interface {
	UpsertDiscord(ctx context.Context, input UpsertDiscordInput) (UpsertDiscordResult, error)
	DeleteDiscord(ctx context.Context, channelID string, messageIDs []string) ([]string, error)
}
