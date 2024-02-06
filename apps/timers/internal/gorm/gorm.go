package gorm

import (
	"context"
	"time"

	cfg "github.com/satont/twir/libs/config"
	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func New(config cfg.Config, lc fx.Lifecycle) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(config.DatabaseUrl))
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
