package timers

type responseDto struct {
	Text       string `validate:"required,max=400" json:"text"`
	IsAnnounce *bool  `validate:"required"         json:"announce"`
}

type timerDto struct {
	Name            string        `validate:"required,max=510"       json:"name"`
	Enabled         *bool         `validate:"required"               json:"enabled"`
	TimeInterval    uint64        `validate:"required,gte=1,lte=120" json:"timeInterval"`
	MessageInterval uint64        `validate:"required,lte=10000"     json:"messagInteval"`
	Responses       []responseDto `validate:"required,dive"          json:"responses"`
}
