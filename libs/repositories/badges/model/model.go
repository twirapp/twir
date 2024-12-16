package model

import (
	"time"

	"github.com/google/uuid"
)

type Badge struct {
	ID        uuid.UUID
	Name      string
	Enabled   bool
	CreatedAt time.Time
	FileName  string
	FFZSlot   int
}

var Nil = Badge{}
