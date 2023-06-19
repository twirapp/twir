package impl_deps

import (
	"github.com/alexedwards/scs/v2"
	"github.com/redis/go-redis/v9"
	cfg "github.com/satont/tsuwari/libs/config"
	"github.com/satont/tsuwari/libs/grpc/generated/bots"
	"github.com/satont/tsuwari/libs/grpc/generated/integrations"
	"github.com/satont/tsuwari/libs/grpc/generated/tokens"
	"gorm.io/gorm"
)

type Grpc struct {
	Tokens       tokens.TokensClient
	Integrations integrations.IntegrationsClient
	Bots         bots.BotsClient
}

type Deps struct {
	Config         *cfg.Config
	Redis          *redis.Client
	Db             *gorm.DB
	Grpc           *Grpc
	SessionManager *scs.SessionManager
}
