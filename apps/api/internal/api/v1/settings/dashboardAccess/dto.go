package dashboardaccess

type addUserDto struct {
	UserName string `validate:"required" json:"username"`
}
