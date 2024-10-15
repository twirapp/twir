package baseapp

import (
	"context"
	"database/sql"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/google/uuid"
	config "github.com/satont/twir/libs/config"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/logger"
	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type selectedDashboardContextKeyType string
type requesterUserIdContextKeyType string

const (
	SelectedDashboardContextKey = selectedDashboardContextKeyType("__selectedDashboard__")
	RequesterUserIdContextKey   = requesterUserIdContextKeyType("__requesterUserId__")
)

func newGorm(withAudit bool) func(
	cfg config.Config,
	l logger.Logger,
	lc fx.Lifecycle,
) (*gorm.DB, error) {
	return func(cfg config.Config, l logger.Logger, lc fx.Lifecycle) (*gorm.DB, error) {
		newLogger := gormlogger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			gormlogger.Config{
				SlowThreshold:             100 * time.Millisecond,
				LogLevel:                  gormlogger.Warn,
				IgnoreRecordNotFoundError: true,
				ParameterizedQueries:      true,
				Colorful:                  true,
			},
		)

		db, err := gorm.Open(
			postgres.Open(cfg.DatabaseUrl),
			&gorm.Config{
				Logger: newLogger,
			},
		)
		if err != nil {
			return nil, err
		}
		d, _ := db.DB()
		d.SetMaxIdleConns(1)
		d.SetMaxOpenConns(10)
		d.SetConnMaxLifetime(time.Hour)

		lc.Append(
			fx.Hook{
				OnStop: func(_ context.Context) error {
					return d.Close()
				},
			},
		)

		h := &auditHooks{
			l,
		}

		if withAudit {
			db.Callback().Create().After("gorm:create").Register(
				"custom_plugin:create_audit_log",
				h.create,
			)
			// db.Callback().Update().After("gorm:update").Register(
			// 	"custom_plugin:update_audit_log",
			// 	updateAuditLog,
			// )
			db.Callback().Delete().Before("gorm:delete").Register(
				"custom_plugin:delete_audit_log",
				h.delete,
			)
		}

		return db, nil
	}
}

type auditHooks struct {
	logger logger.Logger
}

func (c *auditHooks) create(tx *gorm.DB) {
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
		NewValue: sql.Null[string]{
			V:     prepareData(recordMap),
			Valid: true,
		},
		ObjectID: sql.Null[string]{
			V:     objId,
			Valid: true,
		},
		UserID: sql.Null[string]{
			V:     "",
			Valid: false,
		},
		DashboardID: sql.Null[string]{},
	}

	if userID != nil {
		audit.UserID = sql.Null[string]{
			V:     *userID,
			Valid: true,
		}
	}

	if dashboardID != nil {
		audit.DashboardID = sql.Null[string]{
			V:     *dashboardID,
			Valid: true,
		}
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

func (c *auditHooks) delete(tx *gorm.DB) {
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
		OldValue: sql.Null[string]{
			V:     prepareData(recordMap),
			Valid: true,
		},
		NewValue: sql.Null[string]{},
		ObjectID: sql.Null[string]{
			V:     objId,
			Valid: true,
		},
		UserID:      sql.Null[string]{},
		DashboardID: sql.Null[string]{},
	}

	if userID != nil {
		audit.UserID = sql.Null[string]{
			V:     *userID,
			Valid: true,
		}
	}

	if dashboardID != nil {
		audit.DashboardID = sql.Null[string]{
			V:     *dashboardID,
			Valid: true,
		}
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
