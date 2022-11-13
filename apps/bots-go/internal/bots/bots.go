package bots

import (
	"sync"
	cfg "tsuwari/config"
	model "tsuwari/models"
	"tsuwari/twitch"

	irc "github.com/gempir/go-twitch-irc/v3"
	"github.com/nats-io/nats.go"
	"github.com/satont/tsuwari/apps/bots/types"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type NewBotsOpts struct {
	Twitch *twitch.Twitch
	DB     *gorm.DB
	Logger *zap.Logger
	Cfg    *cfg.Config
	Nats   *nats.Conn
}

type botInstance struct {
	bot *irc.Client
	api *twitch.Twitch
}

type BotsService struct {
	Instances []*types.BotClient
}

func NewBotsService(opts *NewBotsOpts) {
	service := BotsService{}
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
		go func(bot model.Bots) {
			instance := newBot(&ClientOpts{
				DB:     opts.DB,
				Cfg:    opts.Cfg,
				Logger: opts.Logger,
				Bot:    &bot,
				Nats:   opts.Nats,
			})
			mu.Lock()
			service.Instances = append(service.Instances, instance)
			mu.Unlock()
		}(bot)
	}
}
