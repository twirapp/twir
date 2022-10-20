package commands

type responsesDto struct {
	ID    *string
	Text  string
	Order int
}

type commandDto struct {
	ID      *string `json:"id"`
	Default *bool   `json:"default"`

	Name         string         `validate:"required,min=1,max=50"    json:"name"`
	Cooldown     *int           `validate:"required,gte=5,lte=86400" json:"cooldown"`
	CooldownType string         `validate:"required"                 json:"cooldownType"`
	Description  *string        `validate:"max=400"                  json:"description"`
	Permission   string         `validate:"required"                 json:"permission"`
	Aliases      []string       `validate:"dive"                     json:"aliases"`
	Visible      *bool          `                                    json:"visible"`
	Enabled      *bool          `                                    json:"enabled"`
	Responses    []responsesDto `validate:"required,dive"            json:"responses"`
	KeepOrder    *bool          `                                    json:"keepOrder"`
}
