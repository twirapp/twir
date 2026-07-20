package match

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/twirapp/kv"
	kvoptions "github.com/twirapp/kv/options"
	"github.com/twirapp/twir/apps/dota/internal/gsi"
	busapi "github.com/twirapp/twir/libs/bus-core/api"
	busdota "github.com/twirapp/twir/libs/bus-core/dota"
	dotarepository "github.com/twirapp/twir/libs/repositories/dota"
	"github.com/twirapp/twir/libs/repositories/dota/model"
)

type updateMatchResultCall struct {
	channelID uuid.UUID
	won       bool
	mmrDelta  int
}

type fakeRepo struct {
	settings      model.ChannelDotaSettings
	updated       model.ChannelDotaSettings
	getErr        error
	updateErr     error
	updateCalls   []updateMatchResultCall
	getByIDCalled int
}

func (f *fakeRepo) GetByChannelID(
	_ context.Context,
	_ uuid.UUID,
) (model.ChannelDotaSettings, error) {
	f.getByIDCalled++
	if f.getErr != nil {
		return model.Nil, f.getErr
	}
	return f.settings, nil
}

func (f *fakeRepo) GetByGsiToken(
	_ context.Context,
	_ string,
) (model.ChannelDotaSettings, error) {
	return model.Nil, dotarepository.ErrNotFound
}

func (f *fakeRepo) Create(
	_ context.Context,
	_ dotarepository.CreateInput,
) (model.ChannelDotaSettings, error) {
	return model.Nil, errors.New("not implemented")
}

func (f *fakeRepo) Update(
	_ context.Context,
	_ uuid.UUID,
	_ dotarepository.UpdateInput,
) (model.ChannelDotaSettings, error) {
	return model.Nil, errors.New("not implemented")
}

func (f *fakeRepo) UpdateMatchResult(
	_ context.Context,
	channelID uuid.UUID,
	won bool,
	mmrDelta int,
) (model.ChannelDotaSettings, error) {
	if f.updateErr != nil {
		return model.Nil, f.updateErr
	}
	f.updateCalls = append(f.updateCalls, updateMatchResultCall{channelID, won, mmrDelta})
	return f.updated, nil
}

func (f *fakeRepo) ResetSession(
	_ context.Context,
	_ uuid.UUID,
) (model.ChannelDotaSettings, error) {
	return model.Nil, errors.New("not implemented")
}

func (f *fakeRepo) RegenerateGsiToken(
	_ context.Context,
	_ uuid.UUID,
) (model.ChannelDotaSettings, error) {
	return model.Nil, errors.New("not implemented")
}

type fakeEmitter struct {
	mu           sync.Mutex
	roshanErr    error
	aegisErr     error
	matchStarted []busdota.MatchStartedMessage
	matchEnded   []busdota.MatchEndedMessage
	roshanKilled []busdota.RoshanKilledMessage
	aegisPickup  []busdota.AegisPickupMessage
	stateUpdates []busapi.DotaStateUpdateMessage
}

func (f *fakeEmitter) MatchStarted(_ context.Context, msg busdota.MatchStartedMessage) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.matchStarted = append(f.matchStarted, msg)
	return nil
}

func (f *fakeEmitter) MatchEnded(_ context.Context, msg busdota.MatchEndedMessage) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.matchEnded = append(f.matchEnded, msg)
	return nil
}

func (f *fakeEmitter) RoshanKilled(_ context.Context, msg busdota.RoshanKilledMessage) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.roshanErr != nil {
		return f.roshanErr
	}
	f.roshanKilled = append(f.roshanKilled, msg)
	return nil
}

func (f *fakeEmitter) AegisPickup(_ context.Context, msg busdota.AegisPickupMessage) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.aegisErr != nil {
		return f.aegisErr
	}
	f.aegisPickup = append(f.aegisPickup, msg)
	return nil
}

