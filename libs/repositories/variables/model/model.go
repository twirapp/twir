package model

import (
	"github.com/google/uuid"
	"github.com/guregu/null"
)

type CustomVarType string

const (
	CustomVarScript = "SCRIPT"
	CustomVarText   = "TEXT"
	CustomVarNumber = "NUMBER"
)

type ScriptLanguage string

const (
	ScriptLanguageJavaScript = "javascript"
	ScriptLanguagePython     = "python"
)

type CustomVariable struct {
	ID             uuid.UUID
	Name           string
	Description    null.String
	Type           CustomVarType
	EvalValue      string
	Response       string
	ChannelID      string
	ScriptLanguage ScriptLanguage
}

var Nil = CustomVariable{}
