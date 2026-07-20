package predictions

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	reservationPrefix       = "pending:"
	compareAndDeleteScript  = `if redis.call("GET", KEYS[1]) == ARGV[1] then return redis.call("DEL", KEYS[1]) end return 0`
	commitReservationScript = `local current = redis.call("GET", KEYS[1]) if current == ARGV[1] then redis.call("SET", KEYS[1], ARGV[2], "PX", ARGV[3]) return 1 end if current == ARGV[2] then return 2 end return 0`
)

type redisPredictionClient interface {
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd
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
	token string,
	record storedPrediction,
	ttl time.Duration,
) error {
	if ttl <= 0 {
		return fmt.Errorf("prediction TTL must be positive")
	}
	data, err := json.Marshal(record)
	if err != nil {
		return err
	}
	updated, err := s.client.Eval(
		ctx,
		commitReservationScript,
		[]string{key},
		reservationPrefix+token,
		string(data),
		ttl.Milliseconds(),
	).Int64()
	if err != nil {
		return fmt.Errorf("commit prediction reservation: %w", err)
	}
	if updated != 1 && updated != 2 {
		return errPredictionReservationLost
	}
	return nil
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
