package sended_messages_store

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
)

type Message struct {
	GuildID   string `redis:"guildId"`
	MessageID string `redis:"messageId"`
	ChannelID string `redis:"channelId"`
}

type SendedMessagesStore struct {
	redis *redis.Client
}

type Opts struct {
	fx.In

	Redis *redis.Client
}

func New(opts Opts) *SendedMessagesStore {
	return &SendedMessagesStore{
		redis: opts.Redis,
	}
}

func (c *SendedMessagesStore) buildMessageKey(id string) string {
	return fmt.Sprintf("discord:sended_messages:%s", id)
}

func (c *SendedMessagesStore) Add(ctx context.Context, msg Message) error {
	key := c.buildMessageKey(msg.MessageID)

	if err := c.redis.HSet(ctx, key, msg).Err(); err != nil {
		return err
	}

	if err := c.redis.Expire(ctx, key, 24*time.Hour).Err(); err != nil {
		return err
	}

	return nil
}

func (c *SendedMessagesStore) Get(ctx context.Context, messageID string) (Message, error) {
	msg := Message{}
	err := c.redis.HGetAll(ctx, c.buildMessageKey(messageID)).Scan(&msg)
	return msg, err
}

func (c *SendedMessagesStore) Delete(ctx context.Context, messageID string) error {
	return c.redis.Del(ctx, c.buildMessageKey(messageID)).Err()
}

func (c *SendedMessagesStore) GetAll(ctx context.Context) ([]Message, error) {
	iter := c.redis.Scan(ctx, 0, c.buildMessageKey("*"), 0).Iterator()

	var keys []string

	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}
	if err := iter.Err(); err != nil {
		return nil, err
	}

	var messages []Message
	for _, key := range keys {
		msgId := strings.Split(key, ":")[2]
		msg, err := c.Get(ctx, msgId)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, nil
}
