package sessions

import (
	"github.com/alexedwards/scs/goredisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/redis/go-redis/v9"
	"time"
)

func New(redisClient *redis.Client) *scs.SessionManager {
	sessionManager := scs.New()
	sessionManager.Lifetime = 24 * time.Hour * 31
	sessionManager.Store = goredisstore.New(redisClient)

	return sessionManager
}
