package commands_group

type groupDto struct {
	Name  string `json:"name" validate:"required"`
	Color string `json:"color" validate:"required"`
}
