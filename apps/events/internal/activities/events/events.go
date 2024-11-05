package events

import (
	"github.com/redis/go-redis/v9"
	"github.com/satont/twir/apps/events/internal/hydrator"
	config "github.com/satont/twir/libs/config"
	bus_core "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/grpc/tokens"
	"github.com/twirapp/twir/libs/grpc/websockets"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In

	Gorm           *gorm.DB
	Redis          *redis.Client
	Cfg            config.Config
	TokensGrpc     tokens.TokensClient
	WebsocketsGrpc websockets.WebsocketClient
	Hydrator       *hydrator.Hydrator
	Bus            *bus_core.Bus
}

func New(opts Opts) *Activity {
	return &Activity{
		db:             opts.Gorm,
		redis:          opts.Redis,
		cfg:            opts.Cfg,
		tokensGrpc:     opts.TokensGrpc,
		websocketsGrpc: opts.WebsocketsGrpc,
		bus:            opts.Bus,
		hydrator:       opts.Hydrator,
	}
}

type Activity struct {
	db             *gorm.DB
	redis          *redis.Client
	cfg            config.Config
	tokensGrpc     tokens.TokensClient
	websocketsGrpc websockets.WebsocketClient
	hydrator       *hydrator.Hydrator
	bus            *bus_core.Bus
}
