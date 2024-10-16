package auditlog

import (
	"time"

	"github.com/google/uuid"
)

type NewAuditLogMessage struct {
	ID            uuid.UUID
	Table         string
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
	AuditOperationTypeCreate AuditOperationType = "CREATE"
	AuditOperationTypeUpdate AuditOperationType = "UPDATE"
	AuditOperationTypeDelete AuditOperationType = "DELETE"
)
