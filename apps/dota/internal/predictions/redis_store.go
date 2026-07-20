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
	reservationPrefix = "pending:"

	reservePendingIntentScript = `
if redis.call("EXISTS", KEYS[1]) == 1 then return 0 end
redis.call("SET", KEYS[1], ARGV[1], "PX", ARGV[3])
redis.call("SET", KEYS[2], ARGV[2], "PX", ARGV[3])
return 1`
	compareAndDeleteScript = `
if redis.call("GET", KEYS[1]) == ARGV[1] then return redis.call("DEL", KEYS[1], KEYS[2]) end
return 0`
	commitReservationScript = `
local current = redis.call("GET", KEYS[1])
if current == ARGV[1] then
  redis.call("SET", KEYS[1], ARGV[2], "PX", ARGV[3])
  redis.call("DEL", KEYS[2])
  return 1
end
if current == ARGV[2] then
  redis.call("DEL", KEYS[2])
  return 2
end
return 0`
)

type redisPredictionClient interface {
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
	intent pendingPredictionIntent,
	ttl time.Duration,
) (bool, error) {
	if ttl <= 0 {
		return false, fmt.Errorf("prediction TTL must be positive")
	}
	if err := intent.validate(); err != nil {
		return false, fmt.Errorf("invalid pending prediction intent: %w", err)
	}
	data, err := json.Marshal(intent)
	if err != nil {
		return false, err
	}
	reserved, err := s.client.Eval(
		ctx,
		reservePendingIntentScript,
		[]string{key, pendingIntentKey(key)},
		reservationPrefix+intent.Token,
		string(data),
		ttl.Milliseconds(),
	).Int64()
	if err != nil {
		return false, fmt.Errorf("reserve prediction intent: %w", err)
	}
	if reserved == 0 {
		return false, nil
	}
	if reserved != 1 {
		return false, fmt.Errorf("unexpected prediction reservation result: %d", reserved)
	}
	return true, nil
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
		[]string{key, pendingIntentKey(key)},
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

func (s *RedisPredictionStore) GetPending(
	ctx context.Context,
	key string,
) (pendingPredictionIntent, error) {
	value, err := s.client.Get(ctx, pendingIntentKey(key)).Result()
	if err == redis.Nil {
		return pendingPredictionIntent{}, errPredictionIntentNotFound
	}
	if err != nil {
		return pendingPredictionIntent{}, err
	}

	var intent pendingPredictionIntent
	if err := json.Unmarshal([]byte(value), &intent); err != nil {
		return pendingPredictionIntent{}, fmt.Errorf("decode pending prediction intent: %w", err)
	}
	return intent, nil
}

func (s *RedisPredictionStore) Release(ctx context.Context, key string, token string) error {
	return s.client.Eval(
		ctx,
		compareAndDeleteScript,
		[]string{key, pendingIntentKey(key)},
		reservationPrefix+token,
	).Err()
}

func (s *RedisPredictionStore) Delete(ctx context.Context, key string) error {
	return s.client.Del(ctx, key, pendingIntentKey(key)).Err()
}

func pendingIntentKey(key string) string {
	return key + ":pending-intent"
}
