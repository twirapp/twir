package services

import (
	"github.com/go-redis/redis"
	config "github.com/satont/tsuwari/libs/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In

	Logger *zap.SugaredLogger
	Config *config.Config
	Redis  *redis.Client
	Gorm   *gorm.DB
}

type Services struct {
	Logger *zap.SugaredLogger
	Config *config.Config
	Redis  *redis.Client
	Gorm   *gorm.DB
}
