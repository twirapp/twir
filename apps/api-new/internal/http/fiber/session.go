package fiber

import (
	"github.com/gofiber/fiber/v2/middleware/session"
	"time"
)

func NewSession(redisStorage *RedisCacheStorage) *session.Store {
	s := session.New(session.Config{
		Storage:    redisStorage,
		Expiration: 7 * 24 * time.Hour,
	})

	return s
}
