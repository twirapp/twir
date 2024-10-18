package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null"
)

type AuditLog struct {
	ID            uuid.UUID          `gorm:"column:id;type:uuid"`
	Table         string             `gorm:"column:table_name;type:varchar(255)"`
	OperationType AuditOperationType `gorm:"column:operation_type;type:varchar(255)"`
	OldValue      null.String        `gorm:"column:old_value;type:text"`
	NewValue      null.String        `gorm:"column:new_value;type:text"`
	ObjectID      null.String        `gorm:"column:object_id;type:uuid"`
	ChannelID     null.String        `gorm:"column:channel_id;type:uuid"`
	UserID        null.String        `gorm:"column:user_id;type:uuid"`
	CreatedAt     time.Time          `gorm:"column:created_at;type:timestamp;default:now()"`

	Channel *Users `gorm:"foreignKey:ChannelID"`
	User    *Users `gorm:"foreignKey:UserID"`
}

func (c *AuditLog) TableName() string {
	return "audit_logs"
}

type AuditOperationType string

const (
	AuditOperationCreate AuditOperationType = "CREATE"
	AuditOperationUpdate AuditOperationType = "UPDATE"
	AuditOperationDelete AuditOperationType = "DELETE"
)
