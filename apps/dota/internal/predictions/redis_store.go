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
	reservationPrefix                 = "pending:"
	terminalClaimPrefix               = "terminal:"
	matchEndedDeliveryPendingPrefix   = "pending:"
	matchEndedDeliveryDeliveredMarker = "delivered"

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
	claimTerminalScript = `
if redis.call("SET", KEYS[1], ARGV[1], "NX", "PX", ARGV[2]) then return 1 end
return 0`
	renewTerminalScript = `
if redis.call("GET", KEYS[1]) == ARGV[1] then return redis.call("PEXPIRE", KEYS[1], ARGV[2]) end
return 0`
	releaseTerminalScript = `
if redis.call("GET", KEYS[1]) == ARGV[1] then return redis.call("DEL", KEYS[1]) end
return 0`
	completeTerminalScript = `
if redis.call("GET", KEYS[3]) == ARGV[1] then
  redis.call("DEL", KEYS[1], KEYS[2], KEYS[3])
  return 1
end
return 0`
	claimMatchEndedDeliveryScript = `
local current = redis.call("GET", KEYS[1])
if not current then
  redis.call("SET", KEYS[1], ARGV[1], "PX", ARGV[2])
  return 1
end
if current == ARGV[3] then return 3 end
return 2`
	completeMatchEndedDeliveryScript = `
if redis.call("GET", KEYS[1]) == ARGV[1] then
  redis.call("SET", KEYS[1], ARGV[2], "PX", ARGV[3])
  return 1
end
return 0`
	renewMatchEndedDeliveryScript = `
local current = redis.call("GET", KEYS[1])
if current == ARGV[1] then
  return redis.call("PEXPIRE", KEYS[1], ARGV[2])
end
return 0`
	releaseMatchEndedDeliveryScript = `
local current = redis.call("GET", KEYS[1])
if current == ARGV[1] then return redis.call("DEL", KEYS[1]) end
return 0`
)

type redisPredictionClient interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Eval(ctx context.Context, script string, keys []string, args ...interface{}) *redis.Cmd
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

func (s *RedisPredictionStore) ClaimTerminal(
	ctx context.Context,
	key string,
	token string,
	ttl time.Duration,
) (bool, error) {
	if ttl <= 0 {
		return false, fmt.Errorf("prediction terminal TTL must be positive")
	}
	if strings.TrimSpace(token) == "" {
		return false, fmt.Errorf("prediction terminal token is required")
	}
	claimed, err := s.client.Eval(
		ctx,
		claimTerminalScript,
		[]string{terminalPredictionKey(key)},
		terminalClaimPrefix+token,
		ttl.Milliseconds(),
	).Int64()
	if err != nil {
		return false, fmt.Errorf("claim prediction terminal: %w", err)
	}
	if claimed == 0 {
		return false, nil
	}
	if claimed != 1 {
		return false, fmt.Errorf("unexpected prediction terminal claim result: %d", claimed)
	}

	return true, nil
}

func (s *RedisPredictionStore) RenewTerminal(
	ctx context.Context,
	key string,
	token string,
	ttl time.Duration,
) (bool, error) {
	if ttl <= 0 {
		return false, fmt.Errorf("prediction terminal TTL must be positive")
	}
	if strings.TrimSpace(token) == "" {
		return false, fmt.Errorf("prediction terminal token is required")
	}
	renewed, err := s.client.Eval(
		ctx,
		renewTerminalScript,
		[]string{terminalPredictionKey(key)},
		terminalClaimPrefix+token,
		ttl.Milliseconds(),
	).Int64()
	if err != nil {
		return false, fmt.Errorf("renew prediction terminal: %w", err)
	}
	if renewed == 0 {
		return false, nil
	}
	if renewed != 1 {
		return false, fmt.Errorf("unexpected prediction terminal renewal result: %d", renewed)
	}

	return true, nil
}

