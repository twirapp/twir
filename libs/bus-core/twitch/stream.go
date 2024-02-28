package twitch

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
