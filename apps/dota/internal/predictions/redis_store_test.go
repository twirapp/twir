package predictions

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

type fakeRedisPredictionClient struct {
	values    map[string]string
	evalCalls []redisEvalCall
}

type redisEvalCall struct {
	script string
	keys   []string
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
	c.evalCalls = append(c.evalCalls, redisEvalCall{
		script: script,
		keys:   append([]string(nil), keys...),
		args:   append([]interface{}(nil), args...),
	})

	switch script {
	case reservePendingIntentScript:
		if len(keys) != 2 {
			return redis.NewCmdResult(nil, errors.New("unexpected Redis reserve keys"))
		}
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
		if len(keys) != 2 {
			return redis.NewCmdResult(nil, errors.New("unexpected Redis release keys"))
		}
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
		if len(keys) != 2 {
			return redis.NewCmdResult(nil, errors.New("unexpected Redis commit keys"))
		}
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

	case claimTerminalScript:
		if len(keys) != 1 || len(args) != 2 {
			return redis.NewCmdResult(nil, errors.New("unexpected Redis terminal claim arguments"))
		}
		token, tokenOK := args[0].(string)
		_, ttlOK := args[1].(int64)
		if !tokenOK || !ttlOK {
			return redis.NewCmdResult(nil, errors.New("unexpected Redis terminal claim values"))
		}
		if _, claimed := c.values[keys[0]]; claimed {
			return redis.NewCmdResult(int64(0), nil)
		}
		c.values[keys[0]] = token
		return redis.NewCmdResult(int64(1), nil)

	case renewTerminalScript:
		if len(keys) != 1 || len(args) != 2 {
			return redis.NewCmdResult(nil, errors.New("unexpected Redis terminal renewal arguments"))
		}
		token, tokenOK := args[0].(string)
		_, ttlOK := args[1].(int64)
		if !tokenOK || !ttlOK {
			return redis.NewCmdResult(nil, errors.New("unexpected Redis terminal renewal values"))
		}
		if c.values[keys[0]] == token {
			return redis.NewCmdResult(int64(1), nil)
		}
		return redis.NewCmdResult(int64(0), nil)

	case releaseTerminalScript:
		if len(keys) != 1 || len(args) != 1 {
			return redis.NewCmdResult(nil, errors.New("unexpected Redis terminal release arguments"))
		}
		token, ok := args[0].(string)
		if !ok {
			return redis.NewCmdResult(nil, errors.New("unexpected Redis terminal release value"))
		}
		if c.values[keys[0]] == token {
			delete(c.values, keys[0])
			return redis.NewCmdResult(int64(1), nil)
		}
		return redis.NewCmdResult(int64(0), nil)

	case completeTerminalScript:
		if len(keys) != 3 || len(args) != 1 {
			return redis.NewCmdResult(nil, errors.New("unexpected Redis terminal completion arguments"))
		}
		token, ok := args[0].(string)
		if !ok {
			return redis.NewCmdResult(nil, errors.New("unexpected Redis terminal completion value"))
		}
		if c.values[keys[2]] == token {
			delete(c.values, keys[0])
			delete(c.values, keys[1])
			delete(c.values, keys[2])
			return redis.NewCmdResult(int64(1), nil)
		}
		return redis.NewCmdResult(int64(0), nil)

	case claimMatchEndedDeliveryScript:
		if len(keys) != 1 || len(args) != 3 {
			return redis.NewCmdResult(nil, errors.New("unexpected Redis match ended delivery claim arguments"))
		}
		pending, pendingOK := args[0].(string)
		_, ttlOK := args[1].(int64)
		delivered, deliveredOK := args[2].(string)
		if !pendingOK || !ttlOK || !deliveredOK {
			return redis.NewCmdResult(nil, errors.New("unexpected Redis match ended delivery claim values"))
		}
		current, exists := c.values[keys[0]]
		if !exists {
			c.values[keys[0]] = pending
			return redis.NewCmdResult(int64(1), nil)
		}
		if current == delivered {
			return redis.NewCmdResult(int64(3), nil)
		}
		return redis.NewCmdResult(int64(2), nil)

	case completeMatchEndedDeliveryScript:
		if len(keys) != 1 || len(args) != 3 {
			return redis.NewCmdResult(nil, errors.New("unexpected Redis match ended delivery completion arguments"))
		}
		pending, pendingOK := args[0].(string)
		delivered, deliveredOK := args[1].(string)
		_, ttlOK := args[2].(int64)
		if !pendingOK || !deliveredOK || !ttlOK {
			return redis.NewCmdResult(nil, errors.New("unexpected Redis match ended delivery completion values"))
		}
		if c.values[keys[0]] == pending {
			c.values[keys[0]] = delivered
			return redis.NewCmdResult(int64(1), nil)
		}
		return redis.NewCmdResult(int64(0), nil)

	case renewMatchEndedDeliveryScript:
		if len(keys) != 1 || len(args) != 2 {
			return redis.NewCmdResult(nil, errors.New("unexpected Redis match ended delivery renewal arguments"))
		}
		pending, pendingOK := args[0].(string)
		_, ttlOK := args[1].(int64)
		if !pendingOK || !ttlOK {
			return redis.NewCmdResult(nil, errors.New("unexpected Redis match ended delivery renewal values"))
		}
		if c.values[keys[0]] == pending {
			return redis.NewCmdResult(int64(1), nil)
		}
		return redis.NewCmdResult(int64(0), nil)

	case releaseMatchEndedDeliveryScript:
		if len(keys) != 1 || len(args) != 1 {
			return redis.NewCmdResult(nil, errors.New("unexpected Redis match ended delivery release arguments"))
		}
		pending, ok := args[0].(string)
		if !ok {
			return redis.NewCmdResult(nil, errors.New("unexpected Redis match ended delivery release value"))
		}
		if c.values[keys[0]] == pending {
			delete(c.values, keys[0])
			return redis.NewCmdResult(int64(1), nil)
		}
		return redis.NewCmdResult(int64(0), nil)

	default:
		return redis.NewCmdResult(nil, errors.New("unexpected Redis EVAL script"))
	}
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
		Title:           "Will the streamer win? [d:AAAAAAAAAAA]",
		Correlation:     "AAAAAAAAAAA",
		YesOutcomeTitle: "Yes",
		NoOutcomeTitle:  "No",
		ReservedAt:      time.Date(2026, time.July, 20, 12, 0, 0, 0, time.UTC),
	}

	reserved, err := store.Reserve(ctx, key, intent, time.Hour)

	require.NoError(t, err)
	require.True(t, reserved)
	require.Equal(t, reservationPrefix+intent.Token, client.values[key])
	var persisted map[string]any
	require.NoError(t, json.Unmarshal([]byte(client.values[pendingIntentKey(key)]), &persisted))
	require.Equal(t, "AAAAAAAAAAA", persisted["correlation"])
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