func (s *RedisPredictionStore) CompleteTerminal(
	ctx context.Context,
	key string,
	token string,
) (bool, error) {
	completed, err := s.client.Eval(
		ctx,
		completeTerminalScript,
		[]string{key, pendingIntentKey(key), terminalPredictionKey(key)},
		terminalClaimPrefix+token,
	).Int64()
	if err != nil {
		return false, fmt.Errorf("complete prediction terminal: %w", err)
	}
	if completed == 0 {
		return false, nil
	}
	if completed != 1 {
		return false, fmt.Errorf("unexpected prediction terminal completion result: %d", completed)
	}

	return true, nil
}

func (s *RedisPredictionStore) ReleaseTerminal(ctx context.Context, key string, token string) error {
	return s.client.Eval(
		ctx,
		releaseTerminalScript,
		[]string{terminalPredictionKey(key)},
		terminalClaimPrefix+token,
	).Err()
}

func (s *RedisPredictionStore) ClaimMatchEndedDelivery(
	ctx context.Context,
	key string,
	token string,
	ttl time.Duration,
) (matchEndedDeliveryState, error) {
	if ttl <= 0 {
		return 0, fmt.Errorf("match ended delivery TTL must be positive")
	}
	if strings.TrimSpace(token) == "" {
		return 0, fmt.Errorf("match ended delivery token is required")
	}

	state, err := s.client.Eval(
		ctx,
		claimMatchEndedDeliveryScript,
		[]string{key},
		matchEndedDeliveryPendingPrefix+token,
		ttl.Milliseconds(),
		matchEndedDeliveryDeliveredMarker,
	).Int64()
	if err != nil {
		return 0, fmt.Errorf("claim match ended delivery: %w", err)
	}
	switch state {
	case 1:
		return matchEndedDeliveryAcquired, nil
	case 2:
		return matchEndedDeliveryPending, nil
	case 3:
		return matchEndedDeliveryDelivered, nil
	default:
		return 0, fmt.Errorf("unexpected match ended delivery claim result: %d", state)
	}
}

func (s *RedisPredictionStore) CompleteMatchEndedDelivery(
	ctx context.Context,
	key string,
	token string,
	ttl time.Duration,
) (bool, error) {
	if ttl <= 0 {
		return false, fmt.Errorf("match ended delivery TTL must be positive")
	}
	if strings.TrimSpace(token) == "" {
		return false, fmt.Errorf("match ended delivery token is required")
	}

	completed, err := s.client.Eval(
		ctx,
		completeMatchEndedDeliveryScript,
		[]string{key},
		matchEndedDeliveryPendingPrefix+token,
		matchEndedDeliveryDeliveredMarker,
		ttl.Milliseconds(),
	).Int64()
	if err != nil {
		return false, fmt.Errorf("complete match ended delivery: %w", err)
	}
	if completed == 0 {
		return false, nil
	}
	if completed != 1 {
		return false, fmt.Errorf("unexpected match ended delivery completion result: %d", completed)
	}

	return true, nil
}

func (s *RedisPredictionStore) RenewMatchEndedDelivery(
	ctx context.Context,
	key string,
	token string,
	ttl time.Duration,
) (bool, error) {
	if ttl <= 0 {
		return false, fmt.Errorf("match ended delivery TTL must be positive")
	}
	if strings.TrimSpace(token) == "" {
		return false, fmt.Errorf("match ended delivery token is required")
	}

	renewed, err := s.client.Eval(
		ctx,
		renewMatchEndedDeliveryScript,
		[]string{key},
		matchEndedDeliveryPendingPrefix+token,
		ttl.Milliseconds(),
	).Int64()
	if err != nil {
		return false, fmt.Errorf("renew match ended delivery: %w", err)
	}
	if renewed == 0 {
		return false, nil
	}
	if renewed != 1 {
		return false, fmt.Errorf("unexpected match ended delivery renewal result: %d", renewed)
	}

	return true, nil
}

func (s *RedisPredictionStore) ReleaseMatchEndedDelivery(ctx context.Context, key string, token string) error {
	return s.client.Eval(
		ctx,
		releaseMatchEndedDeliveryScript,
		[]string{key},
		matchEndedDeliveryPendingPrefix+token,
	).Err()
}

func pendingIntentKey(key string) string {
	return key + ":pending-intent"
}

func terminalPredictionKey(key string) string {
	return key + ":terminal"
}
