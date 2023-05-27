package v1_handlers

import (
	"github.com/satont/tsuwari/apps/api-new/internal/http/fiber"
	"github.com/satont/tsuwari/apps/api-new/internal/http/middlewares"
	config "github.com/satont/tsuwari/libs/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Opts struct {
	fx.In

	Logger       *zap.SugaredLogger
	Middlewares  *middlewares.Middlewares
	Config       *config.Config
	CacheStorage *fiber.RedisCacheStorage
}

type Handlers struct {
	logger       *zap.SugaredLogger
	middlewares  *middlewares.Middlewares
	config       *config.Config
	cacheStorage *fiber.RedisCacheStorage
}

func NewHandlers(opts Opts) *Handlers {
	return &Handlers{
		logger:       opts.Logger,
		middlewares:  opts.Middlewares,
		config:       opts.Config,
		cacheStorage: opts.CacheStorage,
	}
}
