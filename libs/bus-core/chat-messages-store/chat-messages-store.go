package chat_messages_store

type GetChatMessagesByTextRequest struct {
	ChannelID string `json:"channel_id"`
	Text      string `json:"text"`
}

type GetChatMessagesByTextResponse struct {
	Messages []StoredChatMessage `json:"messages"`
}

type StoredChatMessage struct {
	MessageID    string `json:"message_id"`
	ChannelID    string `json:"channel_id"`
	UserID       string `json:"user_id"`
	UserLogin    string `json:"user_login"`
	Text         string `json:"text"`
	CanBeDeleted bool   `json:"can_be_deleted"`
}
