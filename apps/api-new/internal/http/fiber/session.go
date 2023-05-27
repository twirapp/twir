package fiber

import (
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/google/uuid"
	"time"
)

func NewSession(redisStorage *RedisCacheStorage) *session.Store {
	s := session.New(session.Config{
		Storage:    redisStorage,
		Expiration: 7 * 24 * time.Hour,
		KeyGenerator: func() string {
			return "fiber:session:" + uuid.New().String()
		},
	})

	return s
}
