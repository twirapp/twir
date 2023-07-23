package interceptors

import (
	"github.com/alexedwards/scs/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Service struct {
	redis          *redis.Client
	sessionManager *scs.SessionManager
	db             *gorm.DB
}

func New(r *redis.Client, sessionManager *scs.SessionManager, db *gorm.DB) *Service {
	return &Service{redis: r, sessionManager: sessionManager, db: db}
}
