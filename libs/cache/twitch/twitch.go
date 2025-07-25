package twitch

import (
	"github.com/nicklaw5/helix/v2"
	"github.com/redis/go-redis/v9"
	cfg "github.com/twirapp/twir/libs/config"
	twitchlib "github.com/twirapp/twir/libs/twitch"
	buscore "github.com/twirapp/twir/libs/bus-core"
)

type CachedTwitchClient struct {
	config  cfg.Config
	redis   *redis.Client
	client  *helix.Client
	twirBus *buscore.Bus
}

func New(
	config cfg.Config,
	twirBus *buscore.Bus,
	redisClient *redis.Client,
) (
	*CachedTwitchClient,
	error,
) {
	twitchClient, err := twitchlib.NewAppClient(config, twirBus)
	if err != nil {
		return nil, err
	}

	return &CachedTwitchClient{
		client:  twitchClient,
		redis:   redisClient,
		config:  config,
		twirBus: twirBus,
	}, nil
}
