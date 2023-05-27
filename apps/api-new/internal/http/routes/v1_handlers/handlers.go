package v1_handlers

import (
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/satont/tsuwari/apps/api-new/internal/http/fiber"
	"github.com/satont/tsuwari/apps/api-new/internal/http/middlewares"
	config "github.com/satont/tsuwari/libs/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In

	Logger         *zap.SugaredLogger
	Middlewares    *middlewares.Middlewares
	Config         *config.Config
	CacheStorage   *fiber.RedisCacheStorage
	SessionStorage *session.Store
	Gorm           *gorm.DB
}

type Handlers struct {
	logger         *zap.SugaredLogger
	middlewares    *middlewares.Middlewares
	config         *config.Config
	cacheStorage   *fiber.RedisCacheStorage
	sessionStorage *session.Store
	gorm           *gorm.DB
}

func NewHandlers(opts Opts) *Handlers {
	return &Handlers{
		logger:         opts.Logger,
		middlewares:    opts.Middlewares,
		config:         opts.Config,
		cacheStorage:   opts.CacheStorage,
		sessionStorage: opts.SessionStorage,
		gorm:           opts.Gorm,
	}
}
