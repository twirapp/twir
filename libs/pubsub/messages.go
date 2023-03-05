package pubsub

type StreamUpdateMessage struct {
	ChannelID string `json:"channelId"`
	Title     string `json:"title"`
	Category  string `json:"category"`
}

type StreamOnlineMessage struct {
	ChannelID string `json:"channelId"`
	StreamID  string `json:"streamId"`
}

type StreamOfflineMessage struct {
	ChannelID string `json:"channelId"`
}

type UserUpdateMessage struct {
	UserID        string `json:"user_id"`
	UserLogin     string `json:"user_login"`
	UserName      string `json:"user_name"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Description   string `json:"description"`
}
