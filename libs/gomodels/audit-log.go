package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type AuditLog struct {
	ID            uuid.UUID          `gorm:"column:id;type:uuid"`
	Table         string             `gorm:"column:table_name;type:varchar(255)"`
	OperationType AuditOperationType `gorm:"column:operation_type;type:varchar(255)"`
	OldValue      sql.Null[string]   `gorm:"column:old_value;type:text"`
	NewValue      sql.Null[string]   `gorm:"column:new_value;type:text"`
	ObjectID      sql.Null[string]   `gorm:"column:object_id;type:uuid"`
	UserID        sql.Null[string]   `gorm:"column:user_id;type:uuid"`
	DashboardID   sql.Null[string]   `gorm:"column:dashboard_id;type:uuid"`
	CreatedAt     time.Time          `gorm:"column:created_at;type:timestamp;default:now()"`

	User      *Users `gorm:"foreignKey:UserID"`
	Dashboard *Users `gorm:"foreignKey:DashboardID"`
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
