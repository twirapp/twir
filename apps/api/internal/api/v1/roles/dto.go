package roles

type roleDto struct {
	Name  string   `validate:"required" json:"name"`
	Flags []string `json:"flags"`
	Users []string `json:"users"`
}
