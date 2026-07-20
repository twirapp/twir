package predictions

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

type fakeRedisPredictionClient struct {
	values    map[string]string
	evalCalls []redisEvalCall
}

type redisEvalCall struct {
	script string
	args   []interface{}
}

func newFakeRedisPredictionClient() *fakeRedisPredictionClient {
	return &fakeRedisPredictionClient{values: make(map[string]string)}
}

func (c *fakeRedisPredictionClient) SetNX(
	_ context.Context,
	key string,
	value interface{},
	_ time.Duration,
) *redis.BoolCmd {
	if _, exists := c.values[key]; exists {
		return redis.NewBoolResult(false, nil)
	}
	stringValue, ok := value.(string)
	if !ok {
		return redis.NewBoolResult(false, errors.New("unexpected Redis SETNX value"))
	}
	c.values[key] = stringValue
	return redis.NewBoolResult(true, nil)
}

func (c *fakeRedisPredictionClient) Set(
	_ context.Context,
	key string,
	value interface{},
	_ time.Duration,
) *redis.StatusCmd {
	bytesValue, ok := value.([]byte)
	if !ok {
		return redis.NewStatusResult("", errors.New("unexpected Redis SET value"))
	}
	c.values[key] = string(bytesValue)
	return redis.NewStatusResult("OK", nil)
}

func (c *fakeRedisPredictionClient) Get(_ context.Context, key string) *redis.StringCmd {
	value, ok := c.values[key]
	if !ok {
		return redis.NewStringResult("", redis.Nil)
	}
	return redis.NewStringResult(value, nil)
}

func (c *fakeRedisPredictionClient) Eval(
	_ context.Context,
	script string,
	keys []string,
	args ...interface{},
) *redis.Cmd {
	if len(keys) != 2 {
		return redis.NewCmdResult(nil, errors.New("unexpected Redis EVAL keys"))
	}
	c.evalCalls = append(c.evalCalls, redisEvalCall{script: script, args: append([]interface{}(nil), args...)})

	switch script {
	case reservePendingIntentScript:
		if len(args) != 3 {
			return redis.NewCmdResult(nil, errors.New("unexpected Redis reserve arguments"))
		}
		reservation, reservationOK := args[0].(string)
		intent, intentOK := args[1].(string)
		_, ttlOK := args[2].(int64)
		if !reservationOK || !intentOK || !ttlOK {
			return redis.NewCmdResult(nil, errors.New("unexpected Redis reserve values"))
		}
		if _, exists := c.values[keys[0]]; exists {
			return redis.NewCmdResult(int64(0), nil)
		}
		c.values[keys[0]] = reservation
		c.values[keys[1]] = intent
		return redis.NewCmdResult(int64(1), nil)

	case compareAndDeleteScript:
		if len(args) != 1 {
			return redis.NewCmdResult(nil, errors.New("unexpected Redis release arguments"))
		}
		expected, ok := args[0].(string)
		if !ok {
			return redis.NewCmdResult(nil, errors.New("unexpected Redis release value"))
		}
		if c.values[keys[0]] == expected {
			delete(c.values, keys[0])
			delete(c.values, keys[1])
			return redis.NewCmdResult(int64(1), nil)
		}
		return redis.NewCmdResult(int64(0), nil)

	case commitReservationScript:
		if len(args) != 3 {
			return redis.NewCmdResult(nil, errors.New("unexpected Redis commit arguments"))
		}
		expected, expectedOK := args[0].(string)
		value, valueOK := args[1].(string)
		_, ttlOK := args[2].(int64)
		if !expectedOK || !valueOK || !ttlOK {
			return redis.NewCmdResult(nil, errors.New("unexpected Redis commit values"))
		}
		if c.values[keys[0]] != expected {
			if c.values[keys[0]] == value {
				delete(c.values, keys[1])
				return redis.NewCmdResult(int64(2), nil)
			}
			return redis.NewCmdResult(int64(0), nil)
		}
		c.values[keys[0]] = value
		delete(c.values, keys[1])
		return redis.NewCmdResult(int64(1), nil)

	default:
		return redis.NewCmdResult(nil, errors.New("unexpected Redis EVAL script"))
	}
}

