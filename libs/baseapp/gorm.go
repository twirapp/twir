package baseapp

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/twirapp/twir/libs/baseapp/lifecycle"
	config "github.com/twirapp/twir/libs/config"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
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

func newGorm(
	cfg config.Config,
	lc *lifecycle.Lifecycle,
	pool *pgxpool.Pool,
) (*gorm.DB, error) {
	return createGorm(cfg, pool, func(onStop func(context.Context) error) {
		lc.Append(lifecycle.Hook{OnStop: onStop})
	})
}

func createGorm(
	cfg config.Config,
	pool *pgxpool.Pool,
	appendStop func(func(context.Context) error),
) (*gorm.DB, error) {
	newLogger := gormlogger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		gormlogger.Config{
			SlowThreshold:             100 * time.Millisecond,
			LogLevel:                  gormlogger.Error,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
			Colorful:                  true,
		},
	)

	sqlDb := stdlib.OpenDBFromPool(pool)

	dialector := postgres.New(
		postgres.Config{
			Conn: sqlDb,
		},
	)

	db, err := gorm.Open(
		dialector,
		&gorm.Config{
			Logger:                 newLogger,
			SkipDefaultTransaction: true,
			// PrepareStmt:            false,
		},
	)
	if err != nil {
		return nil, err
	}
	d, _ := db.DB()
	d.SetMaxIdleConns(1)
	d.SetMaxOpenConns(100)
	d.SetConnMaxLifetime(5 * time.Minute)
	d.SetConnMaxIdleTime(1 * time.Minute)

	if err := db.Use(otelgorm.NewPlugin()); err != nil {
		return nil, err
	}

	appendStop(func(_ context.Context) error {
		return d.Close()
	})

	return db, nil
}
