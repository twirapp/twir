package faceit

type tokenDto struct {
	Code string `validate:"required" json:"code"`
}
