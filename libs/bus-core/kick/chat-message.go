package kick

const KickChatMessageSubject = "kick.chat-message"

type KickChatMessage struct {
	MessageID            string      `json:"message_id"`
	BroadcasterUserID    string      `json:"broadcaster_user_id"`
	BroadcasterUserLogin string      `json:"broadcaster_user_login"`
	SenderUserID         string      `json:"sender_user_id"`
	SenderUserLogin      string      `json:"sender_user_login"`
	SenderDisplayName    string      `json:"sender_display_name"`
	Text                 string      `json:"text"`
	Color                string      `json:"color"`
	Badges               []KickBadge `json:"badges,omitempty"`
}

type KickBadge struct {
	SetID string `json:"set_id"`
	Text  string `json:"text"`
}
