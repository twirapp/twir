package nats_handlers

import (
	"github.com/satont/tsuwari/apps/bots/internal/bots"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type NatsHandlers struct {
	db          *gorm.DB
	botsService *bots.BotsService
	logger      *zap.Logger
}

type NatsHandlersOpts struct {
	Db          *gorm.DB
	BotsService *bots.BotsService
	Logger      *zap.Logger
}

func NewNatsHandlers(opts *NatsHandlersOpts) *NatsHandlers {
	return &NatsHandlers{}
}
