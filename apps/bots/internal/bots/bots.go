package bots

import (
	"sync"

	cfg "github.com/satont/tsuwari/libs/config"
	"github.com/satont/tsuwari/libs/grpc/generated/parser"

	model "github.com/satont/tsuwari/libs/gomodels"

	"github.com/satont/tsuwari/apps/bots/types"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type NewBotsOpts struct {
	DB         *gorm.DB
	Logger     *zap.Logger
	Cfg        *cfg.Config
	ParserGrpc parser.ParserClient
}

type BotsService struct {
	Instances map[string]*types.BotClient
}

func NewBotsService(opts *NewBotsOpts) *BotsService {
	service := BotsService{
		Instances: make(map[string]*types.BotClient),
	}
	mu := sync.Mutex{}

	bots := []model.Bots{}
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
		instance := newBot(&ClientOpts{
			DB:         opts.DB,
			Cfg:        opts.Cfg,
			Logger:     opts.Logger,
			Model:      &b,
			ParserGrpc: opts.ParserGrpc,
		})

		mu.Lock()
		service.Instances[b.ID] = instance
		mu.Unlock()
	}

	return &service
}
