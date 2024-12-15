package audit_logs

import (
	"context"

	"github.com/twirapp/twir/libs/repositories/audit-logs/model"
)

type Repository interface {
	GetMany(ctx context.Context, input GetManyInput) (
		[]model.AuditLog,
		error,
	)
	Count(ctx context.Context, input GetCountInput) (int, error)
}

type GetManyInput struct {
	ChannelID      *string
	ActorID        *string
	ObjectID       *string
	Limit          int
	Page           int
	Systems        []string
	OperationTypes []model.AuditOperationType
}

type GetCountInput struct {
	ChannelID *string
}
