package auth

import (
	"context"
	"time"

	"github.com/alexedwards/scs/goredisstore"
	"github.com/redis/go-redis/v9"
)

const (
	oauthAttemptClaimPrefix   = "oauth:attempt:claim:"
	oauthAttemptClaimLifetime = 15 * time.Minute
)

type oauthAttemptClaimStore interface {
	ClaimOAuthAttempt(ctx context.Context, state string) (bool, error)
	ReleaseOAuthAttempt(ctx context.Context, state string) error
}

type oauthAttemptRedisStore struct {
	*goredisstore.RedisStore
	client *redis.Client
}

func newOAuthAttemptRedisStore(
	store *goredisstore.RedisStore,
	client *redis.Client,
) *oauthAttemptRedisStore {
	return &oauthAttemptRedisStore{
		RedisStore: store,
		client:     client,
	}
}

func (s *oauthAttemptRedisStore) ClaimOAuthAttempt(ctx context.Context, state string) (bool, error) {
	return s.client.SetNX(
		ctx,
		oauthAttemptClaimPrefix+state,
		"1",
		oauthAttemptClaimLifetime,
	).Result()
}

func (s *oauthAttemptRedisStore) ReleaseOAuthAttempt(ctx context.Context, state string) error {
	return s.client.Del(ctx, oauthAttemptClaimPrefix+state).Err()
}
