package auditlog

import (
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null"
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
	ChannelID     null.String
	UserID        null.String
	CreatedAt     time.Time
}

type Params struct {
	ID            uuid.UUID
	Table         string
	OperationType AuditOperationType
	CreatedAt     time.Time
}

func NewAuditLog(params Params, opts ...Option) AuditLog {
	auditLog := AuditLog{
		ID:            params.ID,
		Table:         params.Table,
		OperationType: params.OperationType,
		CreatedAt:     params.CreatedAt,
	}

	for _, opt := range opts {
		opt(&auditLog)
	}

	return auditLog
}
