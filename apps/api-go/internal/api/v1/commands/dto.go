package commands

type responsesDto struct {
	ID    *string `json:"id"`
	Text  string  `json:"text"  validate:"required,min=1,max=500"`
	Order int     `json:"order" validate:"required,gte=0"`
}

type commandDto struct {
	ID      *string `json:"id"`
	Default *bool   `json:"default"`

	Name         string         `validate:"required,min=1,max=50"    json:"name"`
	Cooldown     *int64         `validate:"required,gte=5,lte=86400" json:"cooldown"`
	CooldownType string         `validate:"required"                 json:"cooldownType"`
	Description  *string        `validate:"max=400"                  json:"description"`
	Permission   string         `validate:"required"                 json:"permission"`
	Aliases      []string       `validate:"dive"                     json:"aliases"`
	Visible      *bool          `validate:"required"                 json:"visible"`
	Enabled      *bool          `validate:"required"                 json:"enabled"`
	Responses    []responsesDto `validate:"required,dive"            json:"responses"`
	KeepOrder    *bool          `validate:"required"                 json:"keepOrder"`
	IsReply      *bool          `validate:"required"                 json:"isReply"`
}
