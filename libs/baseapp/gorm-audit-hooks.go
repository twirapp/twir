package baseapp

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/guregu/null"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/logger"
	auditlogs "github.com/satont/twir/libs/pubsub/audit-logs"
	"gorm.io/gorm"
)

type auditHook func(tx *gorm.DB) model.AuditLog

type gormAuditHooks struct {
	logger          logger.Logger
	auditLogsPubSub auditlogs.PubSub
}

func (c *gormAuditHooks) create(tx *gorm.DB) model.AuditLog {
	if tx.Statement.Schema != nil && tx.Statement.Schema.Table == "audit_logs" || tx.Error != nil {
		return model.AuditLog{}
	}

	recordMap, err := getDataBeforeOperation(tx)
	if err != nil {
		return model.AuditLog{}
	}

	objId := getKeyFromData("id", recordMap)
	if objId == "" {
		return model.AuditLog{}
	}

	ctx := tx.Statement.Context
	userID := getUserIDFromContext(ctx)
	dashboardID := getDashboardIDFromContext(ctx)

	audit := model.AuditLog{
		ID:            uuid.New(),
		Table:         tx.Statement.Schema.Table,
		OperationType: model.AuditOperationCreate,
		NewValue:      null.StringFrom(prepareData(recordMap)),
		ObjectID:      null.StringFrom(objId),
		UserID:        null.StringFromPtr(userID),
		ChannelID:     null.StringFromPtr(dashboardID),
	}

	if err := tx.Session(
		&gorm.Session{
			SkipHooks: true,
			NewDB:     true,
		},
	).Create(&audit).Error; err != nil {
		c.logger.Error("error in audit log creation", slog.Any("err", err))
		return model.AuditLog{}
	}

	return audit
}

func (c *gormAuditHooks) delete(tx *gorm.DB) model.AuditLog {
	if tx.Statement.Schema != nil && tx.Statement.Schema.Table == "audit_logs" || tx.Error != nil {
		return model.AuditLog{}
	}

	recordMap, err := getDataBeforeOperation(tx)
	if err != nil {
		return model.AuditLog{}
	}
	objId := getKeyFromData("id", recordMap)

	if objId == "" {
		return model.AuditLog{}
	}

	ctx := tx.Statement.Context
	userID := getUserIDFromContext(ctx)
	dashboardID := getDashboardIDFromContext(ctx)

	audit := model.AuditLog{
		ID:            uuid.New(),
		Table:         tx.Statement.Schema.Table,
		OperationType: model.AuditOperationDelete,
		OldValue:      null.StringFrom(prepareData(recordMap)),
		ObjectID:      null.StringFrom(objId),
		UserID:        null.StringFromPtr(userID),
		ChannelID:     null.StringFromPtr(dashboardID),
	}

	if err := tx.Session(
		&gorm.Session{
			SkipHooks: true,
			NewDB:     true,
		},
	).Create(&audit).Error; err != nil {
		c.logger.Error("error in audit log creation", slog.Any("err", err))
		return model.AuditLog{}
	}

	return audit
}

func (c *gormAuditHooks) update(tx *gorm.DB) model.AuditLog {
	if tx.Statement.Schema != nil && tx.Statement.Schema.Table == "audit_logs" || tx.Error != nil {
		return model.AuditLog{}
	}

	recordMap, err := getDataBeforeOperation(tx)
	if err != nil {
		return model.AuditLog{}
	}

	objId := getKeyFromData("id", recordMap)
	if objId == "" {
		return model.AuditLog{}
	}

	ctx := tx.Statement.Context
	userID := getUserIDFromContext(ctx)
	dashboardID := getDashboardIDFromContext(ctx)

	audit := model.AuditLog{
		ID:            uuid.New(),
		Table:         tx.Statement.Schema.Table,
		OperationType: model.AuditOperationUpdate,
		NewValue:      null.StringFrom(prepareData(recordMap)),
		ObjectID:      null.StringFrom(objId),
		UserID:        null.StringFromPtr(userID),
		ChannelID:     null.StringFromPtr(dashboardID),
	}

	if err := tx.Session(
		&gorm.Session{
			SkipHooks: true,
			NewDB:     true,
		},
	).Create(&audit).Error; err != nil {
		c.logger.Error("error in audit log creation", slog.Any("err", err))
		return model.AuditLog{}
	}

	return audit
}

func (c *gormAuditHooks) withPublisher(hook auditHook) func(tx *gorm.DB) {
	return func(tx *gorm.DB) {
		auditLog := hook(tx)

		err := c.auditLogsPubSub.Publish(
			context.Background(), auditlogs.AuditLog{
				ID:            auditLog.ID,
				Table:         auditLog.Table,
				OperationType: auditlogs.AuditOperationType(auditLog.OperationType),
				OldValue:      auditLog.OldValue,
				NewValue:      auditLog.NewValue,
				ObjectID:      auditLog.ObjectID,
				ChannelID:     auditLog.ChannelID,
				UserID:        auditLog.UserID,
				CreatedAt:     auditLog.CreatedAt,
			},
		)
		if err != nil {
			c.logger.Error(
				"failed to publish audit log",
				slog.String("audit_log_id", auditLog.ID.String()),
			)
		}
	}
}
