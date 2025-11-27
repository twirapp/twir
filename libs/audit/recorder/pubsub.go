package recorder

import (
	"context"
	"fmt"
	"time"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/twirapp/twir/libs/audit"
	buscoreauditlogs "github.com/twirapp/twir/libs/bus-core/audit-logs"
	auditlogs "github.com/twirapp/twir/libs/pubsub/audit-logs"
)

// PubSub is an [audit.Recorder] implementation with [auditlogs.PubSub].
type PubSub struct {
	pubSub auditlogs.PubSub
}

var _ audit.Recorder = (*PubSub)(nil)

func NewPubSub(pubSub auditlogs.PubSub) PubSub {
	return PubSub{
		pubSub: pubSub,
	}
}

func (p PubSub) RecordCreateOperation(ctx context.Context, operation audit.CreateOperation) error {
	return p.publishOperation(
		ctx,
		operation.Metadata,
		buscoreauditlogs.AuditOperationTypeCreate,
		operation.NewValue,
		nil,
	)
}

func (p PubSub) RecordDeleteOperation(ctx context.Context, operation audit.DeleteOperation) error {
	return p.publishOperation(
		ctx,
		operation.Metadata,
		buscoreauditlogs.AuditOperationTypeDelete,
		nil,
		operation.OldValue,
	)
}

func (p PubSub) RecordUpdateOperation(ctx context.Context, operation audit.UpdateOperation) error {
	return p.publishOperation(
		ctx,
		operation.Metadata,
		buscoreauditlogs.AuditOperationTypeUpdate,
		operation.NewValue,
		operation.OldValue,
	)
}

func (p PubSub) publishOperation(
	ctx context.Context,
	metadata audit.OperationMetadata,
	operationType buscoreauditlogs.AuditOperationType,
	newValue any,
	oldValue any,
) error {
	auditLog := auditlogs.AuditLog{
		ID:            uuid.New(),
		Table:         metadata.System,
		OperationType: operationType,
		ObjectID:      null.StringFromPtr(metadata.ObjectID),
		ChannelID:     null.StringFromPtr(metadata.ChannelID),
		UserID:        null.StringFromPtr(metadata.ActorID),
		CreatedAt:     time.Now(),
	}

	if newValue != nil {
		newValueBytes, err := json.Marshal(newValue)
		if err != nil {
			return fmt.Errorf("marshal operation new value: %w", err)
		}

		auditLog.NewValue = null.StringFrom(string(newValueBytes))
	}

	if oldValue != nil {
		oldValueBytes, err := json.Marshal(oldValue)
		if err != nil {
			return fmt.Errorf("marshal operation old value: %w", err)
		}

		auditLog.OldValue = null.StringFrom(string(oldValueBytes))
	}

	if err := p.pubSub.Publish(ctx, auditLog); err != nil {
		return fmt.Errorf("publish audit log to pubsub: %w", err)
	}

	return nil
}
