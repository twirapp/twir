package types

import (
	"github.com/redis/go-redis/v9"
	cfg "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/grpc/generated/bots"
	"github.com/satont/twir/libs/grpc/generated/events"
	"github.com/satont/twir/libs/grpc/generated/parser"
	"github.com/satont/twir/libs/grpc/generated/tokens"
	"github.com/satont/twir/libs/grpc/generated/websockets"
	"github.com/satont/twir/libs/pubsub"
	"gorm.io/gorm"
)

type GrpcClients struct {
	Tokens     tokens.TokensClient
	Events     events.EventsClient
	Bots       bots.BotsClient
	Parser     parser.ParserClient
	WebSockets websockets.WebsocketClient
}

type Services struct {
	Gorm   *gorm.DB
	Config *cfg.Config
	Grpc   *GrpcClients
	PubSub *pubsub.PubSub
	Redis  *redis.Client
}
