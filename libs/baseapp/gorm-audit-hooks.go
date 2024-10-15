package baseapp

import (
	"log/slog"

	"github.com/google/uuid"
	"github.com/guregu/null"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/logger"
	"gorm.io/gorm"
)

type gormAuditHooks struct {
	logger logger.Logger
}

func (c *gormAuditHooks) create(tx *gorm.DB) {
	if tx.Statement.Schema != nil && tx.Statement.Schema.Table == "audit_logs" || tx.Error != nil {
		return
	}

	recordMap, err := getDataBeforeOperation(tx)
	if err != nil {
		return
	}

	objId := getKeyFromData("id", recordMap)
	if objId == "" {
		return
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
	}
}

func (c *gormAuditHooks) delete(tx *gorm.DB) {
	if tx.Statement.Schema != nil && tx.Statement.Schema.Table == "audit_logs" || tx.Error != nil {
		return
	}

	recordMap, err := getDataBeforeOperation(tx)
	if err != nil {
		return
	}
	objId := getKeyFromData("id", recordMap)

	if objId == "" {
		return
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
	}
}

func (c *gormAuditHooks) update(tx *gorm.DB) {
	if tx.Statement.Schema != nil && tx.Statement.Schema.Table == "audit_logs" || tx.Error != nil {
		return
	}

	recordMap, err := getDataBeforeOperation(tx)
	if err != nil {
		return
	}

	objId := getKeyFromData("id", recordMap)
	if objId == "" {
		return
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
	}
}
