package mod_task_queue

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	config "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/logger"
)

type TaskDistributor interface {
	DistributeModUser(
		ctx context.Context,
		payload *TaskModUserPayload,
		opts ...asynq.Option,
	) error
}

type ModTaskDistributor struct {
	client *asynq.Client
	logger logger.Logger
}

var _ TaskDistributor = (*ModTaskDistributor)(nil)

func NewRedisModTaskDistributor(
	cfg config.Config,
	logger logger.Logger,
) (*ModTaskDistributor, error) {
	url, err := redis.ParseURL(cfg.RedisUrl)
	if err != nil {
		return nil, err
	}

	redisOpt := asynq.RedisClientOpt{
		Addr:     url.Addr,
		Password: url.Password,
		DB:       url.DB,
		Username: url.Username,
		PoolSize: url.PoolSize,
	}

	client := asynq.NewClient(redisOpt)

	distributor := &ModTaskDistributor{
		client: client,
		logger: logger,
	}

	return distributor, nil
}
