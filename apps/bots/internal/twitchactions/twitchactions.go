package twitchactions

import (
	"log/slog"
	"sync"

	"github.com/aidenwallis/go-ratelimiting/redis"
	adapter "github.com/aidenwallis/go-ratelimiting/redis/adapters/go-redis"
	goredis "github.com/redis/go-redis/v9"
	"github.com/twirapp/kv"
	mod_task_queue "github.com/twirapp/twir/apps/bots/internal/mod-task-queue"
	toxicity_check "github.com/twirapp/twir/apps/bots/internal/services/toxicity-check"
	buscore "github.com/twirapp/twir/libs/bus-core"
	channelcache "github.com/twirapp/twir/libs/cache/channel"
	"github.com/twirapp/twir/libs/cache/twitch"
	cfg "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/repositories/channels"
	"github.com/twirapp/twir/libs/repositories/sentmessages"
	"github.com/twirapp/twir/libs/repositories/toxic_messages"
	"gorm.io/gorm"
)

func New(
	logger *slog.Logger,
	sentMessagesRepository sentmessages.Repository,
	channelsRepository channels.Repository,
	toxicMessagesRepository toxic_messages.Repository,
	gormDB *gorm.DB,
	redisClient *goredis.Client,
	toxicityCheck *toxicity_check.Service,
	cfg cfg.Config,
	channelsByTwitchIDCache *channelcache.TwitchUserIDCacher,
	twirBus *buscore.Bus,
	kv kv.KV,
	modTaskDistributor mod_task_queue.TaskDistributor,
	cachedTwitchClient *twitch.CachedTwitchClient,
) *TwitchActions {
	actions := &TwitchActions{
		logger:                  logger,
		config:                  cfg,
		twirBus:                 twirBus,
		gorm:                    gormDB,
		rateLimiter:             redis.NewSlidingWindow(adapter.NewAdapter(redisClient)),
		sentMessagesRepository:  sentMessagesRepository,
		channelsRepository:      channelsRepository,
		toxicityCheck:           toxicityCheck,
		toxicMessagesRepository: toxicMessagesRepository,
		channelsByTwitchIDCache: channelsByTwitchIDCache,
		kv:                      kv,
		modTaskDistributor:      modTaskDistributor,
		cachedTwitchClient:      cachedTwitchClient,
	}

	return actions
}

type TwitchActions struct {
	logger                  *slog.Logger
	twirBus                 *buscore.Bus
	rateLimiter             redis.SlidingWindow
	sentMessagesRepository  sentmessages.Repository
	channelsRepository      channels.Repository
	toxicMessagesRepository toxic_messages.Repository
	gorm                    *gorm.DB
	toxicityCheck           *toxicity_check.Service
	config                  cfg.Config
	channelsByTwitchIDCache *channelcache.TwitchUserIDCacher
	kv                      kv.KV
	modTaskDistributor      mod_task_queue.TaskDistributor
	cachedTwitchClient      *twitch.CachedTwitchClient
	newUserClient           twitchUserClientFactory
	newBotClient            twitchBotClientFactory

	botClientsMu sync.Mutex
	botClients   map[string]cachedBotClient
}
