package task_queue

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	config "github.com/satont/twir/libs/config"
	"go.uber.org/zap"
)

type TaskDistributor interface {
	DistributeModUser(
		ctx context.Context,
		payload *TaskModUserPayload,
		opts ...asynq.Option,
	) error
}

type redisTaskDistributor struct {
	client *asynq.Client
	logger *zap.Logger
}

var _ TaskDistributor = (*redisTaskDistributor)(nil)

func NewRedisTaskDistributor(
	cfg *config.Config,
	logger *zap.Logger,
) *redisTaskDistributor {
	url, err := redis.ParseURL(cfg.RedisUrl)
	if err != nil {
		panic("Wrong redis url")
	}

	redisOpt := asynq.RedisClientOpt{
		Addr:     url.Addr,
		Password: url.Password,
		DB:       url.DB,
		Username: url.Username,
		PoolSize: url.PoolSize,
	}

	client := asynq.NewClient(redisOpt)

	distributor := &redisTaskDistributor{
		client: client,
		logger: logger,
	}

	return distributor
}
