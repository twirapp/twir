package mappers

import (
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/apps/api-gql/internal/gql/gqlmodel"
)

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
