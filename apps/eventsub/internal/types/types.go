package types

import (
	cfg "github.com/satont/tsuwari/libs/config"
	"github.com/satont/tsuwari/libs/grpc/generated/tokens"
	"gorm.io/gorm"
)

type GrpcClients struct {
	Tokens tokens.TokensClient
}

type Services struct {
	Gorm   *gorm.DB
	Config *cfg.Config
	Grpc   *GrpcClients
}
