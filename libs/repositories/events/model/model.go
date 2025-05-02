package model

var Nil = Event{}

type Event struct {
	ID          string           `json:"id"`
	ChannelID   string           `json:"channelId"`
	Type        string           `json:"type"`
	RewardID    *string          `json:"rewardId"`
	CommandID   *string          `json:"commandId"`
	KeywordID   *string          `json:"keywordId"`
	Description string           `json:"description"`
	Enabled     bool             `json:"enabled"`
	OnlineOnly  bool             `json:"onlineOnly"`
	Operations  []EventOperation `json:"operations"`
}

type EventOperation struct {
	ID             string                 `json:"id"`
	Type           string                 `json:"type"`
	Input          *string                `json:"input"`
	Delay          int                    `json:"delay"`
	Repeat         int                    `json:"repeat"`
	UseAnnounce    bool                   `json:"useAnnounce"`
	TimeoutTime    int                    `json:"timeoutTime"`
	TimeoutMessage *string                `json:"timeoutMessage"`
	Target         *string                `json:"target"`
	Enabled        bool                   `json:"enabled"`
	Filters        []EventOperationFilter `json:"filters"`
}

type EventOperationFilter struct {
	ID    string `json:"id"`
	Type  string `json:"type"`
	Left  string `json:"left"`
	Right string `json:"right"`
}
