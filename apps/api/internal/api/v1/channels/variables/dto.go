package variables

import model "github.com/satont/tsuwari/libs/gomodels"

type variableDto struct {
	Name        string              `validate:"required"                   json:"name"`
	Description *string             `                                      json:"description"`
	Type        model.CustomVarType `validate:"required,oneof=SCRIPT TEXT NUMBER" json:"type"        enums:"SCRIPT,TEXT"`
	EvalValue   string              `                                      json:"evalValue"`
	Response    string              `                                      json:"response"`
}
