package types

import (
	"github.com/redis/go-redis/v9"
	cfg "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/pubsub"
	"github.com/twirapp/twir/libs/grpc/bots"
	"github.com/twirapp/twir/libs/grpc/events"
	"github.com/twirapp/twir/libs/grpc/parser"
	"github.com/twirapp/twir/libs/grpc/tokens"
	"github.com/twirapp/twir/libs/grpc/websockets"
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
