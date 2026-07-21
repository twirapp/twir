package match

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

type redisStateEvalCall struct {
	script string
	keys   []string
	args   []interface{}
}

type fakeRedisStateClient struct {
	values    map[string]string
	actions   []string
	getCalls  int
	evalCalls []redisStateEvalCall
	onEval    func(*fakeRedisStateClient, redisStateEvalCall) *redis.Cmd
}

func newFakeRedisStateClient() *fakeRedisStateClient {
	return &fakeRedisStateClient{values: make(map[string]string)}
}

func (c *fakeRedisStateClient) Get(_ context.Context, key string) *redis.StringCmd {
	c.getCalls++
	value, ok := c.values[key]
	if !ok {
		return redis.NewStringResult("", redis.Nil)
	}
	return redis.NewStringResult(value, nil)
}

func (c *fakeRedisStateClient) Eval(
	_ context.Context,
	script string,
	keys []string,
	args ...interface{},
) *redis.Cmd {
	call := redisStateEvalCall{
		script: script,
		keys:   append([]string(nil), keys...),
		args:   append([]interface{}(nil), args...),
	}
	c.evalCalls = append(c.evalCalls, call)

	if c.onEval != nil {
		if result := c.onEval(c, call); result != nil {
			return result
		}
	}

	if script != compareAndSwapSnapshotScript {
		return redis.NewCmdResult(nil, errors.New("unexpected Redis EVAL script"))
	}
	if len(keys) != 2 || len(args) < 4 {
		return redis.NewCmdResult(nil, errors.New("unexpected Redis EVAL arguments"))
	}

	expected, expectedOK := args[0].(string)
	next, nextOK := args[1].(string)
	if !expectedOK || !nextOK {
		return redis.NewCmdResult(nil, errors.New("unexpected Redis state values"))
	}

	current := c.values[keys[0]]
	if current != expected {
		return redis.NewCmdResult(int64(0), nil)
	}

	c.values[keys[0]] = next
	for _, arg := range args[4:] {
		action, ok := arg.(string)
		if !ok {
			return redis.NewCmdResult(nil, errors.New("unexpected Redis action value"))
		}
		c.actions = append(c.actions, action)
	}

	return redis.NewCmdResult(int64(1), nil)
}

func mustMarshalSnapshot(t testing.TB, snapshot Snapshot) string {
	t.Helper()

	data, err := json.Marshal(snapshot)
	require.NoError(t, err)
	return string(data)
}

func TestRedisStateStoreLoadAbsentKeyReturnsIdleRevisionZero(t *testing.T) {
	channelID := uuid.New()
	store := &RedisStateStore{client: newFakeRedisStateClient()}

	snapshot, err := store.Load(context.Background(), channelID)

	require.NoError(t, err)
	require.Equal(t, channelID, snapshot.ChannelID)
	require.Equal(t, StateIdle, snapshot.State)
	require.Zero(t, snapshot.Revision)
}

func TestRedisStateStoreLoadRejectsCorruptOrForeignSnapshot(t *testing.T) {
	channelID := uuid.New()
	otherChannelID := uuid.New()

	for _, raw := range []string{
		"not json",
		fmt.Sprintf(`{"channelId":%q,"state":"idle"}`, otherChannelID.String()),
	} {
		t.Run(raw, func(t *testing.T) {
			client := newFakeRedisStateClient()
			client.values[snapshotKey(channelID)] = raw
			store := &RedisStateStore{client: client}

			_, err := store.Load(context.Background(), channelID)

			require.Error(t, err)
		})
	}
}

func TestRedisStateStoreCompareAndSwapInitialWriteUsesOneAtomicEval(t *testing.T) {
	ctx := context.Background()
	channelID := uuid.New()
	client := newFakeRedisStateClient()
	store := &RedisStateStore{client: client}
	current := Snapshot{
		ChannelID: channelID,
		State:     StateIdle,
	}
	next := current
	next.Revision = 1
	next.State = StateInGame
	next.MatchID = 12345
	action := LifecycleAction{
		Kind:      ActionCreate,
		ChannelID: channelID,
		MatchID:   next.MatchID,
		Revision:  next.Revision,
	}

	swapped, err := store.CompareAndSwap(ctx, current, next, []LifecycleAction{action})

	require.NoError(t, err)
	require.True(t, swapped)
	require.Zero(t, client.getCalls)
	require.Len(t, client.evalCalls, 1)

	call := client.evalCalls[0]
	require.Equal(t, compareAndSwapSnapshotScript, call.script)
	require.Equal(t, []string{snapshotKey(channelID), lifecycleActionStreamKey}, call.keys)
	require.Len(t, call.args, 5)
	require.Equal(t, absentSnapshotSentinel, call.args[0])
	require.Equal(t, snapshotTTL.Milliseconds(), call.args[2])
	require.Equal(t, lifecycleActionStreamMaxLen, call.args[3])

	var saved Snapshot
	require.NoError(t, json.Unmarshal([]byte(call.args[1].(string)), &saved))
	require.Equal(t, next, saved)

	actionJSON, err := json.Marshal(action)
	require.NoError(t, err)
	require.Equal(t, string(actionJSON), call.args[4])
	require.Equal(t, []string{string(actionJSON)}, client.actions)
}

