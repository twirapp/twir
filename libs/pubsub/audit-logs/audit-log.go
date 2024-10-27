package auditlog

import (
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null"
	auditlog "github.com/twirapp/twir/libs/bus-core/audit-logs"
)

type AuditLog struct {
	ID            uuid.UUID
	Table         string
	OperationType auditlog.AuditOperationType
	OldValue      null.String
	NewValue      null.String
	ObjectID      null.String
	ChannelID     null.String
	UserID        null.String
	CreatedAt     time.Time
}