func (f *fakeEmitter) StateUpdate(_ context.Context, msg busapi.DotaStateUpdateMessage) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.stateUpdates = append(f.stateUpdates, msg)
	return nil
}

type fakeValuer struct {
	b   []byte
	err error
}

func (v fakeValuer) Int() (int64, error)     { return 0, v.err }
func (v fakeValuer) String() (string, error) { return string(v.b), v.err }
func (v fakeValuer) Bytes() ([]byte, error)  { return v.b, v.err }
func (v fakeValuer) Bool() (bool, error)     { return false, v.err }
func (v fakeValuer) Float() (float64, error) { return 0, v.err }
func (v fakeValuer) Scan(dest any) error     { return v.err }
func (v fakeValuer) Err() error              { return v.err }

type fakeKV struct {
	mu    sync.Mutex
	store map[string][]byte
	sets  int
}

func newFakeKV() *fakeKV {
	return &fakeKV{store: make(map[string][]byte)}
}

func (f *fakeKV) Get(_ context.Context, key string) kv.Valuer {
	f.mu.Lock()
	defer f.mu.Unlock()
	b, ok := f.store[key]
	if !ok {
		return fakeValuer{err: kv.ErrKeyNil}
	}
	return fakeValuer{b: b}
}

func (f *fakeKV) Set(_ context.Context, key string, value any, _ ...kvoptions.Option) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.sets++
	switch v := value.(type) {
	case []byte:
		f.store[key] = v
	case string:
		f.store[key] = []byte(v)
	default:
		return fmt.Errorf("unsupported value type %T", value)
	}
	return nil
}

func (f *fakeKV) SetMany(_ context.Context, values []kv.SetMany) error {
	for _, v := range values {
		if err := f.Set(context.Background(), v.Key, v.Value, v.Options...); err != nil {
			return err
		}
	}
	return nil
}

func (f *fakeKV) Delete(_ context.Context, key string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.store, key)
	return nil
}

func (f *fakeKV) DeleteMany(_ context.Context, keys []string) error {
	for _, k := range keys {
		delete(f.store, k)
	}
	return nil
}

func (f *fakeKV) Exists(_ context.Context, key string) (bool, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	_, ok := f.store[key]
	return ok, nil
}

func (f *fakeKV) ExistsMany(_ context.Context, keys []string) ([]bool, error) {
	res := make([]bool, len(keys))
	for i, k := range keys {
		ok, _ := f.Exists(context.Background(), k)
		res[i] = ok
	}
	return res, nil
}

func (f *fakeKV) GetKeysByPattern(_ context.Context, _ string) ([]string, error) {
	return nil, nil
}

type fixture struct {
	sm      *StateMachine
	repo    *fakeRepo
	emitter *fakeEmitter
	kv      *fakeKV
	channel uuid.UUID
}

func newFixture(t *testing.T) *fixture {
	t.Helper()
	repo := &fakeRepo{
		settings: model.ChannelDotaSettings{
			Mmr:       3000,
			MmrDelta:  25,
			ChannelID: uuid.New(),
		},
		updated: model.ChannelDotaSettings{
			Mmr:           3025,
			MmrDelta:      25,
			SessionWins:   1,
			SessionLosses: 0,
		},
	}
	emitter := &fakeEmitter{}
	store := newFakeKV()
	channelID := uuid.New()
	repo.settings.ChannelID = channelID
	sm := New(repo, emitter, store, slog.Default())
	return &fixture{
		sm:      sm,
		repo:    repo,
		emitter: emitter,
		kv:      store,
		channel: channelID,
	}
}

func inGamePayload(matchID int64, heroName string) gsi.Payload {
	return gsi.Payload{
		Map: &gsi.Map{
			MatchID:   matchID,
			GameState: gsi.GameStateInProgress,
			GameTime:  100,
		},
		Player: &gsi.Player{
			Activity:  gsi.PlayerActivityPlaying,
			TeamName:  "radiant",
			AccountID: 12345,
		},
		Hero: &gsi.Hero{Name: heroName},
	}
}

