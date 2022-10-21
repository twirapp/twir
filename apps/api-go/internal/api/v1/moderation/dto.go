package moderation

type moderationDto struct {
	Type               string   `validate:"required" json:"type"`
	Enabled            *bool    `validate:"required" json:"enabled"`
	Subscribers        *bool    `validate:"required" json:"subscribers"`
	Vips               *bool    `validate:"required" json:"vips"`
	BanTime            uint64   `validate:"required" json:"banTime"`
	BanMessage         *string  `                    json:"banMessage"`
	WarningMessage     *string  `                    json:"warningMessage"`
	CheckClips         *bool    `                    json:"checkClips"`
	TriggerLength      uint64   `validate:"required" json:"triggerLength"`
	MaxPercentage      uint64   `validate:"required" json:"maxPercentage"`
	BlackListSentences []string `                    json:"blackListSentences"`
}
