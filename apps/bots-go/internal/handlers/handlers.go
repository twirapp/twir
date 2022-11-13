package handlers

import (
	cfg "tsuwari/config"
	model "tsuwari/models"

	"github.com/nats-io/nats.go"
	"github.com/satont/tsuwari/apps/bots/types"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type BotInstance struct {
	BotClient *types.BotClient
	Db        *model.Bots
}

type HandlersOpts struct {
	DB        *gorm.DB
	Logger    *zap.Logger
	Cfg       *cfg.Config
	BotClient *types.BotClient
	Nats      *nats.Conn
}

type Handlers struct {
	db        *gorm.DB
	logger    *zap.Logger
	BotClient *types.BotClient
	nats      *nats.Conn
}

func CreateHandlers(opts *HandlersOpts) *Handlers {
	return &Handlers{
		db:        opts.DB,
		logger:    opts.Logger,
		BotClient: opts.BotClient,
		nats:      opts.Nats,
	}
}
