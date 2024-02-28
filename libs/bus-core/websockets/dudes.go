package websockets

type DudesGrowRequest struct {
	ChannelID       string `json:"channelId"`
	UserID          string `json:"userId"`
	UserDisplayName string `json:"userDisplayName"`
	UserName        string `json:"userName"`
	UserColor       string `json:"userColor"`
}

type DudesChangeUserSettingsRequest struct {
	ChannelID string `json:"channelId"`
	UserID    string `json:"userId"`
}
