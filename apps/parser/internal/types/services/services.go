package services

import (
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/satont/twir/apps/parser/internal/task-queue"
	cfg "github.com/satont/twir/libs/config"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/grpc/dota"
	"github.com/twirapp/twir/libs/grpc/eval"
	"github.com/twirapp/twir/libs/grpc/events"
	"github.com/twirapp/twir/libs/grpc/tokens"
	"github.com/twirapp/twir/libs/grpc/websockets"
	"github.com/twirapp/twir/libs/grpc/ytsr"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Grpc struct {
	WebSockets websockets.WebsocketClient
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
	TaskDistributor task_queue.TaskDistributor
	Bus             *buscore.Bus
}
