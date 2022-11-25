package keywords

type keywordDto struct {
	Text      string `validate:"required,max=100" json:"text"`
	Response  string `validate:"max=400" json:"response"`
	Enabled   *bool  `validate:"required"         json:"enabled"`
	Cooldown  uint64 `validate:"gte=5,lte=86400"  json:"cooldown"`
	IsReply   *bool  `validate:"required"         json:"isReply"`
	IsRegular *bool  `validate:"required"         json:"isRegular"`
	Usages    *int   `validate:"required" json:"usages"`
}
