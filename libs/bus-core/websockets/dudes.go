package websockets

type DudesGrowRequest struct {
	ChannelID string `json:"channelId"`
	UserID    string `json:"userId"`
	UserName  string `json:"userName"`
	UserLogin string `json:"userLogin"`
	Color     string `json:"color"`
}

type DudesChangeColorRequest struct {
	ChannelID string `json:"channelId"`
	UserID    string `json:"userId"`
	UserName  string `json:"userName"`
	UserLogin string `json:"userLogin"`
	Color     string `json:"color"`
}
