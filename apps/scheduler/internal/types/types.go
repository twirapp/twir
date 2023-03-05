package types

import (
	config "github.com/satont/tsuwari/libs/config"
	"github.com/satont/tsuwari/libs/grpc/generated/emotes_cacher"
	"github.com/satont/tsuwari/libs/grpc/generated/parser"
	"github.com/satont/tsuwari/libs/grpc/generated/tokens"
	"github.com/satont/tsuwari/libs/grpc/generated/watched"
	"github.com/satont/tsuwari/libs/pubsub"
	"gorm.io/gorm"
)

type GrpcServices struct {
	Emotes  emotes_cacher.EmotesCacherClient
	Parser  parser.ParserClient
	Tokens  tokens.TokensClient
	Watched watched.WatchedClient
}

type Services struct {
	Grpc   *GrpcServices
	Gorm   *gorm.DB
	Config *config.Config
	PubSub *pubsub.PubSub
}
