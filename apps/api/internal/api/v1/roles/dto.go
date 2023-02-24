package roles

type roleDto struct {
	Name        string   `validate:"required" json:"name"`
	Permissions []string `json:"permissions"`
	Users       []string `json:"users"`
}
