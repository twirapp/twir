package keywords

type keywordDto struct {
	Text      string `validate:"required,max=500" json:"text"`
	Response  string `validate:"max=500" json:"response"`
	Enabled   *bool  `validate:"required"         json:"enabled"`
	Cooldown  uint64 `validate:"gte=5,lte=86400"  json:"cooldown"`
	IsReply   *bool  `validate:"required"         json:"isReply"`
	IsRegular *bool  `validate:"required"         json:"isRegular"`
	Usages    *int   `validate:"required" json:"usages"`
}

type keywordPatchDto struct {
	Enabled *bool `validate:"omitempty,required" json:"enabled,omitempty"`
}
