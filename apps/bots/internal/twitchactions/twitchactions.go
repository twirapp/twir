package twitchactions

import (
	cfg "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/logger"
	"github.com/twirapp/twir/libs/grpc/tokens"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	Logger     logger.Logger
	Config     cfg.Config
	TokensGrpc tokens.TokensClient
}

func New(opts Opts) *TwitchActions {
	return &TwitchActions{
		Logger:     opts.Logger,
		Config:     opts.Config,
		TokensGrpc: opts.TokensGrpc,
	}
}

type TwitchActions struct {
	Logger     logger.Logger
	Config     cfg.Config
	TokensGrpc tokens.TokensClient
}
