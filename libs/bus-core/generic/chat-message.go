package generic

type ChatMessage struct {
	Platform          string             `json:"platform"`
	ChannelID         string             `json:"channel_id"`
	UserID            string             `json:"user_id"`
	PlatformChannelID string             `json:"platform_channel_id"`
	SenderID          string             `json:"sender_id"`
	SenderLogin       string             `json:"sender_login"`
	SenderDisplayName string             `json:"sender_display_name"`
	MessageID         string             `json:"message_id"`
	Text              string             `json:"text"`
	Badges            []ChatMessageBadge `json:"badges,omitempty"`
	Color             string             `json:"color"`
	Emotes            []ChatMessageEmote `json:"emotes,omitempty"`
}

type ChatMessageBadge struct {
	SetID string `json:"set_id"`
	Text  string `json:"text"`
}

type ChatMessageEmote struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}
