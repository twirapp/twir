package mod_task_queue

import (
	"context"
	"log/slog"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"github.com/twirapp/twir/libs/baseapp/lifecycle"
	buscore "github.com/twirapp/twir/libs/bus-core"
	config "github.com/twirapp/twir/libs/config"
	loggerlib "github.com/twirapp/twir/libs/logger"
	"gorm.io/gorm"
)

//goland:noinspection GoNameStartsWithPackageName
const (
	QueueDefault = "default"
)

const (
	TaskModUser = "bots:task:mod_user"
)

type TaskProcessor interface {
	Start() error
	Stop() error

	ProcessDistributeMod(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	config  config.Config
	server  *asynq.Server
	logger  *slog.Logger
	gorm    *gorm.DB
	twirBus *buscore.Bus
}

var _ TaskProcessor = (*RedisTaskProcessor)(nil)

func NewRedisTaskProcessor(
	lc *lifecycle.Lifecycle,
	cfg config.Config,
	logger *slog.Logger,
	gormDB *gorm.DB,
	twirBus *buscore.Bus,
) *RedisTaskProcessor {
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

	server := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Queues: map[string]int{
				QueueDefault: 5,
			},
			ErrorHandler: asynq.ErrorHandlerFunc(
				func(ctx context.Context, task *asynq.Task, err error) {
					logger.Error("error processing task", slog.Any("task", task), loggerlib.Error(err))
				},
			),
			LogLevel: asynq.ErrorLevel,
		},
	)

	processor := &RedisTaskProcessor{
		config:  cfg,
		server:  server,
		logger:  logger,
		gorm:    gormDB,
		twirBus: twirBus,
	}

	lc.Append(
		lifecycle.Hook{
			OnStart: func(ctx context.Context) error {
				go func() {
					logger.Info("Starting mod task processor")
					if err := processor.Start(); err != nil {
						panic(err)
					}
				}()

				return nil
			},
			OnStop: func(ctx context.Context) error {
				return processor.Stop()
			},
		},
	)

	return processor
}

func (p *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(TaskModUser, p.ProcessDistributeMod)

	p.logger.Info("Registered task handler", slog.String("task", TaskModUser))
	return p.server.Start(mux)
}

func (p *RedisTaskProcessor) Stop() error {
	p.server.Stop()
	p.server.Shutdown()

	return nil
}
