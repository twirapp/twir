package types

import (
	"github.com/satont/twir/apps/scheduler/internal/services"
	config "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/grpc/generated/emotes_cacher"
	"github.com/satont/twir/libs/grpc/generated/parser"
	"github.com/satont/twir/libs/grpc/generated/tokens"
	"github.com/satont/twir/libs/grpc/generated/watched"
	"github.com/satont/twir/libs/pubsub"
	"gorm.io/gorm"
)

type GrpcServices struct {
	Emotes  emotes_cacher.EmotesCacherClient
	Parser  parser.ParserClient
	Tokens  tokens.TokensClient
	Watched watched.WatchedClient
}

type Services struct {
	Grpc     *GrpcServices
	Gorm     *gorm.DB
	Config   *config.Config
	PubSub   *pubsub.PubSub
	Commands *services.Commands
	Roles    *services.Roles
}
