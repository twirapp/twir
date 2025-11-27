package recorder

import (
	"context"
	"fmt"

	"github.com/goccy/go-json"
	"github.com/samber/lo"
	"github.com/twirapp/twir/libs/audit"
	auditlogs "github.com/twirapp/twir/libs/repositories/audit_logs"
	"github.com/twirapp/twir/libs/repositories/audit_logs/model"
)

// Database is an [audit.Recorder] implementation with [auditlogs.Repository].
type Database struct {
	repository auditlogs.Repository
}

var _ audit.Recorder = (*Database)(nil)

func NewDatabase(repository auditlogs.Repository) Database {
	return Database{
		repository: repository,
	}
}

func (d Database) RecordCreateOperation(
	ctx context.Context,
	operation audit.CreateOperation,
) error {
	return d.saveOperation(
		ctx,
		operation.Metadata,
		model.AuditOperationCreate,
		operation.NewValue,
		nil,
	)
}

func (d Database) RecordDeleteOperation(
	ctx context.Context,
	operation audit.DeleteOperation,
) error {
	return d.saveOperation(
		ctx,
		operation.Metadata,
		model.AuditOperationDelete,
		nil,
		operation.OldValue,
	)
}

func (d Database) RecordUpdateOperation(
	ctx context.Context,
	operation audit.UpdateOperation,
) error {
	return d.saveOperation(
		ctx,
		operation.Metadata,
		model.AuditOperationUpdate,
		operation.NewValue,
		operation.OldValue,
	)
}

func (d Database) saveOperation(
	ctx context.Context,
	metadata audit.OperationMetadata,
	operationType model.AuditOperationType,
	newValue any,
	oldValue any,
) error {
	auditLog := auditlogs.CreateInput{
		Table:         metadata.System,
		OperationType: operationType,
		ObjectID:      metadata.ObjectID,
		ChannelID:     metadata.ChannelID,
		UserID:        metadata.ActorID,
	}

	if newValue != nil {
		newValueBytes, err := json.Marshal(newValue)
		if err != nil {
			return fmt.Errorf("marshal operation new value: %w", err)
		}

		auditLog.NewValue = lo.ToPtr(string(newValueBytes))
	}

	if oldValue != nil {
		oldValueBytes, err := json.Marshal(oldValue)
		if err != nil {
			return fmt.Errorf("marshal operation old value: %w", err)
		}

		auditLog.OldValue = lo.ToPtr(string(oldValueBytes))
	}

	if err := d.repository.Create(ctx, auditLog); err != nil {
		return fmt.Errorf("create audit log in database: %w", err)
	}

	return nil
}
