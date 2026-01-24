package commandrolecooldownentity

import (
	"time"

	"github.com/google/uuid"
)

type CommandRoleCooldown struct {
	ID        uuid.UUID
	CommandID uuid.UUID
	RoleID    uuid.UUID
	Cooldown  int
	CreatedAt time.Time
	UpdatedAt time.Time

	isNil bool
}

func (c CommandRoleCooldown) IsNil() bool {
	return c.isNil
}

var Nil = CommandRoleCooldown{
	isNil: true,
}
