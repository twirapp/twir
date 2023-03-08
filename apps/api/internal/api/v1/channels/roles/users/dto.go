package roles_users

type roleUserDto struct {
	UserNames []string `validate:"required" json:"userNames"`
}
