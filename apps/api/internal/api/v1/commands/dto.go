package commands

type responsesDto struct {
	Text  string `validate:"min=1,max=500" json:"text"`
	Order uint8  `validate:"gte=0"                                             json:"order"`
}

type commandDto struct {
	ID *string `json:"id,omitempty"`

	Name               string         `validate:"required,min=1,max=100"       json:"name"`
	Cooldown           uint32         `validate:"lte=86400"                   json:"cooldown"`
	CooldownType       string         `validate:"required"                    json:"cooldownType"`
	Description        *string        `validate:"omitempty,max=500"           json:"description,omitempty"`
	Aliases            []string       `validate:"required,dive,max=20" json:"aliases"`
	Visible            *bool          `validate:"omitempty,required"          json:"visible,omitempty"`
	Enabled            *bool          `validate:"omitempty,required"          json:"enabled,omitempty"`
	Responses          []responsesDto `validate:"dive"                        json:"responses"`
	KeepResponsesOrder *bool          `validate:"required"                    json:"keepResponsesOrder"`
	IsReply            *bool          `validate:"omitempty,required"          json:"isReply,omitempty"`
	GroupID            *string        `json:"groupId,omitempty"`
	DeniedUsersIds     []string       `json:"deniedUsersIds"`
	AllowedUsersIds    []string       `json:"allowedUsersIds"`
	RolesIDS           []string       `json:"rolesIds"`
	OnlineOnly         *bool          `validate:"required" json:"onlineOnly"`
}

type commandPatchDto struct {
	Enabled *bool `validate:"omitempty,required" json:"enabled,omitempty"`
}
