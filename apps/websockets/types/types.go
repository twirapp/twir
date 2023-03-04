package types

import (
	"github.com/redis/go-redis/v9"
	"github.com/satont/tsuwari/libs/grpc/generated/bots"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type GrpcClients struct {
	Bots bots.BotsClient
}

type Services struct {
	Gorm   *gorm.DB
	Logger *zap.SugaredLogger
	Redis  *redis.Client
	Grpc   *GrpcClients
}

type WebSocketMessage struct {
	EventName string `json:"eventName"`
	Data      any    `json:"data"`
}
