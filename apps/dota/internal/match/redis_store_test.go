package match

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strings"
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
	values               map[string]string
	actions              []string
	getCalls             int
	evalCalls            []redisStateEvalCall
	failXAdd             bool
	failSet              bool
	commitNextBeforeEval bool
	onEval               func(*fakeRedisStateClient, redisStateEvalCall) *redis.Cmd
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

	appendActions := func() error {
		for _, arg := range args[4:] {
			action, ok := arg.(string)
			if !ok {
				return errors.New("unexpected Redis action value")
			}
			if c.failXAdd {
				return errors.New("simulated Redis XADD failure")
			}
			c.actions = append(c.actions, action)
		}
		return nil
	}
	if c.commitNextBeforeEval {
		if err := appendActions(); err != nil {
			return redis.NewCmdResult(nil, err)
		}
		c.values[keys[0]] = next
		c.commitNextBeforeEval = false
	}

	current := c.values[keys[0]]
	if strings.Contains(script, "if current == ARGV[2]") && current == next {
		return redis.NewCmdResult(int64(2), nil)
	}
	if current != expected {
		return redis.NewCmdResult(int64(0), nil)
	}

	xaddIndex := strings.Index(script, `redis.call("XADD"`)
	setIndex := strings.Index(script, `redis.call("SET"`)
	if xaddIndex == -1 || setIndex == -1 {
		return redis.NewCmdResult(nil, errors.New("missing Redis state mutation commands"))
	}

	if xaddIndex < setIndex {
		if err := appendActions(); err != nil {
			return redis.NewCmdResult(nil, err)
		}
		if c.failSet {
			return redis.NewCmdResult(nil, errors.New("simulated Redis SET failure"))
		}
		c.values[keys[0]] = next
		return redis.NewCmdResult(int64(1), nil)
	}

	if c.failSet {
		return redis.NewCmdResult(nil, errors.New("simulated Redis SET failure"))
	}
	c.values[keys[0]] = next
	if err := appendActions(); err != nil {
		return redis.NewCmdResult(nil, err)
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
	expectedSnapshotKey := "cache:twir:dota:{dota}:matchstate:" + channelID.String()
	expectedStreamKey := "stream:twir:dota:{dota}:lifecycle-actions"
	require.Equal(t, expectedSnapshotKey, snapshotKey(channelID))
	require.Equal(t, expectedStreamKey, lifecycleActionStreamKey)
	require.Equal(t, []string{expectedSnapshotKey, expectedStreamKey}, call.keys)
	require.Len(t, call.args, 5)
	require.Equal(t, absentSnapshotSentinel, call.args[0])
	require.Equal(t, snapshotTTL.Milliseconds(), call.args[2])
	require.Equal(t, lifecycleActionStreamMaxLen, call.args[3])

	var saved Snapshot
	require.NoError(t, json.Unmarshal([]byte(call.args[1].(string)), &saved))
	expected := next
	expected.MutationID = saved.MutationID
	require.NotEmpty(t, saved.MutationID)
	require.Equal(t, expected, saved)

	expectedAction := action
	expectedAction.MutationID = saved.MutationID
	actionJSON, err := json.Marshal(expectedAction)
	require.NoError(t, err)
	require.Equal(t, string(actionJSON), call.args[4])
	require.Equal(t, []string{string(actionJSON)}, client.actions)
}

func TestRedisStateStoreCompareAndSwapKeepsStateWhenActionAppendFails(t *testing.T) {
	ctx := context.Background()
	channelID := uuid.New()
	client := newFakeRedisStateClient()
	client.failXAdd = true
	store := &RedisStateStore{client: client}
	current := Snapshot{
		ChannelID: channelID,
		State:     StateInGame,
		MatchID:   12345,
		Revision:  1,
	}
	client.values[snapshotKey(channelID)] = mustMarshalSnapshot(t, current)
	next := current
	next.Revision++
	next.State = StateIdle
	next.MatchID = 0
	action := LifecycleAction{
		Kind:      ActionCancel,
		ChannelID: channelID,
		MatchID:   current.MatchID,
		Revision:  next.Revision,
	}

	swapped, err := store.CompareAndSwap(ctx, current, next, []LifecycleAction{action})

	require.Error(t, err)
	require.False(t, swapped)
	require.Equal(t, mustMarshalSnapshot(t, current), client.values[snapshotKey(channelID)])
	require.Empty(t, client.actions)
}

