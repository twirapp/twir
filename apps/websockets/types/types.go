package types

import (
	"github.com/redis/go-redis/v9"
	"github.com/satont/twir/libs/grpc/generated/bots"
	"github.com/satont/twir/libs/grpc/generated/parser"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type GrpcClients struct {
	Bots   bots.BotsClient
	Parser parser.ParserClient
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
