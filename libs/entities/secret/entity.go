package secret

import (
	"time"

	"github.com/google/uuid"
)

type Secret struct {
	ID          uuid.UUID
	Name        string
	Description *string
	Value       string
	ChannelID   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	isNil       bool
}

func (s Secret) IsNil() bool {
	return s.isNil
}

var Nil = Secret{
	isNil: true,
}