func TestRedisStateStoreCompareAndSwapTreatsCommittedMutationRetryAsSuccess(t *testing.T) {
	ctx := context.Background()
	channelID := uuid.New()
	client := newFakeRedisStateClient()
	client.commitNextBeforeEval = true
	store := &RedisStateStore{client: client}
	current := Snapshot{
		ChannelID: channelID,
		State:     StateInGame,
		MatchID:   12345,
		Revision:  1,
	}
	client.values[snapshotKey(channelID)] = mustMarshalSnapshot(t, current)
	next := current
	next.Revision++
	next.State = StateIdle
	next.MatchID = 0
	action := LifecycleAction{
		Kind:      ActionCancel,
		ChannelID: channelID,
		MatchID:   current.MatchID,
		Revision:  next.Revision,
	}

	swapped, err := store.CompareAndSwap(ctx, current, next, []LifecycleAction{action})

	require.NoError(t, err)
	require.True(t, swapped)
	require.Len(t, client.evalCalls, 1)
	require.Len(t, client.actions, 1)
	var saved Snapshot
	require.NoError(t, json.Unmarshal([]byte(client.values[snapshotKey(channelID)]), &saved))
	require.NotEmpty(t, saved.MutationID)
	var streamedAction LifecycleAction
	require.NoError(t, json.Unmarshal([]byte(client.actions[0]), &streamedAction))
	require.True(t, ActionMatchesSnapshot(streamedAction, saved))
}

func TestRedisStateStoreCompareAndSwapPreservesMultipleActionOrder(t *testing.T) {
	ctx := context.Background()
	channelID := uuid.New()
	client := newFakeRedisStateClient()
	store := &RedisStateStore{client: client}
	current := Snapshot{
		ChannelID: channelID,
		State:     StateInGame,
		MatchID:   12345,
		Revision:  4,
	}
	client.values[snapshotKey(channelID)] = mustMarshalSnapshot(t, current)
	next := current
	next.Revision++
	next.MatchID = 67890
	actions := []LifecycleAction{
		{
			Kind:      ActionCancel,
			ChannelID: channelID,
			MatchID:   current.MatchID,
			Revision:  next.Revision,
		},
		{
			Kind:      ActionCreate,
			ChannelID: channelID,
			MatchID:   next.MatchID,
			Revision:  next.Revision,
		},
	}

	swapped, err := store.CompareAndSwap(ctx, current, next, actions)

	require.NoError(t, err)
	require.True(t, swapped)
	require.Len(t, client.evalCalls, 1)

	call := client.evalCalls[0]
	require.Len(t, call.args, 6)
	var saved Snapshot
	require.NoError(t, json.Unmarshal([]byte(call.args[1].(string)), &saved))
	expectedActions := append([]LifecycleAction(nil), actions...)
	for index := range expectedActions {
		expectedActions[index].MutationID = saved.MutationID
	}
	cancelJSON, err := json.Marshal(expectedActions[0])
	require.NoError(t, err)
	createJSON, err := json.Marshal(expectedActions[1])
	require.NoError(t, err)
	require.Equal(t, []interface{}{string(cancelJSON), string(createJSON)}, call.args[4:])

	decoded := make([]LifecycleAction, 0, len(actions))
	for _, rawAction := range call.args[4:] {
		var action LifecycleAction
		require.NoError(t, json.Unmarshal([]byte(rawAction.(string)), &action))
		decoded = append(decoded, action)
	}
	require.Equal(t, expectedActions, decoded)
}

