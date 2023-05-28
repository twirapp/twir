package greetings

type greetingsDto struct {
	ID *string `json:"id,omitempty"`

	Username string `validate:"required,min=1" json:"userName"`
	Text     string `validate:"max=400"        json:"text"`
	Enabled  *bool  `validate:"required"       json:"enabled"`
	IsReply  *bool  `validate:"required"       json:"isReply"`
}

type greetingsPatchDto struct {
	Enabled *bool `validate:"required" json:"enabled"`
}
