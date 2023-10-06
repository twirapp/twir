package sended_messages_store

import (
	"github.com/satont/twir/libs/utils"
)

type Message struct {
	GuildID   string
	MessageID string
	ChannelID string
}

type SendedMessagesStore struct {
	store *utils.SyncMap[Message]
}

func New() *SendedMessagesStore {
	return &SendedMessagesStore{
		store: utils.NewSyncMap[Message](),
	}
}

func (c *SendedMessagesStore) Add(msg Message) {
	c.store.Add(msg.MessageID, msg)
}

func (c *SendedMessagesStore) Get(messageID string) (Message, bool) {
	return c.store.Get(messageID)
}

func (c *SendedMessagesStore) Delete(messageID string) {
	c.store.Delete(messageID)
}

func (c *SendedMessagesStore) GetAll() []Message {
	return c.store.GetAll()
}