func (c *fakeRedisPredictionClient) Del(_ context.Context, keys ...string) *redis.IntCmd {
	var deleted int64
	for _, key := range keys {
		if _, exists := c.values[key]; exists {
			delete(c.values, key)
			deleted++
		}
	}
	return redis.NewIntResult(deleted, nil)
}

func TestRedisPredictionStoreCommitRequiresCurrentReservationOwner(t *testing.T) {
	ctx := context.Background()
	key := "cache:twir:dota:prediction:channel:match"
	client := newFakeRedisPredictionClient()
	store := &RedisPredictionStore{client: client}
	record := storedPrediction{
		PredictionID: "prediction-id",
		YesOutcomeID: "yes-outcome-id",
		NoOutcomeID:  "no-outcome-id",
	}

	client.values[key] = reservationPrefix + "owner"
	client.values[pendingIntentKey(key)] = "pending-intent"
	err := store.Commit(ctx, key, "stale-owner", record, time.Hour)
	require.ErrorIs(t, err, errPredictionReservationLost)
	require.Equal(t, reservationPrefix+"owner", client.values[key])
	require.Equal(t, "pending-intent", client.values[pendingIntentKey(key)])

	require.NoError(t, store.Release(ctx, key, "stale-owner"))
	require.Equal(t, reservationPrefix+"owner", client.values[key])
	require.Equal(t, "pending-intent", client.values[pendingIntentKey(key)])

	encodedRecord, err := json.Marshal(record)
	require.NoError(t, err)
	require.NoError(t, store.Commit(ctx, key, "owner", record, time.Hour))
	require.Equal(t, string(encodedRecord), client.values[key])
	_, pendingExists := client.values[pendingIntentKey(key)]
	require.False(t, pendingExists)
	require.NoError(t, store.Commit(ctx, key, "owner", record, time.Hour))
	require.Equal(t, string(encodedRecord), client.values[key])

	otherRecord := storedPrediction{
		PredictionID: "other-prediction-id",
		YesOutcomeID: "other-yes-outcome-id",
		NoOutcomeID:  "other-no-outcome-id",
	}
	err = store.Commit(ctx, key, "owner", otherRecord, time.Hour)
	require.ErrorIs(t, err, errPredictionReservationLost)
	require.Equal(t, string(encodedRecord), client.values[key])
	require.Equal(t, []string{
		commitReservationScript,
		compareAndDeleteScript,
		commitReservationScript,
		commitReservationScript,
		commitReservationScript,
	}, []string{
		client.evalCalls[0].script,
		client.evalCalls[1].script,
		client.evalCalls[2].script,
		client.evalCalls[3].script,
		client.evalCalls[4].script,
	})
	require.Equal(t, time.Hour.Milliseconds(), client.evalCalls[0].args[2])
}

func TestRedisPredictionStoreAtomicallyPersistsPendingIntent(t *testing.T) {
	ctx := context.Background()
	key := "cache:twir:dota:prediction:channel:match"
	client := newFakeRedisPredictionClient()
	store := &RedisPredictionStore{client: client}
	intent := pendingPredictionIntent{
		Version:         pendingIntentVersion,
		Token:           "owner",
		Title:           "Will the streamer win?",
		YesOutcomeTitle: "Yes",
		NoOutcomeTitle:  "No",
		ReservedAt:      time.Date(2026, time.July, 20, 12, 0, 0, 0, time.UTC),
	}

	reserved, err := store.Reserve(ctx, key, intent, time.Hour)

	require.NoError(t, err)
	require.True(t, reserved)
	require.Equal(t, reservationPrefix+intent.Token, client.values[key])
	encodedIntent, err := json.Marshal(intent)
	require.NoError(t, err)
	require.Equal(t, string(encodedIntent), client.values[pendingIntentKey(key)])
	require.Equal(t, reservePendingIntentScript, client.evalCalls[0].script)

	pending, err := store.GetPending(ctx, key)
	require.NoError(t, err)
	require.Equal(t, intent, pending)

	require.NoError(t, store.Release(ctx, key, "stale-owner"))
	require.Equal(t, reservationPrefix+intent.Token, client.values[key])
	require.Equal(t, string(encodedIntent), client.values[pendingIntentKey(key)])

	require.NoError(t, store.Release(ctx, key, intent.Token))
	_, err = store.GetPending(ctx, key)
	require.ErrorIs(t, err, errPredictionIntentNotFound)
}
