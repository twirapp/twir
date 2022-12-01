package handlers

import (
	cfg "github.com/satont/tsuwari/libs/config"
	"github.com/satont/tsuwari/libs/grpc/generated/parser"

	model "github.com/satont/tsuwari/libs/gomodels"

	"github.com/satont/tsuwari/apps/bots/types"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type BotInstance struct {
	BotClient *types.BotClient
	Db        *model.Bots
}

type HandlersOpts struct {
	DB         *gorm.DB
	Logger     *zap.Logger
	Cfg        *cfg.Config
	BotClient  *types.BotClient
	ParserGrpc parser.ParserClient
}

type Handlers struct {
	db         *gorm.DB
	logger     *zap.Logger
	BotClient  *types.BotClient
	cfg        *cfg.Config
	parserGrpc parser.ParserClient
}

func CreateHandlers(opts *HandlersOpts) *Handlers {
	handlersService := &Handlers{
		db:         opts.DB,
		logger:     opts.Logger,
		BotClient:  opts.BotClient,
		cfg:        opts.Cfg,
		parserGrpc: opts.ParserGrpc,
	}

	return handlersService
}
