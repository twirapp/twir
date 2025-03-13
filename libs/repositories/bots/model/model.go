package model

import (
	"github.com/google/uuid"
	tokenmodel "github.com/twirapp/twir/libs/repositories/tokens/model"
)

type Bot struct {
	ID      string
	Type    string
	TokenID *uuid.UUID
	Token   *tokenmodel.Token
}

var Nil = Bot{}
