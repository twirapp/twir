package auditlog

import (
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v5"
)

type AuditOperationType string

const (
	AuditOperationTypeCreate AuditOperationType = "CREATE"
	AuditOperationTypeUpdate AuditOperationType = "UPDATE"
	AuditOperationTypeDelete AuditOperationType = "DELETE"
)

type AuditLog struct {
	ID            uuid.UUID
	Table         string
	OperationType AuditOperationType
	OldValue      null.String
	NewValue      null.String
	ObjectID      null.String
	UserID        null.String
	CreatedAt     time.Time
}

type Options struct {
	ID            uuid.UUID
	Table         string
	OperationType AuditOperationType
	CreatedAt     time.Time
}

func NewAuditLog(opts Options, optionals ...Option) AuditLog {
	auditLog := AuditLog{
		ID:            opts.ID,
		Table:         opts.Table,
		OperationType: opts.OperationType,
		CreatedAt:     opts.CreatedAt,
	}

	for _, opt := range optionals {
		opt(&auditLog)
	}

	return auditLog
}
