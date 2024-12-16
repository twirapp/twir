package entity

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
	FileURL   string
}

type BadgeUser struct {
	ID        uuid.UUID
	BadgeID   uuid.UUID
	UserID    string
	CreatedAt time.Time
}

type BadgeWithUsers struct {
	Badge
	Users []string
}

var BadgeNil = Badge{}
var BadgeUserNil = BadgeUser{}
