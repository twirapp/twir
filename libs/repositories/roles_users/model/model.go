package model

import (
	"github.com/google/uuid"
)

type RoleUser struct {
	ID     uuid.UUID
	UserID uuid.UUID
	RoleID uuid.UUID
}

var RoleUserNil = RoleUser{}
