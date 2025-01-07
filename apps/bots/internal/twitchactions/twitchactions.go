package twitchactions

import (
	"github.com/aidenwallis/go-ratelimiting/redis"
	adapter "github.com/aidenwallis/go-ratelimiting/redis/adapters/go-redis"
	goredis "github.com/redis/go-redis/v9"
	mod_task_queue "github.com/satont/twir/apps/bots/internal/mod-task-queue"
	cfg "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/logger"
	"github.com/twirapp/twir/libs/grpc/tokens"
	"github.com/twirapp/twir/libs/repositories/sentmessages"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In

	Logger                 logger.Logger
	Config                 cfg.Config
	TokensGrpc             tokens.TokensClient
	Gorm                   *gorm.DB
	Redis                  *goredis.Client
	ModTaskDistributor     mod_task_queue.TaskDistributor
	SentMessagesRepository sentmessages.Repository
}

func New(opts Opts) *TwitchActions {
	actions := &TwitchActions{
		logger:                 opts.Logger,
		config:                 opts.Config,
		tokensGrpc:             opts.TokensGrpc,
		gorm:                   opts.Gorm,
		rateLimiter:            redis.NewSlidingWindow(adapter.NewAdapter(opts.Redis)),
		modTaskDistributor:     opts.ModTaskDistributor,
		sentMessagesRepository: opts.SentMessagesRepository,
	}

	return actions
}

type TwitchActions struct {
	logger                 logger.Logger
	config                 cfg.Config
	tokensGrpc             tokens.TokensClient
	gorm                   *gorm.DB
	rateLimiter            redis.SlidingWindow
	modTaskDistributor     mod_task_queue.TaskDistributor
	sentMessagesRepository sentmessages.Repository
}
