package redis

import (
	"context"
	"tsuwari/parser/internal/config/cfg"

	"github.com/go-redis/redis/v9"
)

var (
	RdbCtx  = context.Background()
	Rdb redis.Client
)

func Connect() {
	url, err := redis.ParseURL(cfg.Cfg.RedisUrl)

	if err != nil {
		panic("Wrong redis url")
	}

	client := redis.NewClient(&redis.Options{
		Addr:     url.Addr,
		Password: url.Password,
		DB:       url.DB,
		Username: url.Username,
	})

	Rdb = *client
}