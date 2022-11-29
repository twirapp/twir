package handlers

import (
	cfg "github.com/satont/tsuwari/libs/config"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type HandlersOpts struct {
	DB     *gorm.DB
	Cfg    *cfg.Config
	Logger *zap.Logger
}

type Handlers struct {
	db     *gorm.DB
	cfg    *cfg.Config
	logger *zap.Logger
}

func NewHandlers(opts HandlersOpts) *Handlers {
	return &Handlers{
		db:     opts.DB,
		cfg:    opts.Cfg,
		logger: opts.Logger,
	}
}
