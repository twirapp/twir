package entity

var EventNil = Event{}

type Event struct {
	ID          string
	ChannelID   string
	Type        string
	RewardID    *string
	CommandID   *string
	KeywordID   *string
	Description string
	Enabled     bool
	OnlineOnly  bool
	Operations  []EventOperation
}

type EventOperation struct {
	ID             string
	Type           string
	Input          *string
	Delay          int
	Repeat         int
	UseAnnounce    bool
	TimeoutTime    int
	TimeoutMessage *string
	Target         *string
	Enabled        bool
	Filters        []EventOperationFilter
}

type EventOperationFilter struct {
	ID    string
	Type  string
	Left  string
	Right string
}
