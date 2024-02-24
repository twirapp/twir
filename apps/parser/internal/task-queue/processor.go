package task_queue

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	config "github.com/satont/twir/libs/config"
	"github.com/twirapp/twir/libs/grpc/tokens"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

//goland:noinspection GoNameStartsWithPackageName
const (
	QueueDefault = "default"
)

const (
	TaskModUser = "task:mod_user"
)

type TaskProcessor interface {
	Start() error
	Stop() error

	ProcessDistributeMod(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	config     config.Config
	server     *asynq.Server
	logger     *zap.Logger
	gorm       *gorm.DB
	tokensGrpc tokens.TokensClient
}

var _ TaskProcessor = (*RedisTaskProcessor)(nil)

type RedisTaskProcessorOpts struct {
	Cfg        config.Config
	Logger     *zap.Logger
	Gorm       *gorm.DB
	TokensGrpc tokens.TokensClient
}

func NewRedisTaskProcessor(opts RedisTaskProcessorOpts) *RedisTaskProcessor {
	url, err := redis.ParseURL(opts.Cfg.RedisUrl)
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

	server := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Queues: map[string]int{
				QueueDefault: 5,
			},
			ErrorHandler: asynq.ErrorHandlerFunc(
				func(ctx context.Context, task *asynq.Task, err error) {
					opts.Logger.Error("error processing task", zap.Any("task", task), zap.Error(err))
				},
			),
		},
	)

	processor := &RedisTaskProcessor{
		config:     opts.Cfg,
		server:     server,
		logger:     opts.Logger,
		gorm:       opts.Gorm,
		tokensGrpc: opts.TokensGrpc,
	}

	return processor
}

func (p *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(TaskModUser, p.ProcessDistributeMod)

	return p.server.Start(mux)
}

func (p *RedisTaskProcessor) Stop() error {
	p.server.Stop()
	p.server.Shutdown()

	return nil
}
