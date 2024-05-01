package commands

import (
	"time"

	"github.com/redis/go-redis/v9"
	cfg "github.com/satont/twir/libs/config"
	"gorm.io/gorm"
)

const channelCommandsCacheKey = "cache:twir:commands:channel:"
const channelCommandsCacheTTL = 24 * time.Hour

func buildChannelCommandsCacheKey(channelID string) string {
	return channelCommandsCacheKey + channelID
}

type CachedCommandsClient struct {
	redis  *redis.Client
	config cfg.Config
	db     *gorm.DB
}

func New(
	config cfg.Config,
	redisClient *redis.Client,
	db *gorm.DB,
) (
	*CachedCommandsClient,
	error,
) {
	return &CachedCommandsClient{
		redis:  redisClient,
		config: config,
		db:     db,
	}, nil
}
