package interceptors

import "github.com/redis/go-redis/v9"

type Service struct {
	redis *redis.Client
}

func New(r *redis.Client) *Service {
	return &Service{redis: r}
}
