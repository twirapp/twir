package baseapp

import (
	"context"
	"log"
	"os"
	"time"

	config "github.com/satont/twir/libs/config"
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
		
		auditHooks := &gormAuditHooks{
			l,
		}

		if withAudit {
			db.Callback().Create().After("gorm:create").Register(
				"custom_plugin:create_audit_log",
				auditHooks.create,
			)
			db.Callback().Update().After("gorm:update").Register(
				"custom_plugin:update_audit_log",
				auditHooks.update,
			)
			db.Callback().Delete().Before("gorm:delete").Register(
				"custom_plugin:delete_audit_log",
				auditHooks.delete,
			)
		}

		return db, nil
	}
}