func TestRedisStateStoreCompareAndSwapFencesActionsAfterSetFailure(t *testing.T) {
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
	client.values[snapshotKey(channelID)] = mustMarshalSnapshot(t, current)
	next := current
	next.Revision++
	next.State = StateIdle
	next.MatchID = 0
	action := LifecycleAction{
		Kind:      ActionCancel,
		ChannelID: channelID,
		MatchID:   current.MatchID,
		Revision:  next.Revision,
	}

	client.failSet = true
	swapped, err := store.CompareAndSwap(ctx, current, next, []LifecycleAction{action})
	require.Error(t, err)
	require.False(t, swapped)
	require.Equal(t, mustMarshalSnapshot(t, current), client.values[snapshotKey(channelID)])
	require.Len(t, client.actions, 1)

	client.failSet = false
	swapped, err = store.CompareAndSwap(ctx, current, next, []LifecycleAction{action})
	require.NoError(t, err)
	require.True(t, swapped)
	require.Len(t, client.actions, 2)

	var staleAction LifecycleAction
	require.NoError(t, json.Unmarshal([]byte(client.actions[0]), &staleAction))
	var committedAction LifecycleAction
	require.NoError(t, json.Unmarshal([]byte(client.actions[1]), &committedAction))
	var committedSnapshot Snapshot
	require.NoError(t, json.Unmarshal([]byte(client.values[snapshotKey(channelID)]), &committedSnapshot))
	require.NotEmpty(t, staleAction.MutationID)
	require.NotEmpty(t, committedAction.MutationID)
	require.NotEqual(t, staleAction.MutationID, committedAction.MutationID)
	require.False(t, ActionMatchesSnapshot(staleAction, committedSnapshot))
	require.True(t, ActionMatchesSnapshot(committedAction, committedSnapshot))
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
		{
			name: "action kind is invalid",
			next: validNext,
			actions: []LifecycleAction{{
				Kind:      ActionKind("invalid"),
				ChannelID: channelID,
				MatchID:   validAction.MatchID,
				Revision:  validAction.Revision,
			}},
		},
		{
			name: "action revision does not match next snapshot",
			next: validNext,
			actions: []LifecycleAction{{
				Kind:      ActionCreate,
				ChannelID: channelID,
				MatchID:   validAction.MatchID,
				Revision:  current.Revision,
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

func TestRedisStateStoreCompareAndSwapRejectsMaxRevisionBeforeEval(t *testing.T) {
	channelID := uuid.New()
	client := newFakeRedisStateClient()
	store := &RedisStateStore{client: client}
	current := Snapshot{
		ChannelID: channelID,
		State:     StateIdle,
		Revision:  math.MaxUint64,
	}
	next := current
	next.Revision++

	swapped, err := store.CompareAndSwap(context.Background(), current, next, nil)

	require.Error(t, err)
	require.False(t, swapped)
	require.Empty(t, client.evalCalls)
}

func TestRedisStateStoreCompareAndSwapRejectsInvalidCreateTransitionBeforeEval(t *testing.T) {
	channelID := uuid.New()
	current := Snapshot{
		ChannelID: channelID,
		State:     StateIdle,
		Revision:  1,
	}
	validNext := current
	validNext.Revision++
	validNext.State = StateInGame
	validNext.MatchID = 12345

	tests := []struct {
		name   string
		next   Snapshot
		action LifecycleAction
	}{
		{
			name: "next match is missing",
			next: Snapshot{
				ChannelID: channelID,
				State:     StateInGame,
				Revision:  validNext.Revision,
			},
			action: LifecycleAction{
				Kind:      ActionCreate,
				ChannelID: channelID,
				MatchID:   validNext.MatchID,
				Revision:  validNext.Revision,
			},
		},
		{
			name: "action match differs from next match",
			next: validNext,
			action: LifecycleAction{
				Kind:      ActionCreate,
				ChannelID: channelID,
				MatchID:   validNext.MatchID + 1,
				Revision:  validNext.Revision,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := newFakeRedisStateClient()
			store := &RedisStateStore{client: client}

			swapped, err := store.CompareAndSwap(context.Background(), current, test.next, []LifecycleAction{test.action})

			require.Error(t, err)
			require.False(t, swapped)
			require.Empty(t, client.evalCalls)
		})
	}
}

func TestRedisStateStoreCompareAndSwapRejectsInvalidTerminalTransitionBeforeEval(t *testing.T) {
	channelID := uuid.New()
	current := Snapshot{
		ChannelID: channelID,
		State:     StateInGame,
		MatchID:   12345,
		Revision:  1,
	}
	validNext := current
	validNext.Revision++
	validNext.State = StateIdle
	validNext.MatchID = 0

	tests := []struct {
		name    string
		current Snapshot
		next    Snapshot
		action  LifecycleAction
	}{
		{
			name: "resolve without current match",
			current: Snapshot{
				ChannelID: channelID,
				State:     StateInGame,
				Revision:  current.Revision,
			},
			next: Snapshot{
				ChannelID: channelID,
				State:     StateIdle,
				Revision:  validNext.Revision,
			},
			action: LifecycleAction{
				Kind:      ActionResolve,
				ChannelID: channelID,
				MatchID:   current.MatchID,
				Revision:  validNext.Revision,
			},
		},
		{
			name: "cancel without current match",
			current: Snapshot{
				ChannelID: channelID,
				State:     StateInGame,
				Revision:  current.Revision,
			},
			next: Snapshot{
				ChannelID: channelID,
				State:     StateIdle,
				Revision:  validNext.Revision,
			},
			action: LifecycleAction{
				Kind:      ActionCancel,
				ChannelID: channelID,
				MatchID:   current.MatchID,
				Revision:  validNext.Revision,
			},
		},
		{
			name:    "resolve match differs from current match",
			current: current,
			next:    validNext,
			action: LifecycleAction{
				Kind:      ActionResolve,
				ChannelID: channelID,
				MatchID:   current.MatchID + 1,
				Revision:  validNext.Revision,
			},
		},
		{
			name:    "cancel match differs from current match",
			current: current,
			next:    validNext,
			action: LifecycleAction{
				Kind:      ActionCancel,
				ChannelID: channelID,
				MatchID:   current.MatchID + 1,
				Revision:  validNext.Revision,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := newFakeRedisStateClient()
			store := &RedisStateStore{client: client}

			swapped, err := store.CompareAndSwap(context.Background(), test.current, test.next, []LifecycleAction{test.action})

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
	require.Empty(t, snapshot.MutationID)
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
	expected.MutationID = saved.MutationID
	require.Equal(t, expected, saved)
}

func TestRedisStateStoreUpdateStatsSupportsSequentialCompareAndSwap(t *testing.T) {
	ctx := context.Background()
	channelID := uuid.New()
	initial := Snapshot{
		ChannelID: channelID,
		State:     StateInGame,
		MatchID:   12345,
		Revision:  4,
		Mmr:       1000,
	}
	client := newFakeRedisStateClient()
	client.values[snapshotKey(channelID)] = mustMarshalSnapshot(t, initial)
	store := &RedisStateStore{client: client}

	require.NoError(t, store.UpdateStats(ctx, channelID, 2000, 1, 0))
	require.NoError(t, store.UpdateStats(ctx, channelID, 2100, 2, 0))

	require.Len(t, client.evalCalls, 2)
	require.Equal(t, client.evalCalls[0].args[1], client.evalCalls[1].args[0])
	var saved Snapshot
	require.NoError(t, json.Unmarshal([]byte(client.values[snapshotKey(channelID)]), &saved))
	require.Equal(t, initial.Revision+2, saved.Revision)
	require.Equal(t, 2100, saved.Mmr)
	require.Equal(t, 2, saved.SessionWins)
	require.Zero(t, saved.SessionLosses)
}

func TestRedisStateStoreCompareAndSwapRejectsCallerProvidedActionMutationIDBeforeEval(t *testing.T) {
	channelID := uuid.New()
	client := newFakeRedisStateClient()
	store := &RedisStateStore{client: client}
	current := Snapshot{
		ChannelID: channelID,
		State:     StateIdle,
		Revision:  1,
	}
	next := current
	next.Revision++
	next.State = StateInGame
	next.MatchID = 12345
	action := LifecycleAction{
		Kind:       ActionCreate,
		ChannelID:  channelID,
		MatchID:    next.MatchID,
		Revision:   next.Revision,
		MutationID: "forged",
	}

	swapped, err := store.CompareAndSwap(context.Background(), current, next, []LifecycleAction{action})

	require.Error(t, err)
	require.False(t, swapped)
	require.Empty(t, client.evalCalls)
}

func TestLifecycleActionJSONRoundTripPreservesResolveFields(t *testing.T) {
	action := LifecycleAction{
		Kind:       ActionResolve,
		ChannelID:  uuid.New(),
		MatchID:    12345,
		Revision:   8,
		MutationID: "mutation-id",
		Win:        true,
		HeroName:   "axe",
	}

	data, err := json.Marshal(action)
	require.NoError(t, err)

	var decoded LifecycleAction
	require.NoError(t, json.Unmarshal(data, &decoded))
	require.Equal(t, action, decoded)
	require.True(t, decoded.Win)
	require.Equal(t, "axe", decoded.HeroName)
	require.Equal(t, "mutation-id", decoded.MutationID)
}

func TestActionMatchesSnapshot(t *testing.T) {
	channelID := uuid.New()
	snapshot := Snapshot{
		ChannelID:  channelID,
		Revision:   8,
		MutationID: "mutation-id",
	}
	action := LifecycleAction{
		ChannelID:  channelID,
		Revision:   snapshot.Revision,
		MutationID: snapshot.MutationID,
	}

	require.True(t, ActionMatchesSnapshot(action, snapshot))

	for _, mutate := range []func(*LifecycleAction){
		func(action *LifecycleAction) { action.ChannelID = uuid.New() },
		func(action *LifecycleAction) { action.Revision++ },
		func(action *LifecycleAction) { action.MutationID = "other-mutation" },
		func(action *LifecycleAction) { action.MutationID = "" },
	} {
		mismatched := action
		mutate(&mismatched)
		require.False(t, ActionMatchesSnapshot(mismatched, snapshot))
	}
}
