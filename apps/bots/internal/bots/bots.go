package bots

import (
	"sync"

	cfg "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/grpc/generated/parser"

	model "github.com/satont/twir/libs/gomodels"

	"github.com/satont/twir/apps/bots/types"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type NewBotsOpts struct {
	DB         *gorm.DB
	Logger     *zap.Logger
	Cfg        *cfg.Config
	ParserGrpc parser.ParserClient
}

type Service struct {
	Instances map[string]*types.BotClient
}

func NewBotsService(opts *NewBotsOpts) *Service {
	service := Service{
		Instances: make(map[string]*types.BotClient),
	}
	mu := sync.Mutex{}

	var bots []model.Bots
	err := opts.DB.
		Preload("Token").
		Preload("Channels").
		Find(&bots).
		Error
	if err != nil {
		panic(err)
	}

	for _, bot := range bots {
		b := bot
		instance := newBot(
			&ClientOpts{
				DB:         opts.DB,
				Cfg:        opts.Cfg,
				Logger:     opts.Logger,
				Model:      &b,
				ParserGrpc: opts.ParserGrpc,
			},
		)

		mu.Lock()
		service.Instances[b.ID] = instance
		mu.Unlock()
	}

	return &service
}
