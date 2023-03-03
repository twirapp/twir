package types

import (
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Services struct {
	Gorm   *gorm.DB
	Logger *zap.SugaredLogger
	Redis  *redis.Client
}

type WebSocketMessage struct {
	EventName string `json:"eventName"`
	Data      any    `json:"data"`
}
