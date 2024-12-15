package entity

import (
	"time"

	"github.com/google/uuid"
)

type AuditLog struct {
	ID            uuid.UUID
	TableName     string
	OperationType AuditOperationType
	OldValue      *string
	NewValue      *string
	ObjectID      *string
	ChannelID     *string
	UserID        *string
	CreatedAt     time.Time
}

type AuditOperationType string

const (
	AuditOperationCreate      AuditOperationType = "CREATE"
	AuditOperationUpdate      AuditOperationType = "UPDATE"
	AuditOperationDelete      AuditOperationType = "DELETE"
	AuditOperationTypeUnknown AuditOperationType = "UNKNOWN"
)
