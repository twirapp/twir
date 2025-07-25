package interceptors

import (
	"github.com/alexedwards/scs/v2"
	"github.com/redis/go-redis/v9"
	"github.com/twirapp/twir/libs/logger"
	"gorm.io/gorm"
)

type Service struct {
	redis          *redis.Client
	sessionManager *scs.SessionManager
	db             *gorm.DB
	logger         logger.Logger
}

func New(
	r *redis.Client,
	sessionManager *scs.SessionManager,
	db *gorm.DB,
	l logger.Logger,
) *Service {
	return &Service{redis: r, sessionManager: sessionManager, db: db, logger: l}
}
