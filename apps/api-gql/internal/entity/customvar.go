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

type CustomVariable struct {
	ID          uuid.UUID
	ChannelID   string
	Name        string
	Description *string
	Type        CustomVarType
	EvalValue   string
	Response    string
}

var CustomVarNil = CustomVariable{}
