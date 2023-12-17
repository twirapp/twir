package events

import (
	"github.com/redis/go-redis/v9"
	"github.com/satont/twir/apps/events/internal/hydrator"
	config "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/grpc/generated/bots"
	"github.com/satont/twir/libs/grpc/generated/tokens"
	"github.com/satont/twir/libs/grpc/generated/websockets"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In

	Gorm           *gorm.DB
	Redis          *redis.Client
	Cfg            config.Config
	TokensGrpc     tokens.TokensClient
	BotsGrpc       bots.BotsClient
	WebsocketsGrpc websockets.WebsocketClient
	Hydrator       *hydrator.Hydrador
}

func New(opts Opts) *Activity {
	return &Activity{
		db:             opts.Gorm,
		redis:          opts.Redis,
		cfg:            opts.Cfg,
		tokensGrpc:     opts.TokensGrpc,
		botsGrpc:       opts.BotsGrpc,
		websocketsGrpc: opts.WebsocketsGrpc,
	}
}

type Activity struct {
	db             *gorm.DB
	redis          *redis.Client
	cfg            config.Config
	tokensGrpc     tokens.TokensClient
	botsGrpc       bots.BotsClient
	websocketsGrpc websockets.WebsocketClient
	hydrator       *hydrator.Hydrador
}