func TestRedisStateStoreCompareAndSwapLossDoesNotAppendAction(t *testing.T) {
	ctx := context.Background()
	channelID := uuid.New()
	client := newFakeRedisStateClient()
	store := &RedisStateStore{client: client}
	current := Snapshot{
		ChannelID: channelID,
		State:     StateInGame,
		MatchID:   12345,
		Revision:  1,
	}
	next := current
	next.Revision++
	next.State = StatePostGame
	action := LifecycleAction{
		Kind:      ActionResolve,
		ChannelID: channelID,
		MatchID:   current.MatchID,
		Revision:  next.Revision,
		Win:       true,
		HeroName:  "axe",
	}
	client.values[snapshotKey(channelID)] = `{"newer":true}`

	swapped, err := store.CompareAndSwap(ctx, current, next, []LifecycleAction{action})

	require.NoError(t, err)
	require.False(t, swapped)
	require.Len(t, client.evalCalls, 1)
	require.Empty(t, client.actions)
}

func TestRedisStateStoreCompareAndSwapRejectsInvalidInputBeforeEval(t *testing.T) {
	channelID := uuid.New()
	otherChannelID := uuid.New()
	current := Snapshot{
		ChannelID: channelID,
		State:     StateIdle,
		Revision:  4,
	}
	validNext := current
	validNext.Revision++
	validAction := LifecycleAction{
		Kind:      ActionCreate,
		ChannelID: channelID,
		MatchID:   12345,
		Revision:  validNext.Revision,
	}

	tests := []struct {
		name    string
		next    Snapshot
		actions []LifecycleAction
	}{
		{
			name: "revision does not increment",
			next: current,
		},
		{
			name: "next channel does not match",
			next: Snapshot{
				ChannelID: otherChannelID,
				State:     StateIdle,
				Revision:  current.Revision + 1,
			},
		},
		{
			name: "action channel is empty",
			next: validNext,
			actions: []LifecycleAction{{
				Kind:     ActionCreate,
				MatchID:  validAction.MatchID,
				Revision: validAction.Revision,
			}},
		},
		{
			name: "action channel does not match",
			next: validNext,
			actions: []LifecycleAction{{
				Kind:      ActionCreate,
				ChannelID: otherChannelID,
				MatchID:   validAction.MatchID,
				Revision:  validAction.Revision,
			}},
		},
		{
			name: "action match is invalid",
			next: validNext,
			actions: []LifecycleAction{{
				Kind:      ActionCreate,
				ChannelID: channelID,
				Revision:  validAction.Revision,
			}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := newFakeRedisStateClient()
			store := &RedisStateStore{client: client}

			swapped, err := store.CompareAndSwap(context.Background(), current, test.next, test.actions)

			require.Error(t, err)
			require.False(t, swapped)
			require.Empty(t, client.evalCalls)
		})
	}
}

func TestSnapshotLoadOldJSONDefaultsOrderingFields(t *testing.T) {
	channelID := uuid.New()
	client := newFakeRedisStateClient()
	client.values[snapshotKey(channelID)] = fmt.Sprintf(
		`{"channelId":%q,"state":"in_game","matchId":12345,"gameTime":321}`,
		channelID.String(),
	)
	store := &RedisStateStore{client: client}

	snapshot, err := store.Load(context.Background(), channelID)

	require.NoError(t, err)
	require.Equal(t, channelID, snapshot.ChannelID)
	require.Equal(t, StateInGame, snapshot.State)
	require.Equal(t, int64(12345), snapshot.MatchID)
	require.Zero(t, snapshot.Revision)
	require.Zero(t, snapshot.LastProviderTimestamp)
	require.Zero(t, snapshot.LastGameTime)
}

func TestRedisStateStoreUpdateStatsRetriesAndPreservesConcurrentState(t *testing.T) {
	ctx := context.Background()
	channelID := uuid.New()
	initial := Snapshot{
		ChannelID: channelID,
		State:     StateIdle,
		Revision:  4,
		Mmr:       1000,
	}
	concurrent := initial
	concurrent.Revision++
	concurrent.State = StateInGame
	concurrent.MatchID = 54321
	concurrent.HeroName = "axe"
	concurrent.GameTime = 678
	concurrent.LastProviderTimestamp = 1_234_567
	concurrent.LastGameTime = 678
	concurrent.SeenEvents = []string{"roshan_killed:500"}

	client := newFakeRedisStateClient()
	client.values[snapshotKey(channelID)] = mustMarshalSnapshot(t, initial)
	client.onEval = func(client *fakeRedisStateClient, call redisStateEvalCall) *redis.Cmd {
		if len(client.evalCalls) != 1 {
			return nil
		}

		client.values[call.keys[0]] = mustMarshalSnapshot(t, concurrent)
		return redis.NewCmdResult(int64(0), nil)
	}
	store := &RedisStateStore{client: client}

	err := store.UpdateStats(ctx, channelID, 3200, 7, 3)

	require.NoError(t, err)
	require.Len(t, client.evalCalls, 2)
	require.Empty(t, client.actions)

	var saved Snapshot
	require.NoError(t, json.Unmarshal([]byte(client.values[snapshotKey(channelID)]), &saved))
	expected := concurrent
	expected.Revision++
	expected.Mmr = 3200
	expected.SessionWins = 7
	expected.SessionLosses = 3
	require.Equal(t, expected, saved)
}

func TestLifecycleActionJSONRoundTripPreservesResolveFields(t *testing.T) {
	action := LifecycleAction{
		Kind:      ActionResolve,
		ChannelID: uuid.New(),
		MatchID:   12345,
		Revision:  8,
		Win:       true,
		HeroName:  "axe",
	}

	data, err := json.Marshal(action)
	require.NoError(t, err)

	var decoded LifecycleAction
	require.NoError(t, json.Unmarshal(data, &decoded))
	require.Equal(t, action, decoded)
	require.True(t, decoded.Win)
	require.Equal(t, "axe", decoded.HeroName)
}
