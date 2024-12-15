package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null"
)

type AuditLog struct {
	ID            uuid.UUID
	Table         string
	OperationType AuditOperationType
	OldValue      null.String
	NewValue      null.String
	ObjectID      null.String
	ChannelID     null.String
	UserID        null.String
	CreatedAt     time.Time
}

type AuditOperationType string

const (
	AuditOperationCreate      AuditOperationType = "CREATE"
	AuditOperationUpdate      AuditOperationType = "UPDATE"
	AuditOperationDelete      AuditOperationType = "DELETE"
	AuditOperationTypeUnknown AuditOperationType = "UNKNOWN"
)
