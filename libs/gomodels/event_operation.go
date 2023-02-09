package model

import (
	"database/sql"
	"time"

	"github.com/guregu/null"
	uuid "github.com/satori/go.uuid"
)

var (
	_ = time.Second
	_ = sql.LevelDefault
	_ = null.Bool{}
	_ = uuid.UUID{}
)

type EventOperation struct {
	ID      string   `gorm:"primary_key;AUTO_INCREMENT;column:id;type:TEXT;" json:"id"`
	Type    string   `gorm:"column:type;type:TEXT;"                     json:"type"`
	Delay   null.Int `gorm:"column:delay;type:int" json:"delay"`
	EventID string   `gorm:"column:eventId;type:string" json:"eventId"`

	Input null.String `gorm:"column:input;type:string" json:"input"`
}

func (c *EventOperation) TableName() string {
	return "channels_events_operations"
}
