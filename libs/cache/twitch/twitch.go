package twitch

import (
	"github.com/nicklaw5/helix/v2"
	"github.com/redis/go-redis/v9"
	cfg "github.com/satont/twir/libs/config"
	twitchlib "github.com/satont/twir/libs/twitch"
	"github.com/twirapp/twir/libs/grpc/tokens"
)

type CachedTwitchClient struct {
	client *helix.Client
	redis  *redis.Client
}

func New(
	config cfg.Config,
	tokensClient tokens.TokensClient,
	redisClient *redis.Client,
) (
	*CachedTwitchClient,
	error,
) {
	twitchClient, err := twitchlib.NewAppClient(config, tokensClient)
	if err != nil {
		return nil, err
	}

	return &CachedTwitchClient{
		client: twitchClient,
		redis:  redisClient,
	}, nil
}
