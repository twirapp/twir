package interceptors

import (
	"github.com/alexedwards/scs/v2"
	"github.com/redis/go-redis/v9"
)

type Service struct {
	redis          *redis.Client
	sessionManager *scs.SessionManager
}

func New(r *redis.Client, sessionManager *scs.SessionManager) *Service {
	return &Service{redis: r, sessionManager: sessionManager}
}
