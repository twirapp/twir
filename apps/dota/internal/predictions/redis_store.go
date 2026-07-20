package predictions

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	reservationPrefix      = "pending:"
	compareAndDeleteScript = `if redis.call("GET", KEYS[1]) == ARGV[1] then return redis.call("DEL", KEYS[1]) end return 0`
)

type redisPredictionClient interface {
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Get(ctx context.Context, key string) *redis.StringCmd
	Eval(ctx context.Context, script string, keys []string, args ...interface{}) *redis.Cmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
}

type RedisPredictionStore struct {
	client redisPredictionClient
}

var _ Store = (*RedisPredictionStore)(nil)

func NewRedisPredictionStore(client *redis.Client) *RedisPredictionStore {
	return &RedisPredictionStore{client: client}
}

func (s *RedisPredictionStore) Reserve(
	ctx context.Context,
	key string,
	token string,
	ttl time.Duration,
) (bool, error) {
	return s.client.SetNX(ctx, key, reservationPrefix+token, ttl).Result()
}

func (s *RedisPredictionStore) Commit(
	ctx context.Context,
	key string,
	record storedPrediction,
	ttl time.Duration,
) error {
	data, err := json.Marshal(record)
	if err != nil {
		return err
	}
	return s.client.Set(ctx, key, data, ttl).Err()
}

func (s *RedisPredictionStore) Get(ctx context.Context, key string) (storedPrediction, error) {
	value, err := s.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return storedPrediction{}, errPredictionNotFound
	}
	if err != nil {
		return storedPrediction{}, err
	}
	if strings.HasPrefix(value, reservationPrefix) {
		return storedPrediction{}, errPredictionPending
	}

	var record storedPrediction
	if err := json.Unmarshal([]byte(value), &record); err != nil {
		return storedPrediction{}, err
	}
	return record, nil
}

func (s *RedisPredictionStore) Release(ctx context.Context, key string, token string) error {
	return s.client.Eval(
		ctx,
		compareAndDeleteScript,
		[]string{key},
		reservationPrefix+token,
	).Err()
}

func (s *RedisPredictionStore) Delete(ctx context.Context, key string) error {
	return s.client.Del(ctx, key).Err()
}
