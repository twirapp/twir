package redis

import (
	"context"
	"tsuwari/parser/internal/config/cfg"

	"github.com/go-redis/redis/v9"
	"github.com/rueian/rueidis"
)

var RedisCtx = context.Background()
var Rdb rueidis.Client

func Connect() {
	if cfg.Cfg.RedisUrl == nil {
		panic("No Redis url setuped")
	}

	url, err := redis.ParseURL(*cfg.Cfg.RedisUrl)

	if err != nil {
		panic("Wrong redis url")
	}

	Rdb, err = rueidis.NewClient(rueidis.ClientOption{
		InitAddress: []string{url.Addr},
		Username: url.Username,
		SelectDB: url.DB,
		Password: url.Password,
	})

	if err != nil {
		panic(err)
	}
}