func TestRedisPredictionStoreTerminalOwnershipPreservesNewOwnerData(t *testing.T) {
	ctx := context.Background()
	key := "cache:twir:dota:prediction:channel:match"
	client := newFakeRedisPredictionClient()
	store := &RedisPredictionStore{client: client}
	record := storedPrediction{
		PredictionID: "prediction-id",
		YesOutcomeID: "yes-outcome-id",
		NoOutcomeID:  "no-outcome-id",
	}
	encodedRecord, err := json.Marshal(record)
	require.NoError(t, err)
	client.values[key] = string(encodedRecord)

	claimed, err := store.ClaimTerminal(ctx, key, "first-owner", time.Hour)
	require.NoError(t, err)
	require.True(t, claimed)

	delete(client.values, terminalPredictionKey(key)) // Simulate expiration before another worker claims it.
	claimed, err = store.ClaimTerminal(ctx, key, "second-owner", time.Hour)
	require.NoError(t, err)
	require.True(t, claimed)

	require.NoError(t, store.ReleaseTerminal(ctx, key, "first-owner"))
	require.Equal(t, terminalClaimPrefix+"second-owner", client.values[terminalPredictionKey(key)])
	require.Equal(t, string(encodedRecord), client.values[key])

	completed, err := store.CompleteTerminal(ctx, key, "first-owner")
	require.NoError(t, err)
	require.False(t, completed)
	require.Equal(t, terminalClaimPrefix+"second-owner", client.values[terminalPredictionKey(key)])
	require.Equal(t, string(encodedRecord), client.values[key])

	completed, err = store.CompleteTerminal(ctx, key, "second-owner")
	require.NoError(t, err)
	require.True(t, completed)
	_, recordExists := client.values[key]
	require.False(t, recordExists)
	_, claimExists := client.values[terminalPredictionKey(key)]
	require.False(t, claimExists)
	require.Equal(t, []string{
		claimTerminalScript,
		claimTerminalScript,
		releaseTerminalScript,
		completeTerminalScript,
		completeTerminalScript,
	}, []string{
		client.evalCalls[0].script,
		client.evalCalls[1].script,
		client.evalCalls[2].script,
		client.evalCalls[3].script,
		client.evalCalls[4].script,
	})
	require.Equal(t, time.Hour.Milliseconds(), client.evalCalls[0].args[1])
}

