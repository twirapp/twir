package timers

type responseDto struct {
	Text       string `validate:"required,max=400" json:"text"`
	IsAnnounce *bool  `validate:"required"         json:"isAnnounce"`
}

type timerDto struct {
	Name            string        `validate:"required,max=510" json:"name"`
	Enabled         *bool         `validate:"required"         json:"enabled"`
	TimeInterval    uint64        `validate:"gte=1,lte=120"    json:"timeInterval"`
	MessageInterval uint64        `validate:"gte=0,lte=10000"  json:"messageInterval"`
	Responses       []responseDto `validate:"required,dive"    json:"responses"`
}

type timerPatchDto struct {
	Enabled *bool `validate:"required,omitempty" json:"enabled,omitempty"`
}
