package services

import (
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/satont/twir/apps/parser/internal/queue"
	cfg "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/grpc/generated/bots"
	"github.com/satont/twir/libs/grpc/generated/dota"
	"github.com/satont/twir/libs/grpc/generated/eval"
	"github.com/satont/twir/libs/grpc/generated/events"
	"github.com/satont/twir/libs/grpc/generated/tokens"
	"github.com/satont/twir/libs/grpc/generated/websockets"
	"github.com/satont/twir/libs/grpc/generated/ytsr"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Grpc struct {
	WebSockets websockets.WebsocketClient
	Bots       bots.BotsClient
	Dota       dota.DotaClient
	Eval       eval.EvalClient
	Tokens     tokens.TokensClient
	Events     events.EventsClient
	Ytsr       ytsr.YtsrClient
}

type Services struct {
	Config          *cfg.Config
	Logger          *zap.Logger
	Gorm            *gorm.DB
	Sqlx            *sqlx.DB
	Redis           *redis.Client
	GrpcClients     *Grpc
	TaskDistributor queue.TaskDistributor
}
