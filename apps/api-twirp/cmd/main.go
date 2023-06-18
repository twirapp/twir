package main

import (
	"github.com/redis/go-redis/v9"
	"github.com/satont/tsuwari/apps/api-twirp/internal/twirp_handler"
	cfg "github.com/satont/tsuwari/libs/config"
	"go.uber.org/zap"
	"net/http"
)

func main() {
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)

	config, err := cfg.New()
	if err != nil {
		panic(err)
	}

	redisOpts, err := redis.ParseURL(config.RedisUrl)
	if err != nil {
		panic(err)
	}
	redisClient := redis.NewClient(redisOpts)

	twirpPathPrefix, twirpHandler := twirp_handler.New(redisClient)

	mux := http.NewServeMux()
	mux.Handle(twirpPathPrefix, twirpHandler)
	http.ListenAndServe(":3002", mux)
}
