package variables

type variableDto struct {
	Name        string  `validate:"required"                   json:"name"`
	Description *string `                                      json:"description"`
	Type        string  `validate:"required,oneof=SCRIPT TEXT" json:"type"`
	EvalValue   *string `validate:"max=10000"                  json:"evalValue"`
	Response    *string `validate:"max=500"                    json:"response"`
}
