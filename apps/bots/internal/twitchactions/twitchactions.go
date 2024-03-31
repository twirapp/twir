package twitchactions

import (
	"github.com/aidenwallis/go-ratelimiting/redis"
	adapter "github.com/aidenwallis/go-ratelimiting/redis/adapters/go-redis"
	goredis "github.com/redis/go-redis/v9"
	cfg "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/logger"
	"github.com/twirapp/twir/libs/grpc/tokens"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In

	Logger     logger.Logger
	Config     cfg.Config
	TokensGrpc tokens.TokensClient
	Gorm       *gorm.DB
	Redis      *goredis.Client
}

func New(opts Opts) *TwitchActions {
	actions := &TwitchActions{
		logger:      opts.Logger,
		config:      opts.Config,
		tokensGrpc:  opts.TokensGrpc,
		gorm:        opts.Gorm,
		rateLimiter: redis.NewSlidingWindow(adapter.NewAdapter(opts.Redis)),
	}

	return actions
}

type TwitchActions struct {
	logger      logger.Logger
	config      cfg.Config
	tokensGrpc  tokens.TokensClient
	gorm        *gorm.DB
	rateLimiter redis.SlidingWindow
}