func TestIdleToInGameEmitsMatchStartedOnce(t *testing.T) {
	f := newFixture(t)
	ctx := context.Background()

	err := f.sm.Process(ctx, f.channel, inGamePayload(111, "npc_dota_hero_antimage"))
	require.NoError(t, err)

	require.Len(t, f.emitter.matchStarted, 1)
	require.Equal(t, "antimage", f.emitter.matchStarted[0].HeroName)
	require.Equal(t, f.channel.String(), f.emitter.matchStarted[0].ChannelID)

	snap, err := f.sm.GetSnapshot(ctx, f.channel)
	require.NoError(t, err)
	require.Equal(t, StateInGame, snap.State)
	require.True(t, snap.InGame)
	require.Equal(t, int64(111), snap.MatchID)
	require.Equal(t, "antimage", snap.HeroName)
	require.True(t, snap.IsRadiant)

	require.NotEmpty(t, f.emitter.stateUpdates)
	last := f.emitter.stateUpdates[len(f.emitter.stateUpdates)-1]
	require.True(t, last.InGame)
	require.Equal(t, "antimage", last.HeroName)
	require.Equal(t, int64(111), last.MatchID)
}

func TestSameMatchIDDedupesMatchStarted(t *testing.T) {
	f := newFixture(t)
	ctx := context.Background()

	payload := inGamePayload(222, "npc_dota_hero_antimage")
	require.NoError(t, f.sm.Process(ctx, f.channel, payload))
	payload.Map.GameTime = 200
	require.NoError(t, f.sm.Process(ctx, f.channel, payload))
	payload.Map.GameTime = 300
	require.NoError(t, f.sm.Process(ctx, f.channel, payload))

	require.Len(t, f.emitter.matchStarted, 1)
}

func TestNilMapKeepsIdleAndEmitsNothing(t *testing.T) {
	f := newFixture(t)
	ctx := context.Background()

	require.NoError(t, f.sm.Process(ctx, f.channel, gsi.Payload{}))

	require.Empty(t, f.emitter.matchStarted)
	require.Empty(t, f.emitter.matchEnded)
	require.Empty(t, f.emitter.roshanKilled)
	require.Empty(t, f.emitter.aegisPickup)

	snap, err := f.sm.GetSnapshot(ctx, f.channel)
	require.NoError(t, err)
	require.Equal(t, StateIdle, snap.State)
	require.False(t, snap.InGame)
}

func TestNilMapAfterGameTransitionsToIdle(t *testing.T) {
	f := newFixture(t)
	ctx := context.Background()

	require.NoError(t, f.sm.Process(ctx, f.channel, inGamePayload(333, "npc_dota_hero_pudge")))
	require.NoError(t, f.sm.Process(ctx, f.channel, gsi.Payload{}))

	snap, err := f.sm.GetSnapshot(ctx, f.channel)
	require.NoError(t, err)
	require.Equal(t, StateIdle, snap.State)
	require.False(t, snap.InGame)
}

func postGamePayload(matchID int64, winTeam gsi.WinTeam) gsi.Payload {
	return gsi.Payload{
		Map: &gsi.Map{
			MatchID:   matchID,
			GameState: gsi.GameStatePostGame,
			WinTeam:   winTeam,
		},
		Player: &gsi.Player{
			Activity:  gsi.PlayerActivityPlaying,
			TeamName:  "radiant",
			AccountID: 12345,
		},
	}
}

