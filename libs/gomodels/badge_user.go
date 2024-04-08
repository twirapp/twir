package model

import (
	"time"

	"github.com/google/uuid"
)

type BadgeUser struct {
	ID        uuid.UUID `gorm:"primaryKey;column:id;type:uuid;" json:"id"`
	BadgeID   uuid.UUID `gorm:"column:badge_id;type:uuid;" json:"badge_id"`
	UserID    string    `gorm:"column:user_id;type:uuid;" json:"user_id"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp without time zone;default:now()" json:"created_at"`

	Badge Badge `gorm:"foreignKey:BadgeID;references:ID" json:"badge"`
}

func (BadgeUser) TableName() string {
	return "badges_users"
}
