package types

import (
	cfg "github.com/satont/tsuwari/libs/config"
	"github.com/satont/tsuwari/libs/grpc/generated/bots"
	"github.com/satont/tsuwari/libs/grpc/generated/events"
	"github.com/satont/tsuwari/libs/grpc/generated/parser"
	"github.com/satont/tsuwari/libs/grpc/generated/tokens"
	"github.com/satont/tsuwari/libs/pubsub"
	"gorm.io/gorm"
)

type GrpcClients struct {
	Tokens tokens.TokensClient
	Events events.EventsClient
	Bots   bots.BotsClient
	Parser parser.ParserClient
}

type Services struct {
	Gorm   *gorm.DB
	Config *cfg.Config
	Grpc   *GrpcClients
	PubSub *pubsub.PubSub
}
