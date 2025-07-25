package sended_messages_store

import (
	"context"
	"time"

	"github.com/google/uuid"
	model "github.com/twirapp/twir/libs/gomodels"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Message struct {
	GuildID          string `redis:"guildId"`
	MessageID        string `redis:"messageId"`
	TwitchChannelID  string `redis:"channelId"`
	DiscordChannelID string `redis:"discordChannelId"`
}

type SendedMessagesStore struct {
	db *gorm.DB
}

type Opts struct {
	fx.In

	DB *gorm.DB
}

func New(opts Opts) *SendedMessagesStore {
	return &SendedMessagesStore{
		db: opts.DB,
	}
}

func (c *SendedMessagesStore) Add(ctx context.Context, msg Message) error {
	entity := model.DiscordSendedNotification{
		ID:               uuid.NewString(),
		GuildID:          msg.GuildID,
		MessageID:        msg.MessageID,
		TwitchChannelID:  msg.TwitchChannelID,
		DiscordChannelID: msg.DiscordChannelID,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	return c.db.WithContext(ctx).Create(&entity).Error
}

func (c *SendedMessagesStore) GetByMessageId(ctx context.Context, messageID string) (
	Message,
	error,
) {
	msg := Message{}
	entity := model.DiscordSendedNotification{}
	err := c.db.WithContext(ctx).Where(
		`"message_id" = ?`,
		messageID,
	).First(&entity).Error
	if err != nil {
		return msg, err
	}

	msg.GuildID = entity.GuildID
	msg.MessageID = entity.MessageID
	msg.TwitchChannelID = entity.TwitchChannelID
	msg.DiscordChannelID = entity.DiscordChannelID

	return msg, nil
}

func (c *SendedMessagesStore) GetByChannelId(ctx context.Context, channelId string) (
	[]Message,
	error,
) {
	var result []Message
	var entities []model.DiscordSendedNotification
	err := c.db.WithContext(ctx).Where(
		`"channel_id" = ?`,
		channelId,
	).Find(&entities).Error
	if err != nil {
		return result, err
	}

	for _, e := range entities {
		result = append(
			result,
			Message{
				GuildID:          e.GuildID,
				MessageID:        e.MessageID,
				TwitchChannelID:  e.TwitchChannelID,
				DiscordChannelID: e.DiscordChannelID,
			},
		)
	}

	return result, nil
}

func (c *SendedMessagesStore) GetByGuildId(ctx context.Context, messageID string) (
	[]Message,
	error,
) {
	var result []Message
	var entities []model.DiscordSendedNotification
	err := c.db.WithContext(ctx).Where(
		`"guild_id" = ?`,
		messageID,
	).Find(&entities).Error
	if err != nil {
		return result, err
	}

	for _, e := range entities {
		result = append(
			result, Message{
				GuildID:          e.GuildID,
				MessageID:        e.MessageID,
				TwitchChannelID:  e.TwitchChannelID,
				DiscordChannelID: e.DiscordChannelID,
			},
		)
	}

	return result, nil
}

func (c *SendedMessagesStore) DeleteByMessageId(ctx context.Context, messageID string) error {
	return c.db.WithContext(ctx).Where(
		`"message_id" = ?`,
		messageID,
	).Delete(&model.DiscordSendedNotification{}).Error
}

func (c *SendedMessagesStore) DeleteByChannelId(ctx context.Context, messageID string) error {
	return c.db.WithContext(ctx).Where(
		`"channel_id" = ?`,
		messageID,
	).Delete(&model.DiscordSendedNotification{}).Error
}

func (c *SendedMessagesStore) DeleteByGuildId(ctx context.Context, messageID string) error {
	return c.db.WithContext(ctx).Where(
		`"guild_id" = ?`,
		messageID,
	).Delete(&model.DiscordSendedNotification{}).Error
}

func (c *SendedMessagesStore) GetAll(ctx context.Context) ([]Message, error) {
	var entities []model.DiscordSendedNotification

	err := c.db.WithContext(ctx).Find(&entities).Error
	if err != nil {
		return nil, err
	}

	messages := make([]Message, len(entities))
	for i, entity := range entities {
		messages[i] = Message{
			GuildID:          entity.GuildID,
			MessageID:        entity.MessageID,
			TwitchChannelID:  entity.TwitchChannelID,
			DiscordChannelID: entity.DiscordChannelID,
		}
	}

	return messages, nil
}
