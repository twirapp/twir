package spotify

type tokenDto struct {
	Code string `validate:"required" json:"code"`
}
