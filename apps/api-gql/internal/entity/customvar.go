package entity

import (
	"github.com/google/uuid"
)

type CustomVarType string

const (
	CustomVarScript CustomVarType = "SCRIPT"
	CustomVarText   CustomVarType = "TEXT"
	CustomVarNumber CustomVarType = "NUMBER"
)

type CustomVarScriptLanguage string

func (s CustomVarScriptLanguage) String() string {
	return string(s)
}

const (
	ScriptLanguageJavaScript = "javascript"
	ScriptLanguagePython     = "python"
)

type CustomVariable struct {
	ID             uuid.UUID
	ChannelID      string
	Name           string
	Description    *string
	Type           CustomVarType
	EvalValue      string
	Response       string
	ScriptLanguage CustomVarScriptLanguage
}

var CustomVarNil = CustomVariable{}
