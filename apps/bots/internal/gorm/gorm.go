package gorm

import (
	"context"
	"log"
	"os"
	"time"

	cfg "github.com/satont/twir/libs/config"
	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func New(config cfg.Config, lc fx.Lifecycle) (*gorm.DB, error) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             300 * time.Millisecond,
			LogLevel:                  logger.Silent,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      false,
			Colorful:                  false,
		},
	)

	db, err := gorm.Open(
		postgres.Open(config.DatabaseUrl), &gorm.Config{
			Logger: newLogger,
		},
	)
	if err != nil {
		return nil, err
	}

	d, err := db.DB()
	if err != nil {
		return nil, err
	}
	d.SetMaxIdleConns(1)
	d.SetMaxOpenConns(10)
	d.SetConnMaxLifetime(time.Hour)

	lc.Append(
		fx.Hook{
			OnStart: nil,
			OnStop: func(ctx context.Context) error {
				return d.Close()
			},
		},
	)

	return db, nil
}
