package discord_sended_notifications

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, input CreateInput) error
	GetByMessageID(ctx context.Context, messageID string) (DiscordSendedNotification, error)
	GetByChannelID(ctx context.Context, channelID string) ([]DiscordSendedNotification, error)
	GetByGuildID(ctx context.Context, guildID string) ([]DiscordSendedNotification, error)
	GetAll(ctx context.Context) ([]DiscordSendedNotification, error)
	DeleteByMessageID(ctx context.Context, messageID string) error
	DeleteByChannelID(ctx context.Context, channelID string) error
	DeleteByGuildID(ctx context.Context, guildID string) error
}

type CreateInput struct {
	GuildID          string
	MessageID        string
	TwitchChannelID  string
	DiscordChannelID string
}