func TestPostGameWinUpdatesResultAndEndsMatch(t *testing.T) {
	f := newFixture(t)
	ctx := context.Background()

	require.NoError(t, f.sm.Process(ctx, f.channel, inGamePayload(444, "npc_dota_hero_pudge")))
	require.NoError(t, f.sm.Process(ctx, f.channel, postGamePayload(444, gsi.WinTeamRadiant)))

	require.Len(t, f.repo.updateCalls, 1)
	require.True(t, f.repo.updateCalls[0].won)
	require.Equal(t, 25, f.repo.updateCalls[0].mmrDelta)
	require.Equal(t, f.channel, f.repo.updateCalls[0].channelID)

	require.Len(t, f.emitter.matchEnded, 1)
	ended := f.emitter.matchEnded[0]
	require.True(t, ended.Win)
	require.Equal(t, "pudge", ended.HeroName)
	require.Equal(t, 3025, ended.Mmr)
	require.Equal(t, 1, ended.SessionWins)
	require.Equal(t, 0, ended.SessionLosses)

	snap, err := f.sm.GetSnapshot(ctx, f.channel)
	require.NoError(t, err)
	require.Equal(t, StateIdle, snap.State)
	require.False(t, snap.InGame)
	require.Equal(t, int64(0), snap.MatchID)
}

func TestPostGameLossUpdatesResult(t *testing.T) {
	f := newFixture(t)
	ctx := context.Background()

	f.repo.updated = model.ChannelDotaSettings{
		Mmr:           2975,
		MmrDelta:      25,
		SessionWins:   0,
		SessionLosses: 1,
	}

	require.NoError(t, f.sm.Process(ctx, f.channel, inGamePayload(555, "npc_dota_hero_pudge")))
	require.NoError(t, f.sm.Process(ctx, f.channel, postGamePayload(555, gsi.WinTeamDire)))

	require.Len(t, f.repo.updateCalls, 1)
	require.False(t, f.repo.updateCalls[0].won)
	require.Equal(t, -25, f.repo.updateCalls[0].mmrDelta)

	require.Len(t, f.emitter.matchEnded, 1)
	require.False(t, f.emitter.matchEnded[0].Win)
	require.Equal(t, 2975, f.emitter.matchEnded[0].Mmr)
	require.Equal(t, 1, f.emitter.matchEnded[0].SessionLosses)
}

func TestRoshanKilledEmittedOnce(t *testing.T) {
	f := newFixture(t)
	ctx := context.Background()

	payload := inGamePayload(666, "npc_dota_hero_pudge")
	payload.Events = []gsi.Event{
		{EventType: "roshan_killed", KillerTeam: "dire", GameTime: 500},
	}

	require.NoError(t, f.sm.Process(ctx, f.channel, payload))
	require.NoError(t, f.sm.Process(ctx, f.channel, payload))

	require.Len(t, f.emitter.roshanKilled, 1)
	require.Equal(t, "dire", f.emitter.roshanKilled[0].Team)
	require.Equal(t, 500, f.emitter.roshanKilled[0].GameTime)
	require.Equal(t, f.channel.String(), f.emitter.roshanKilled[0].ChannelID)
}

func TestAegisPickupEmittedOnce(t *testing.T) {
	f := newFixture(t)
	ctx := context.Background()

	slot := 3
	payload := inGamePayload(777, "npc_dota_hero_pudge")
	payload.Events = []gsi.Event{
		{EventType: "aegis_picked_up", Player: &slot, GameTime: 600},
	}

	require.NoError(t, f.sm.Process(ctx, f.channel, payload))
	require.NoError(t, f.sm.Process(ctx, f.channel, payload))

	require.Len(t, f.emitter.aegisPickup, 1)
	require.Equal(t, 600, f.emitter.aegisPickup[0].GameTime)
	require.Equal(t, f.channel.String(), f.emitter.aegisPickup[0].ChannelID)
}

