package auth

import (
	config "github.com/satont/twir/libs/config"
	sessions "github.com/twirapp/twir/apps/api-gql/internal/auth"
	"github.com/twirapp/twir/apps/api-gql/internal/httpserver"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/grpc/tokens"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In

	Server     *httpserver.Server
	Gorm       *gorm.DB
	Config     config.Config
	TokensGrpc tokens.TokensClient
	Bus        *buscore.Bus
	Sessions   *sessions.Auth
}

type Auth struct {
	gorm       *gorm.DB
	config     config.Config
	tokensGrpc tokens.TokensClient
	bus        *buscore.Bus
	sessions   *sessions.Auth
}

func New(opts Opts) *Auth {
	p := &Auth{
		gorm:       opts.Gorm,
		config:     opts.Config,
		tokensGrpc: opts.TokensGrpc,
		bus:        opts.Bus,
		sessions:   opts.Sessions,
	}

	opts.Server.POST("/auth", p.handleAuthPostCode)

	return p
}
