package mod_task_queue

import (
	"context"
	"log/slog"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	buscore "github.com/twirapp/twir/libs/bus-core"
	config "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/logger"
	"go.uber.org/fx"
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

type RedisTaskProcessorOpts struct {
	fx.In
	LC fx.Lifecycle

	Cfg     config.Config
	Logger  *slog.Logger
	Gorm    *gorm.DB
	TwirBus *buscore.Bus
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
					opts.Logger.Error("error processing task", slog.Any("task", task), logger.Error(err))
				},
			),
			LogLevel: asynq.ErrorLevel,
		},
	)

	processor := &RedisTaskProcessor{
		config:  opts.Cfg,
		server:  server,
		logger:  opts.Logger,
		gorm:    opts.Gorm,
		twirBus: opts.TwirBus,
	}

	opts.LC.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				go func() {
					opts.Logger.Info("Starting mod task processor")
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
