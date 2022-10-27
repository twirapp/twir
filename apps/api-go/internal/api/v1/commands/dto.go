package commands

type responsesDto struct {
	Text  string `validate:"required,min=1,max=http.StatusInternalServerError" json:"text"`
	Order uint   `validate:"gte=0"                                             json:"order"`
}

type commandDto struct {
	ID *string `json:"id,omitempty"`

	Name               string         `validate:"required,min=1,max=50" json:"name"`
	Cooldown           uint64         `validate:"gte=5,lte=86400"       json:"cooldown"`
	CooldownType       string         `validate:"required"              json:"cooldownType"`
	Description        *string        `validate:"omitempty,max=400"     json:"description,omitempty"`
	Permission         string         `validate:"required"              json:"permission"`
	Aliases            []string       `validate:"dive"                  json:"aliases"`
	Visible            *bool          `validate:"omitempty,required"    json:"visible,omitempty"`
	Enabled            *bool          `validate:"omitempty,required"    json:"enabled,omitempty"`
	Responses          []responsesDto `validate:"required,min=1,dive"   json:"responses"`
	KeepResponsesOrder *bool          `validate:"omitempty,required"    json:"keepResponsesOrder,omitempty"`
	IsReply            *bool          `validate:"omitempty,required"    json:"isReply,omitempty"`
}
