package sended_messages_store

import (
	"context"

	discordsendednotifications "github.com/twirapp/twir/libs/repositories/discord_sended_notifications"
	"go.uber.org/fx"
)

type Message struct {
	GuildID          string `redis:"guildId"`
	MessageID        string `redis:"messageId"`
	TwitchChannelID  string `redis:"channelId"`
	DiscordChannelID string `redis:"discordChannelId"`
}

type SendedMessagesStore struct {
	repo discordsendednotifications.Repository
}

type Opts struct {
	fx.In

	Repo discordsendednotifications.Repository
}

func New(opts Opts) *SendedMessagesStore {
	return &SendedMessagesStore{
		repo: opts.Repo,
	}
}

func (c *SendedMessagesStore) Add(ctx context.Context, msg Message) error {
	return c.repo.Create(ctx, discordsendednotifications.CreateInput{
		GuildID:          msg.GuildID,
		MessageID:        msg.MessageID,
		TwitchChannelID:  msg.TwitchChannelID,
		DiscordChannelID: msg.DiscordChannelID,
	})
}

func (c *SendedMessagesStore) GetByMessageId(ctx context.Context, messageID string) (
	Message,
	error,
) {
	entity, err := c.repo.GetByMessageID(ctx, messageID)
	if err != nil {
		return Message{}, err
	}

	if entity.IsNil() {
		return Message{}, nil
	}

	return Message{
		GuildID:          entity.GuildID,
		MessageID:        entity.MessageID,
		TwitchChannelID:  entity.TwitchChannelID,
		DiscordChannelID: entity.DiscordChannelID,
	}, nil
}

func (c *SendedMessagesStore) GetByChannelId(ctx context.Context, channelId string) (
	[]Message,
	error,
) {
	entities, err := c.repo.GetByChannelID(ctx, channelId)
	if err != nil {
		return nil, err
	}

	result := make([]Message, 0, len(entities))
	for _, e := range entities {
		result = append(result, Message{
			GuildID:          e.GuildID,
			MessageID:        e.MessageID,
			TwitchChannelID:  e.TwitchChannelID,
			DiscordChannelID: e.DiscordChannelID,
		})
	}

	return result, nil
}

func (c *SendedMessagesStore) GetByGuildId(ctx context.Context, guildID string) (
	[]Message,
	error,
) {
	entities, err := c.repo.GetByGuildID(ctx, guildID)
	if err != nil {
		return nil, err
	}

	result := make([]Message, 0, len(entities))
	for _, e := range entities {
		result = append(result, Message{
			GuildID:          e.GuildID,
			MessageID:        e.MessageID,
			TwitchChannelID:  e.TwitchChannelID,
			DiscordChannelID: e.DiscordChannelID,
		})
	}

	return result, nil
}

func (c *SendedMessagesStore) DeleteByMessageId(ctx context.Context, messageID string) error {
	return c.repo.DeleteByMessageID(ctx, messageID)
}

func (c *SendedMessagesStore) DeleteByChannelId(ctx context.Context, channelID string) error {
	return c.repo.DeleteByChannelID(ctx, channelID)
}

func (c *SendedMessagesStore) DeleteByGuildId(ctx context.Context, guildID string) error {
	return c.repo.DeleteByGuildID(ctx, guildID)
}

func (c *SendedMessagesStore) GetAll(ctx context.Context) ([]Message, error) {
	entities, err := c.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	messages := make([]Message, 0, len(entities))
	for _, entity := range entities {
		messages = append(messages, Message{
			GuildID:          entity.GuildID,
			MessageID:        entity.MessageID,
			TwitchChannelID:  entity.TwitchChannelID,
			DiscordChannelID: entity.DiscordChannelID,
		})
	}

	return messages, nil
}
