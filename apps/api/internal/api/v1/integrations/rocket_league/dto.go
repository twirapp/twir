package rocket_league

type createOrUpdateDTO struct {
	Username string `validate:"required" json:"username"`
	Code     string `validate:"required" json:"code"`
}