func TestRedisPredictionStoreRenewTerminalRequiresCurrentOwner(t *testing.T) {
	ctx := context.Background()
	key := "cache:twir:dota:prediction:channel:match"
	client := newFakeRedisPredictionClient()
	store := &RedisPredictionStore{client: client}
	client.values[terminalPredictionKey(key)] = terminalClaimPrefix + "owner"

	renewed, err := store.RenewTerminal(ctx, key, "stale-owner", time.Hour)
	require.NoError(t, err)
	require.False(t, renewed)
	require.Equal(t, terminalClaimPrefix+"owner", client.values[terminalPredictionKey(key)])

	renewed, err = store.RenewTerminal(ctx, key, "owner", time.Hour)
	require.NoError(t, err)
	require.True(t, renewed)
	require.Equal(t, terminalClaimPrefix+"owner", client.values[terminalPredictionKey(key)])
	require.Equal(t, []string{renewTerminalScript, renewTerminalScript}, []string{
		client.evalCalls[0].script,
		client.evalCalls[1].script,
	})
	require.Equal(t, time.Hour.Milliseconds(), client.evalCalls[1].args[1])
}

func TestRedisPredictionStoreRenewMatchEndedDeliveryRequiresCurrentPendingOwner(t *testing.T) {
	ctx := context.Background()
	key := matchEndedDeliveryKey(uuid.New(), 41)
	client := newFakeRedisPredictionClient()
	store := &RedisPredictionStore{client: client}
	client.values[key] = matchEndedDeliveryPendingPrefix + "owner"

	renewed, err := store.RenewMatchEndedDelivery(ctx, key, "stale-owner", 2*time.Minute)
	require.NoError(t, err)
	require.False(t, renewed)
	require.Equal(t, matchEndedDeliveryPendingPrefix+"owner", client.values[key])

	renewed, err = store.RenewMatchEndedDelivery(ctx, key, "owner", 2*time.Minute)
	require.NoError(t, err)
	require.True(t, renewed)
	require.Equal(t, matchEndedDeliveryPendingPrefix+"owner", client.values[key])

	client.values[key] = matchEndedDeliveryDeliveredMarker
	renewed, err = store.RenewMatchEndedDelivery(ctx, key, "owner", 2*time.Minute)
	require.NoError(t, err)
	require.False(t, renewed)
	require.Equal(t, matchEndedDeliveryDeliveredMarker, client.values[key])
	require.Equal(t, []string{
		renewMatchEndedDeliveryScript,
		renewMatchEndedDeliveryScript,
		renewMatchEndedDeliveryScript,
	}, []string{
		client.evalCalls[0].script,
		client.evalCalls[1].script,
		client.evalCalls[2].script,
	})
	require.Equal(t, (2 * time.Minute).Milliseconds(), client.evalCalls[1].args[1])
}

