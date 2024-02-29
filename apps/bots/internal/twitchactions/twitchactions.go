package twitchactions

import (
	cfg "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/logger"
	"github.com/twirapp/twir/libs/grpc/tokens"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In

	Logger     logger.Logger
	Config     cfg.Config
	TokensGrpc tokens.TokensClient
	Gorm       *gorm.DB
}

func New(opts Opts) *TwitchActions {
	return &TwitchActions{
		logger:     opts.Logger,
		config:     opts.Config,
		tokensGrpc: opts.TokensGrpc,
		gorm:       opts.Gorm,
	}
}

type TwitchActions struct {
	logger     logger.Logger
	config     cfg.Config
	tokensGrpc tokens.TokensClient
	gorm       *gorm.DB
}