func TestSnapshotPersistedAndRestored(t *testing.T) {
	f := newFixture(t)
	ctx := context.Background()

	payload := inGamePayload(888, "npc_dota_hero_antimage")
	payload.Events = []gsi.Event{
		{EventType: "roshan_killed", KillerTeam: "radiant", GameTime: 300},
	}
	require.NoError(t, f.sm.Process(ctx, f.channel, payload))
	require.Greater(t, f.kv.sets, 0)

	restored := New(f.repo, f.emitter, f.kv, slog.Default())

	require.NoError(t, restored.Process(ctx, f.channel, payload))
	require.Len(t, f.emitter.matchStarted, 1, "restored machine must not re-emit MatchStarted")
	require.Len(t, f.emitter.roshanKilled, 1, "restored machine must not re-emit roshan")

	snap, err := restored.GetSnapshot(ctx, f.channel)
	require.NoError(t, err)
	require.Equal(t, StateInGame, snap.State)
	require.Equal(t, int64(888), snap.MatchID)
	require.Equal(t, "antimage", snap.HeroName)
	require.True(t, snap.IsRadiant)
}

func TestRoshanEmitFailureRetriesOnNextPayload(t *testing.T) {
	f := newFixture(t)
	ctx := context.Background()

	payload := inGamePayload(101, "npc_dota_hero_pudge")
	payload.Events = []gsi.Event{
		{EventType: "roshan_killed", KillerTeam: "dire", GameTime: 700},
	}

	f.emitter.roshanErr = errors.New("publish failed")
	err := f.sm.Process(ctx, f.channel, payload)
	require.Error(t, err)
	require.Empty(t, f.emitter.roshanKilled)

	f.emitter.roshanErr = nil
	require.NoError(t, f.sm.Process(ctx, f.channel, payload))
	require.Len(t, f.emitter.roshanKilled, 1)
	require.Equal(t, "dire", f.emitter.roshanKilled[0].Team)
	require.Equal(t, 700, f.emitter.roshanKilled[0].GameTime)

	require.NoError(t, f.sm.Process(ctx, f.channel, payload))
	require.Len(t, f.emitter.roshanKilled, 1, "event must not be emitted twice after retry")
}

func TestFinishMatchBailsOnSettingsLoadFailure(t *testing.T) {
	f := newFixture(t)
	ctx := context.Background()

	require.NoError(t, f.sm.Process(ctx, f.channel, inGamePayload(202, "npc_dota_hero_pudge")))

	f.repo.getErr = errors.New("db down")
	err := f.sm.Process(ctx, f.channel, postGamePayload(202, gsi.WinTeamRadiant))
	require.Error(t, err)
	require.Empty(t, f.repo.updateCalls, "UpdateMatchResult must not be called when settings are unavailable")
	require.Empty(t, f.emitter.matchEnded)

	snap, err := f.sm.GetSnapshot(ctx, f.channel)
	require.NoError(t, err)
	require.Equal(t, int64(202), snap.MatchID, "match must stay tracked for retry")

	f.repo.getErr = nil
	require.NoError(t, f.sm.Process(ctx, f.channel, postGamePayload(202, gsi.WinTeamRadiant)))
	require.Len(t, f.repo.updateCalls, 1)
	require.True(t, f.repo.updateCalls[0].won)
	require.Equal(t, 25, f.repo.updateCalls[0].mmrDelta)
	require.Len(t, f.emitter.matchEnded, 1)

	snap, err = f.sm.GetSnapshot(ctx, f.channel)
	require.NoError(t, err)
	require.Equal(t, StateIdle, snap.State)
}

func TestGetSnapshotReturnsCurrentState(t *testing.T) {
	f := newFixture(t)
	ctx := context.Background()

	snap, err := f.sm.GetSnapshot(ctx, f.channel)
	require.NoError(t, err)
	require.Equal(t, StateIdle, snap.State)
	require.Equal(t, f.channel, snap.ChannelID)

	require.NoError(t, f.sm.Process(ctx, f.channel, inGamePayload(999, "npc_dota_hero_pudge")))

	snap, err = f.sm.GetSnapshot(ctx, f.channel)
	require.NoError(t, err)
	require.Equal(t, StateInGame, snap.State)
	require.Equal(t, int64(999), snap.MatchID)
}
