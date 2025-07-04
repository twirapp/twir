package audit_logs

import (
	"context"

	"github.com/twirapp/twir/libs/repositories/audit_logs/model"
)

type Repository interface {
	GetMany(ctx context.Context, input GetManyInput) (
		[]model.AuditLog,
		error,
	)
	Count(ctx context.Context, input GetCountInput) (uint64, error)
	Create(ctx context.Context, input CreateInput) error
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

type CreateInput struct {
	Table         string
	OperationType model.AuditOperationType
	OldValue      *string
	NewValue      *string
	ObjectID      *string
	ChannelID     *string
	UserID        *string
}
