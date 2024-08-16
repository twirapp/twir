package chat_messages_store

type GetChatMessagesByTextRequest struct {
	ChannelID string
	Text      string
}

type GetChatMessagesByTextResponse struct {
	Messages []StoredChatMessage
}

type StoredChatMessage struct {
	RedisID string

	MessageID    string
	ChannelID    string
	UserID       string
	UserLogin    string
	Text         string
	CanBeDeleted bool
}

type RemoveMessagesRequest struct {
	MessagesRedisIDS []string
}
