package mappers

import (
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/apps/api-gql/internal/gql/gqlmodel"
	auditlogs "github.com/twirapp/twir/libs/audit-logs"
)

func AuditLogToGql(auditLog auditlogs.AuditLog) *gqlmodel.AuditLog {
	return &gqlmodel.AuditLog{
		ID:            auditLog.ID,
		Table:         auditLog.Table,
		OperationType: AuditLogOperationTypeToGql(auditLog.OperationType),
		OldValue:      auditLog.OldValue.Ptr(),
		NewValue:      auditLog.NewValue.Ptr(),
		ObjectID:      auditLog.ObjectID.Ptr(),
		UserID:        auditLog.UserID.Ptr(),
		CreatedAt:     auditLog.CreatedAt,
	}
}

func AuditLogOperationTypeToGql(t auditlogs.AuditOperationType) gqlmodel.AuditOperationType {
	switch t {
	case auditlogs.AuditOperationTypeUpdate:
		return gqlmodel.AuditOperationTypeUpdate
	case auditlogs.AuditOperationTypeCreate:
		return gqlmodel.AuditOperationTypeCreate
	case auditlogs.AuditOperationTypeDelete:
		return gqlmodel.AuditOperationTypeDelete
	default:
		return ""
	}
}

func AuditTypeModelToGql(t model.AuditOperationType) gqlmodel.AuditOperationType {
	switch t {
	case model.AuditOperationUpdate:
		return gqlmodel.AuditOperationTypeUpdate
	case model.AuditOperationCreate:
		return gqlmodel.AuditOperationTypeCreate
	case model.AuditOperationDelete:
		return gqlmodel.AuditOperationTypeDelete
	default:
		return ""
	}
}

func AuditTypeGqlToModel(t gqlmodel.AuditOperationType) model.AuditOperationType {
	switch t {
	case gqlmodel.AuditOperationTypeUpdate:
		return model.AuditOperationUpdate
	case gqlmodel.AuditOperationTypeCreate:
		return model.AuditOperationCreate
	case gqlmodel.AuditOperationTypeDelete:
		return model.AuditOperationDelete
	default:
		return ""
	}
}
