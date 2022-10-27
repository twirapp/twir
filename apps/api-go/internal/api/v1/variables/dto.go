package variables

type variableDto struct {
	Name        string  `validate:"required"                   json:"name"`
	Description *string `                                      json:"description"`
	Type        string  `validate:"required,oneof=SCRIPT TEXT" json:"type"`
	EvalValue   *string `                                      json:"evalValue"`
	Response    *string `                                      json:"response"`
}
