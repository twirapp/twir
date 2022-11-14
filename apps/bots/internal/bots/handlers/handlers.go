package handlers

import (
	cfg "github.com/satont/tsuwari/libs/config"

	model "github.com/satont/tsuwari/libs/gomodels"

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
	cfg       *cfg.Config
}

func CreateHandlers(opts *HandlersOpts) *Handlers {
	handlersService := &Handlers{
		db:        opts.DB,
		logger:    opts.Logger,
		BotClient: opts.BotClient,
		nats:      opts.Nats,
		cfg:       opts.Cfg,
	}

	return handlersService
}