func TestRedisPredictionStoreMatchEndedDeliveryStateMachinePreservesTokenOwnership(t *testing.T) {
	ctx := context.Background()
	channelID := uuid.New()
	key := matchEndedDeliveryKey(channelID, 42)
	require.Equal(t, key, matchEndedDeliveryKey(channelID, 42))
	require.NotEqual(t, terminalPredictionKey(predictionKey(channelID, 42)), key)
	client := newFakeRedisPredictionClient()
	store := &RedisPredictionStore{client: client}
	pendingTTL := 2 * time.Minute
	deliveredTTL := 7 * 24 * time.Hour

	state, err := store.ClaimMatchEndedDelivery(ctx, key, "first-owner", pendingTTL)
	require.NoError(t, err)
	require.Equal(t, matchEndedDeliveryAcquired, state)

	state, err = store.ClaimMatchEndedDelivery(ctx, key, "second-owner", pendingTTL)
	require.NoError(t, err)
	require.Equal(t, matchEndedDeliveryPending, state)

	completed, err := store.CompleteMatchEndedDelivery(ctx, key, "stale-owner", deliveredTTL)
	require.NoError(t, err)
	require.False(t, completed)
	require.Equal(t, matchEndedDeliveryPendingPrefix+"first-owner", client.values[key])

	completed, err = store.CompleteMatchEndedDelivery(ctx, key, "first-owner", deliveredTTL)
	require.NoError(t, err)
	require.True(t, completed)
	require.Equal(t, matchEndedDeliveryDeliveredMarker, client.values[key])

	state, err = store.ClaimMatchEndedDelivery(ctx, key, "second-owner", pendingTTL)
	require.NoError(t, err)
	require.Equal(t, matchEndedDeliveryDelivered, state)

	require.NoError(t, store.ReleaseMatchEndedDelivery(ctx, key, "first-owner"))
	require.Equal(t, matchEndedDeliveryDeliveredMarker, client.values[key])

	reclaimedKey := matchEndedDeliveryKey(channelID, 43)
	state, err = store.ClaimMatchEndedDelivery(ctx, reclaimedKey, "first-owner", pendingTTL)
	require.NoError(t, err)
	require.Equal(t, matchEndedDeliveryAcquired, state)
	delete(client.values, reclaimedKey) // Simulate expiry before a replacement owner claims it.
	state, err = store.ClaimMatchEndedDelivery(ctx, reclaimedKey, "second-owner", pendingTTL)
	require.NoError(t, err)
	require.Equal(t, matchEndedDeliveryAcquired, state)

	require.NoError(t, store.ReleaseMatchEndedDelivery(ctx, reclaimedKey, "first-owner"))
	completed, err = store.CompleteMatchEndedDelivery(ctx, reclaimedKey, "first-owner", deliveredTTL)
	require.NoError(t, err)
	require.False(t, completed)
	require.Equal(t, matchEndedDeliveryPendingPrefix+"second-owner", client.values[reclaimedKey])
	require.NoError(t, store.ReleaseMatchEndedDelivery(ctx, reclaimedKey, "second-owner"))
	_, pendingExists := client.values[reclaimedKey]
	require.False(t, pendingExists)
	require.Equal(t, []string{
		claimMatchEndedDeliveryScript,
		claimMatchEndedDeliveryScript,
		completeMatchEndedDeliveryScript,
		completeMatchEndedDeliveryScript,
		claimMatchEndedDeliveryScript,
		releaseMatchEndedDeliveryScript,
		claimMatchEndedDeliveryScript,
		claimMatchEndedDeliveryScript,
		releaseMatchEndedDeliveryScript,
		completeMatchEndedDeliveryScript,
		releaseMatchEndedDeliveryScript,
	}, []string{
		client.evalCalls[0].script,
		client.evalCalls[1].script,
		client.evalCalls[2].script,
		client.evalCalls[3].script,
		client.evalCalls[4].script,
		client.evalCalls[5].script,
		client.evalCalls[6].script,
		client.evalCalls[7].script,
		client.evalCalls[8].script,
		client.evalCalls[9].script,
		client.evalCalls[10].script,
	})
	require.Equal(t, pendingTTL.Milliseconds(), client.evalCalls[0].args[1])
	require.Equal(t, deliveredTTL.Milliseconds(), client.evalCalls[3].args[2])
}
