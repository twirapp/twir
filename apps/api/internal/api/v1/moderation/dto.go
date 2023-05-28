package moderation

type moderationItemDto struct {
	ID                 string   `validate:"required" json:"id"`
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

type moderationDto struct {
	Items []moderationItemDto `validate:"required,dive"`
}

type postTitleDto struct {
	Title string `validate:"required" json:"title"`
}

type postTitleResponse struct {
	Title string `json:"title"`
}

type postCategoryDto struct {
	CategoryId string `validate:"required" json:"categoryId"`
}

type postCategoryResponse struct {
	CategoryId string `json:"categoryId"`
}
