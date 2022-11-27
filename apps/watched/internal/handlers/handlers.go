package handlers

import (
	cfg "github.com/satont/tsuwari/libs/config"
	"gorm.io/gorm"
)

type HandlersOpts struct {
	DB  *gorm.DB
	Cfg *cfg.Config
}

type Handlers struct {
	db  *gorm.DB
	cfg *cfg.Config
}

func NewHandlers(opts HandlersOpts) *Handlers {
	return &Handlers{
		db:  opts.DB,
		cfg: opts.Cfg,
	}
}
