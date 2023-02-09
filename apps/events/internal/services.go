package internal

import (
	cfg "github.com/satont/tsuwari/libs/config"
	"github.com/satont/tsuwari/libs/grpc/generated/bots"
	"github.com/satont/tsuwari/libs/grpc/generated/tokens"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Services struct {
	DB     *gorm.DB
	Logger *zap.Logger
	Cfg    *cfg.Config

	BotsGrpc   bots.BotsClient
	TokensGrpc tokens.TokensClient
}
