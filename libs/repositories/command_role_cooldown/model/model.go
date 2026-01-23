package model

import (
	"time"

	"github.com/google/uuid"
)

type CommandRoleCooldown struct {
	ID        uuid.UUID `db:"id"`
	CommandID uuid.UUID `db:"command_id"`
	RoleID    uuid.UUID `db:"role_id"`
	Cooldown  int       `db:"cooldown"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`

	isNil bool
}

func (c CommandRoleCooldown) IsNil() bool {
	return c.isNil
}

var Nil = CommandRoleCooldown{
	isNil: true,
}
