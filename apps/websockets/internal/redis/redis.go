package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
	cfg "github.com/satont/twir/libs/config"
	"go.uber.org/fx"
)

func New(config cfg.Config, lc fx.Lifecycle) (*redis.Client, error) {
	opts, err := redis.ParseURL(config.RedisUrl)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opts)

	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				return client.Ping(ctx).Err()
			},
			OnStop: func(ctx context.Context) error {
				return client.Close()
			},
		},
	)

	return client, nil
}
