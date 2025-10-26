package twitchactions

import (
	"github.com/aidenwallis/go-ratelimiting/redis"
	adapter "github.com/aidenwallis/go-ratelimiting/redis/adapters/go-redis"
	goredis "github.com/redis/go-redis/v9"
	"github.com/twirapp/kv"
	mod_task_queue "github.com/twirapp/twir/apps/bots/internal/mod-task-queue"
	toxicity_check "github.com/twirapp/twir/apps/bots/internal/services/toxicity-check"
	buscore "github.com/twirapp/twir/libs/bus-core"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	cfg "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/logger"
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
	SentMessagesRepository  sentmessages.Repository
	ChannelsRepository      channels.Repository
	ToxicMessagesRepository toxic_messages.Repository
	Gorm                    *gorm.DB
	Redis                   *goredis.Client
	ToxicityCheck           *toxicity_check.Service
	Config                  cfg.Config
	ChannelsCache           *generic_cacher.GenericCacher[channelmodel.Channel]
	TwirBus                 *buscore.Bus
	KV                      kv.KV
	ModTaskDistributor      mod_task_queue.TaskDistributor
}

func New(opts Opts) *TwitchActions {
	actions := &TwitchActions{
		logger:                  opts.Logger,
		config:                  opts.Config,
		twirBus:                 opts.TwirBus,
		gorm:                    opts.Gorm,
		rateLimiter:             redis.NewSlidingWindow(adapter.NewAdapter(opts.Redis)),
		sentMessagesRepository:  opts.SentMessagesRepository,
		channelsRepository:      opts.ChannelsRepository,
		toxicityCheck:           opts.ToxicityCheck,
		toxicMessagesRepository: opts.ToxicMessagesRepository,
		channelsCache:           opts.ChannelsCache,
		kv:                      opts.KV,
		modTaskDistributor:      opts.ModTaskDistributor,
	}

	return actions
}

type TwitchActions struct {
	logger                  logger.Logger
	twirBus                 *buscore.Bus
	rateLimiter             redis.SlidingWindow
	sentMessagesRepository  sentmessages.Repository
	channelsRepository      channels.Repository
	toxicMessagesRepository toxic_messages.Repository
	gorm                    *gorm.DB
	toxicityCheck           *toxicity_check.Service
	config                  cfg.Config
	channelsCache           *generic_cacher.GenericCacher[channelmodel.Channel]
	kv                      kv.KV
	modTaskDistributor      mod_task_queue.TaskDistributor
}
