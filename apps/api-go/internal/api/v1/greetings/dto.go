package greetings

type greetingsDto struct {
	ID *string `json:"id,omitempty"`

	Username string `validate:"required,min=1" json:"username"`
	Text     string `validate:"max=400"        json:"text"`
	Enabled  *bool  `validate:"required"       json:"enabled"`
	IsReply  *bool  `validate:"required"       json:"isReply"`
}
