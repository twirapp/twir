package twitch

import (
	"context"

	"github.com/google/uuid"
	"github.com/kvizyx/twitchy/helix"
	"github.com/redis/go-redis/v9"
	buscore "github.com/twirapp/twir/libs/bus-core"
	cfg "github.com/twirapp/twir/libs/config"
	twitchlib "github.com/twirapp/twir/libs/twitch"
)

type CachedTwitchClient struct {
	config        cfg.Config
	redis         *redis.Client
	client        *helix.Client
	twirBus       *buscore.Bus
	newUserClient func(context.Context, uuid.UUID) (*helix.Client, error)
}

func (c *CachedTwitchClient) createUserClient(ctx context.Context, userID uuid.UUID) (*helix.Client, error) {
	if c.newUserClient != nil {
		return c.newUserClient(ctx, userID)
	}

	return twitchlib.NewUserClientWithContext(ctx, userID, c.config, c.twirBus)
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
