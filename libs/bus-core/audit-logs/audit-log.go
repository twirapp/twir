package auditlog

import (
	"time"

	"github.com/google/uuid"
	auditlogs "github.com/twirapp/twir/libs/audit-logs"
)

type NewAuditLogMessage struct {
	ID            uuid.UUID                    `json:"id"`
	Table         string                       `json:"table"`
	OperationType auditlogs.AuditOperationType `json:"operation_type"`
	OldValue      *string                      `json:"old_value,omitempty"`
	NewValue      *string                      `json:"new_value,omitempty"`
	ObjectID      *string                      `json:"object_id,omitempty"`
	UserID        *string                      `json:"user_id,omitempty"`
	CreatedAt     time.Time                    `json:"created_at"`
}
