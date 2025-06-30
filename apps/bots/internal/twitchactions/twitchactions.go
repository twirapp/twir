package twitchactions

import (
	"github.com/aidenwallis/go-ratelimiting/redis"
	adapter "github.com/aidenwallis/go-ratelimiting/redis/adapters/go-redis"
	goredis "github.com/redis/go-redis/v9"
	mod_task_queue "github.com/satont/twir/apps/bots/internal/mod-task-queue"
	toxicity_check "github.com/satont/twir/apps/bots/internal/services/toxicity-check"
	cfg "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/logger"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	"github.com/twirapp/twir/libs/grpc/tokens"
	"github.com/twirapp/twir/libs/repositories/channels"
	channelmodel "github.com/twirapp/twir/libs/repositories/channels/model"
	"github.com/twirapp/twir/libs/repositories/sentmessages"
	"github.com/twirapp/twir/libs/repositories/toxic_messages"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In

	Logger                  logger.Logger
	TokensGrpc              tokens.TokensClient
	ModTaskDistributor      mod_task_queue.TaskDistributor
	SentMessagesRepository  sentmessages.Repository
	ChannelsRepository      channels.Repository
	ToxicMessagesRepository toxic_messages.Repository
	Gorm                    *gorm.DB
	Redis                   *goredis.Client
	ToxicityCheck           *toxicity_check.Service
	Config                  cfg.Config
	ChannelsCache           *generic_cacher.GenericCacher[channelmodel.Channel]
}

func New(opts Opts) *TwitchActions {
	actions := &TwitchActions{
		logger:                  opts.Logger,
		config:                  opts.Config,
		tokensGrpc:              opts.TokensGrpc,
		gorm:                    opts.Gorm,
		rateLimiter:             redis.NewSlidingWindow(adapter.NewAdapter(opts.Redis)),
		modTaskDistributor:      opts.ModTaskDistributor,
		sentMessagesRepository:  opts.SentMessagesRepository,
		channelsRepository:      opts.ChannelsRepository,
		toxicityCheck:           opts.ToxicityCheck,
		toxicMessagesRepository: opts.ToxicMessagesRepository,
		channelsCache:           opts.ChannelsCache,
	}

	return actions
}

type TwitchActions struct {
	logger                  logger.Logger
	tokensGrpc              tokens.TokensClient
	rateLimiter             redis.SlidingWindow
	modTaskDistributor      mod_task_queue.TaskDistributor
	sentMessagesRepository  sentmessages.Repository
	channelsRepository      channels.Repository
	toxicMessagesRepository toxic_messages.Repository
	gorm                    *gorm.DB
	toxicityCheck           *toxicity_check.Service
	config                  cfg.Config
	channelsCache           *generic_cacher.GenericCacher[channelmodel.Channel]
}
