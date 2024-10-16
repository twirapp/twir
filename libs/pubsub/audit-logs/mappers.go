package auditlog

import (
	"github.com/guregu/null"
	busauditlog "github.com/twirapp/twir/libs/bus-core/audit-logs"
)

func fromBusNewAuditLogMessage(msg busauditlog.NewAuditLogMessage) AuditLog {
	return AuditLog{
		ID:            msg.ID,
		Table:         msg.Table,
		OperationType: msg.OperationType,
		OldValue:      null.StringFromPtr(msg.OldValue),
		NewValue:      null.StringFromPtr(msg.NewValue),
		ObjectID:      null.StringFromPtr(msg.ObjectID),
		UserID:        null.StringFromPtr(msg.UserID),
		CreatedAt:     msg.CreatedAt,
		ChannelID:     null.StringFromPtr(msg.ChannelID),
	}
}

func toBusNewAuditLogMessage(auditLog AuditLog) busauditlog.NewAuditLogMessage {
	return busauditlog.NewAuditLogMessage{
		ID:            auditLog.ID,
		Table:         auditLog.Table,
		OperationType: auditLog.OperationType,
		OldValue:      auditLog.OldValue.Ptr(),
		NewValue:      auditLog.NewValue.Ptr(),
		ObjectID:      auditLog.ObjectID.Ptr(),
		UserID:        auditLog.UserID.Ptr(),
		CreatedAt:     auditLog.CreatedAt,
		ChannelID:     auditLog.ChannelID.Ptr(),
	}
}
