package model

import (
	"github.com/google/uuid"
)

type RoleUser struct {
	ID     uuid.UUID `db:"id"`
	UserID uuid.UUID `db:"userId"`
	RoleID uuid.UUID `db:"roleId"`
}

var RoleUserNil = RoleUser{}
