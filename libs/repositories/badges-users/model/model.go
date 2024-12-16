package model

import (
	"time"

	"github.com/google/uuid"
)

type BadgeUser struct {
	ID        uuid.UUID
	BadgeID   uuid.UUID
	UserID    string
	CreatedAt time.Time
}

var Nil = BadgeUser{}
