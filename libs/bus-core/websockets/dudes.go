package websockets

const (
	DudesGrowSubject          = "websockets.dudes.grow"
	DudesUserSettingsSubjsect = "websockets.dudes.change_color"
	DudesLeaveSubject         = "websockets.dudes.leave"
)

type DudesGrowRequest struct {
	ChannelID string `json:"channelId"`
	UserID    string `json:"userId"`
}

type DudesChangeUserSettingsRequest struct {
	ChannelID string `json:"channelId"`
	UserID    string `json:"userId"`
}

type DudesLeaveRequest struct {
	ChannelID string `json:"channelId"`
	UserID    string `json:"userId"`
}
