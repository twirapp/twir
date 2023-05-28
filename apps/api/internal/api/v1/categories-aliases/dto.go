package categories_aliases

type categoryAliasDto struct {
	ID         *string `json:"id,omitempty"`
	Category   string  `validate:"required,min=1,max=300" json:"category,omitempty"`
	CategoryId string  `validate:"required" json:"categoryId,omitempty"`
	Alias      string  `validate:"required,min=1,max=300" json:"alias,omitempty"`
}

type categoryAliasPatchDto struct {
	Category   *string `validate:"omitempty,min=1,max=300" json:"category,omitempty"`
	CategoryId *string `validate:"omitempty" json:"categoryId,omitempty"`
	Alias      *string `validate:"omitempty,min=1,max=300" json:"alias,omitempty"`
}
