package model

import (
	"time"

	"github.com/google/uuid"
)

type Badge struct {
	ID        uuid.UUID `gorm:"primaryKey;column:id;type:uuid;" json:"id"`
	Name      string    `gorm:"column:name;type:varchar(255)" json:"name"`
	Enabled   bool      `gorm:"column:enabled;type:boolean;default:true" json:"enabled"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp" json:"created_at"`
	FileName  string    `gorm:"column:file_name;type:varchar(255)" json:"file_name"`
	FFZSlot   int       `gorm:"column:ffz_slot;type:int" json:"ffz_slot"`

	Users []BadgeUser `gorm:"foreignKey:BadgeID;references:ID" json:"users"`
}

func (Badge) TableName() string {
	return "badges"
}
