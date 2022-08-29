package redis

import (
	"github.com/go-redis/redis/v9"
)

func New(redisUrl string) *redis.Client {
	url, err := redis.ParseURL(redisUrl)

	if err != nil {
		panic("Wrong redis url")
	}

	client := redis.NewClient(&redis.Options{
		Addr:     url.Addr,
		Password: url.Password,
		DB:       url.DB,
		Username: url.Username,
	})

	client.Conn()

	return client
}
