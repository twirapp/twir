package commands

type responsesDto struct {
	Text  string `validate:"required,min=1,max=http.StatusInternalServerError" json:"text"`
	Order uint8  `validate:"gte=0"                                             json:"order"`
}

type restrictionsDto struct {
	Type  string `validate:"oneof=WATCHED MESSAGES" json:"type"  enums:"WATCHED,MESSAGES"`
	Value string `validate:"max=50"                 json:"value"`
}

type commandDto struct {
	ID *string `json:"id,omitempty"`

	Name               string            `validate:"required,min=1,max=50"       json:"name"`
	Cooldown           uint32            `validate:"lte=86400"                   json:"cooldown"`
	CooldownType       string            `validate:"required"                    json:"cooldownType"`
	Description        *string           `validate:"omitempty,max=400"           json:"description,omitempty"`
	Permission         string            `validate:"required"                    json:"permission"`
	Aliases            []string          `validate:"max=20,dive,required,max=50" json:"aliases"`
	Visible            *bool             `validate:"omitempty,required"          json:"visible,omitempty"`
	Enabled            *bool             `validate:"omitempty,required"          json:"enabled,omitempty"`
	Responses          []responsesDto    `                                       json:"responses"`
	KeepResponsesOrder *bool             `validate:"required"                    json:"keepResponsesOrder"`
	IsReply            *bool             `validate:"omitempty,required"          json:"isReply,omitempty"`
	Restrictions       []restrictionsDto `validate:"dive"                        json:"restrictions"`
}
