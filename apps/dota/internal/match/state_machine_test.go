package match

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"sync"
	"testing"
	"time"

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

func (f *fakeRepo) ApplyMatchResultOnce(
	_ context.Context,
	input dotarepository.ApplyMatchResultInput,
) (model.ChannelDotaSettings, error) {
	if f.updateErr != nil {
		return model.Nil, f.updateErr
	}
	f.updateCalls = append(f.updateCalls, updateMatchResultCall{
		channelID: input.ChannelID,
		won:       input.Won,
		mmrDelta:  input.MmrDelta,
	})
	return f.updated, nil
}

func (f *fakeRepo) GetMatchState(_ context.Context, _ uuid.UUID) (model.MatchState, error) {
	return model.MatchState{}, errors.New("not implemented")
}

func (f *fakeRepo) ApplyMatchStateTransition(
	_ context.Context,
	_ dotarepository.ApplyMatchStateTransitionInput,
) (bool, error) {
	return false, errors.New("not implemented")
}

func (f *fakeRepo) ClaimPredictionActions(
	_ context.Context,
	_ dotarepository.ClaimPredictionActionsInput,
) ([]model.ClaimedOutboxAction, error) {
	return nil, errors.New("not implemented")
}

func (f *fakeRepo) CompletePredictionAction(_ context.Context, _ uuid.UUID, _ uuid.UUID) error {
	return errors.New("not implemented")
}

func (f *fakeRepo) RetryPredictionAction(
	_ context.Context,
	_ uuid.UUID,
	_ uuid.UUID,
	_ time.Time,
) error {
	return errors.New("not implemented")
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
	mu             sync.Mutex
	roshanErr      error
	aegisErr       error
	abandonedErr   error
	matchStarted   []busdota.MatchStartedMessage
	matchEnded     []busdota.MatchEndedMessage
	matchAbandoned []busdota.MatchAbandonedMessage
	roshanKilled   []busdota.RoshanKilledMessage
	aegisPickup    []busdota.AegisPickupMessage
	stateUpdates   []busapi.DotaStateUpdateMessage
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

func (f *fakeEmitter) MatchAbandoned(_ context.Context, msg busdota.MatchAbandonedMessage) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.abandonedErr != nil {
		return f.abandonedErr
	}
	f.matchAbandoned = append(f.matchAbandoned, msg)
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
	require.Equal(t, int64(111), f.emitter.matchStarted[0].MatchID)
	require.True(t, f.emitter.matchStarted[0].TeamKnown)

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
	require.Equal(t, int64(444), ended.MatchID)

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

	playerID := 3
	payload := inGamePayload(777, "npc_dota_hero_pudge")
	payload.Events = []gsi.Event{
		{EventType: "aegis_picked_up", PlayerID: &playerID, GameTime: 600},
	}

	require.NoError(t, f.sm.Process(ctx, f.channel, payload))
	require.NoError(t, f.sm.Process(ctx, f.channel, payload))

	require.Len(t, f.emitter.aegisPickup, 1)
	require.Equal(t, 600, f.emitter.aegisPickup[0].GameTime)
	require.Equal(t, f.channel.String(), f.emitter.aegisPickup[0].ChannelID)
	require.NotNil(t, f.emitter.aegisPickup[0].PlayerID)
	require.Equal(t, playerID, *f.emitter.aegisPickup[0].PlayerID)
}

func TestAegisPickupUsesDecodedPlayerID(t *testing.T) {
	tests := []struct {
		name      string
		eventJSON string
		playerID  int
	}{
		{
			name:      "player_id",
			eventJSON: `{"event_type":"aegis_picked_up","player_id":2,"game_time":600}`,
			playerID:  2,
		},
		{
			name:      "legacy player",
			eventJSON: `{"event_type":"aegis_picked_up","player":3,"game_time":600}`,
			playerID:  3,
		},
		{
			name:      "player_id takes precedence",
			eventJSON: `{"event_type":"aegis_picked_up","player_id":4,"player":5,"game_time":600}`,
			playerID:  4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var event gsi.Event
			require.NoError(t, json.Unmarshal([]byte(tt.eventJSON), &event))

			f := newFixture(t)
			payload := inGamePayload(777, "npc_dota_hero_pudge")
			payload.Events = []gsi.Event{event}

			require.NoError(t, f.sm.Process(context.Background(), f.channel, payload))
			require.Len(t, f.emitter.aegisPickup, 1)
			require.NotNil(t, f.emitter.aegisPickup[0].PlayerID)
			require.Equal(t, tt.playerID, *f.emitter.aegisPickup[0].PlayerID)
		})
	}
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

func TestPostGameWithUnknownTeamDoesNotUpdateMatchResult(t *testing.T) {
	f := newFixture(t)
	ctx := context.Background()

	payload := inGamePayload(2_003, "npc_dota_hero_pudge")
	payload.Player.TeamName = "spectator"
	require.NoError(t, f.sm.Process(ctx, f.channel, payload))
	require.False(t, f.emitter.matchStarted[0].TeamKnown)

	postGame := postGamePayload(2_003, gsi.WinTeamDire)
	postGame.Player.TeamName = "spectator"
	require.NoError(t, f.sm.Process(ctx, f.channel, postGame))

	require.Empty(t, f.repo.updateCalls)
	require.Empty(t, f.emitter.matchEnded)
	snap, err := f.sm.GetSnapshot(ctx, f.channel)
	require.NoError(t, err)
	require.Equal(t, StatePostGame, snap.State)
	require.Equal(t, int64(2_003), snap.MatchID)
}

func TestLeavingTrackedMatchEmitsAbandonedBeforeClearingSnapshot(t *testing.T) {
	f := newFixture(t)
	ctx := context.Background()

	require.NoError(t, f.sm.Process(ctx, f.channel, inGamePayload(2_004, "npc_dota_hero_pudge")))
	require.NoError(t, f.sm.Process(ctx, f.channel, gsi.Payload{}))

	require.Equal(t, []busdota.MatchAbandonedMessage{{
		ChannelID: f.channel.String(),
		MatchID:   2_004,
	}}, f.emitter.matchAbandoned)
	snap, err := f.sm.GetSnapshot(ctx, f.channel)
	require.NoError(t, err)
	require.Equal(t, StateIdle, snap.State)
	require.Zero(t, snap.MatchID)
}

func TestAbandonedPublishFailurePreservesSnapshotForRetry(t *testing.T) {
	f := newFixture(t)
	ctx := context.Background()

	require.NoError(t, f.sm.Process(ctx, f.channel, inGamePayload(2_005, "npc_dota_hero_pudge")))
	f.emitter.abandonedErr = errors.New("publish failed")
	require.Error(t, f.sm.Process(ctx, f.channel, gsi.Payload{}))

	snap, err := f.sm.GetSnapshot(ctx, f.channel)
	require.NoError(t, err)
	require.Equal(t, int64(2_005), snap.MatchID)
	require.Equal(t, StateInGame, snap.State)

	f.emitter.abandonedErr = nil
	require.NoError(t, f.sm.Process(ctx, f.channel, gsi.Payload{}))
	require.Len(t, f.emitter.matchAbandoned, 1)
	require.Equal(t, int64(2_005), f.emitter.matchAbandoned[0].MatchID)
}

func TestReplacingTrackedMatchAbandonsOldMatchBeforeStartingNewOne(t *testing.T) {
	f := newFixture(t)
	ctx := context.Background()

	require.NoError(t, f.sm.Process(ctx, f.channel, inGamePayload(2_006, "npc_dota_hero_pudge")))
	require.NoError(t, f.sm.Process(ctx, f.channel, inGamePayload(2_007, "npc_dota_hero_axe")))

	require.Equal(t, []busdota.MatchAbandonedMessage{{
		ChannelID: f.channel.String(),
		MatchID:   2_006,
	}}, f.emitter.matchAbandoned)
	require.Len(t, f.emitter.matchStarted, 2)
	require.Equal(t, int64(2_007), f.emitter.matchStarted[1].MatchID)
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

func TestUpdateWinProbabilityPersistsAndEmitsStateUpdate(t *testing.T) {
	f := newFixture(t)
	ctx := context.Background()

	require.NoError(t, f.sm.Process(ctx, f.channel, inGamePayload(1001, "npc_dota_hero_pudge")))
	updatesBefore := len(f.emitter.stateUpdates)
	persistsBefore := f.kv.sets

	require.NoError(t, f.sm.UpdateWinProbability(ctx, f.channel, 1001, 0.625))

	snap, err := f.sm.GetSnapshot(ctx, f.channel)
	require.NoError(t, err)
	require.Equal(t, 0.625, snap.WinProbability)
	require.Equal(t, persistsBefore+1, f.kv.sets)
	require.Len(t, f.emitter.stateUpdates, updatesBefore+1)
	require.Equal(t, 0.625, f.emitter.stateUpdates[updatesBefore].WinProbability)

	restored := New(f.repo, f.emitter, f.kv, slog.Default())
	persisted, err := restored.GetSnapshot(ctx, f.channel)
	require.NoError(t, err)
	require.Equal(t, 0.625, persisted.WinProbability)
}

func TestUpdateWinProbabilityDiscardsStaleMatchResult(t *testing.T) {
	f := newFixture(t)
	ctx := context.Background()

	require.NoError(t, f.sm.Process(ctx, f.channel, inGamePayload(1006, "npc_dota_hero_pudge")))
	require.NoError(t, f.sm.Process(ctx, f.channel, inGamePayload(1007, "npc_dota_hero_pudge")))
	updatesBefore := len(f.emitter.stateUpdates)
	persistsBefore := f.kv.sets

	require.NoError(t, f.sm.UpdateWinProbability(ctx, f.channel, 1006, 0.625))

	snap, err := f.sm.GetSnapshot(ctx, f.channel)
	require.NoError(t, err)
	require.Equal(t, int64(1007), snap.MatchID)
	require.Zero(t, snap.WinProbability)
	require.Equal(t, persistsBefore, f.kv.sets)
	require.Len(t, f.emitter.stateUpdates, updatesBefore)
}

func TestUpdateWinProbabilityEmitsExactThresholdChange(t *testing.T) {
	f := newFixture(t)
	ctx := context.Background()

	require.NoError(t, f.sm.Process(ctx, f.channel, inGamePayload(1008, "npc_dota_hero_pudge")))
	require.NoError(t, f.sm.UpdateWinProbability(ctx, f.channel, 1008, 0.30))
	updatesBefore := len(f.emitter.stateUpdates)
	persistsBefore := f.kv.sets

	require.NoError(t, f.sm.UpdateWinProbability(ctx, f.channel, 1008, 0.35))

	snap, err := f.sm.GetSnapshot(ctx, f.channel)
	require.NoError(t, err)
	require.Equal(t, 0.35, snap.WinProbability)
	require.Equal(t, persistsBefore+1, f.kv.sets)
	require.Len(t, f.emitter.stateUpdates, updatesBefore+1)
	require.Equal(t, 0.35, f.emitter.stateUpdates[updatesBefore].WinProbability)
}

func TestUpdateWinProbabilityDoesNotEmitBelowThreshold(t *testing.T) {
	f := newFixture(t)
	ctx := context.Background()

	require.NoError(t, f.sm.Process(ctx, f.channel, inGamePayload(1009, "npc_dota_hero_pudge")))
	require.NoError(t, f.sm.UpdateWinProbability(ctx, f.channel, 1009, 0.30))
	updatesBefore := len(f.emitter.stateUpdates)
	persistsBefore := f.kv.sets

	require.NoError(t, f.sm.UpdateWinProbability(ctx, f.channel, 1009, 0.349))

	snap, err := f.sm.GetSnapshot(ctx, f.channel)
	require.NoError(t, err)
	require.Equal(t, 0.30, snap.WinProbability)
	require.Equal(t, persistsBefore, f.kv.sets)
	require.Len(t, f.emitter.stateUpdates, updatesBefore)
}

func TestUpdateWinProbabilityDoesNotEmitFor499BasisPointChange(t *testing.T) {
	f := newFixture(t)
	ctx := context.Background()
	initialProbability := math.Nextafter(0.29995, math.Inf(1))
	nextProbability := math.Nextafter(0.34995, 0)

	require.NoError(t, f.sm.Process(ctx, f.channel, inGamePayload(1010, "npc_dota_hero_pudge")))
	require.NoError(t, f.sm.UpdateWinProbability(ctx, f.channel, 1010, initialProbability))
	updatesBefore := len(f.emitter.stateUpdates)
	persistsBefore := f.kv.sets

	require.NoError(t, f.sm.UpdateWinProbability(ctx, f.channel, 1010, nextProbability))

	snap, err := f.sm.GetSnapshot(ctx, f.channel)
	require.NoError(t, err)
	require.Equal(t, initialProbability, snap.WinProbability)
	require.Equal(t, persistsBefore, f.kv.sets)
	require.Len(t, f.emitter.stateUpdates, updatesBefore)
}

func TestUpdateWinProbabilityThrottlesSmallChanges(t *testing.T) {
	f := newFixture(t)
	ctx := context.Background()

	require.NoError(t, f.sm.Process(ctx, f.channel, inGamePayload(1002, "npc_dota_hero_pudge")))
	require.NoError(t, f.sm.UpdateWinProbability(ctx, f.channel, 1002, 0.625))
	updatesBefore := len(f.emitter.stateUpdates)
	persistsBefore := f.kv.sets

	require.NoError(t, f.sm.UpdateWinProbability(ctx, f.channel, 1002, 0.65))

	snap, err := f.sm.GetSnapshot(ctx, f.channel)
	require.NoError(t, err)
	require.Equal(t, 0.625, snap.WinProbability)
	require.Equal(t, persistsBefore, f.kv.sets)
	require.Len(t, f.emitter.stateUpdates, updatesBefore)
}

func TestUpdateWinProbabilityRejectsInvalidValues(t *testing.T) {
	for _, test := range []struct {
		name        string
		probability float64
	}{
		{name: "negative", probability: -0.01},
		{name: "greater_than_one", probability: 1.01},
		{name: "nan", probability: math.NaN()},
		{name: "positive_infinity", probability: math.Inf(1)},
		{name: "negative_infinity", probability: math.Inf(-1)},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := newFixture(t)
			ctx := context.Background()

			require.NoError(t, f.sm.Process(ctx, f.channel, inGamePayload(1003, "npc_dota_hero_pudge")))
			before, err := f.sm.GetSnapshot(ctx, f.channel)
			require.NoError(t, err)
			updatesBefore := len(f.emitter.stateUpdates)
			persistsBefore := f.kv.sets

			err = f.sm.UpdateWinProbability(ctx, f.channel, 1003, test.probability)
			require.Error(t, err)
			require.ErrorContains(t, err, "win probability")

			after, err := f.sm.GetSnapshot(ctx, f.channel)
			require.NoError(t, err)
			require.Equal(t, before, after)
			require.Equal(t, persistsBefore, f.kv.sets)
			require.Len(t, f.emitter.stateUpdates, updatesBefore)
		})
	}
}

func TestUpdateWinProbabilityRejectsInvalidValuesWhileIdle(t *testing.T) {
	for _, test := range []struct {
		name        string
		probability float64
	}{
		{name: "negative", probability: -0.01},
		{name: "greater_than_one", probability: 1.01},
		{name: "nan", probability: math.NaN()},
		{name: "positive_infinity", probability: math.Inf(1)},
		{name: "negative_infinity", probability: math.Inf(-1)},
	} {
		t.Run(test.name, func(t *testing.T) {
			f := newFixture(t)
			ctx := context.Background()
			before, err := f.sm.GetSnapshot(ctx, f.channel)
			require.NoError(t, err)
			updatesBefore := len(f.emitter.stateUpdates)
			persistsBefore := f.kv.sets

			err = f.sm.UpdateWinProbability(ctx, f.channel, 1003, test.probability)
			require.Error(t, err)
			require.ErrorContains(t, err, "win probability")

			after, err := f.sm.GetSnapshot(ctx, f.channel)
			require.NoError(t, err)
			require.Equal(t, before, after)
			require.Equal(t, persistsBefore, f.kv.sets)
			require.Len(t, f.emitter.stateUpdates, updatesBefore)
		})
	}
}

func TestUpdateWinProbabilityDoesNothingWhileIdle(t *testing.T) {
	f := newFixture(t)
	ctx := context.Background()

	require.NoError(t, f.sm.UpdateWinProbability(ctx, f.channel, 1003, 0.625))

	snap, err := f.sm.GetSnapshot(ctx, f.channel)
	require.NoError(t, err)
	require.Equal(t, StateIdle, snap.State)
	require.Zero(t, snap.WinProbability)
	require.Zero(t, f.kv.sets)
	require.Empty(t, f.emitter.stateUpdates)
}

func TestLeavingGameClearsWinProbability(t *testing.T) {
	f := newFixture(t)
	ctx := context.Background()

	require.NoError(t, f.sm.Process(ctx, f.channel, inGamePayload(1004, "npc_dota_hero_pudge")))
	require.NoError(t, f.sm.UpdateWinProbability(ctx, f.channel, 1004, 0.625))
	require.NoError(t, f.sm.Process(ctx, f.channel, gsi.Payload{}))

	snap, err := f.sm.GetSnapshot(ctx, f.channel)
	require.NoError(t, err)
	require.Zero(t, snap.WinProbability)
	require.Equal(t, 0.0, f.emitter.stateUpdates[len(f.emitter.stateUpdates)-1].WinProbability)
}

func TestFinishingMatchClearsWinProbability(t *testing.T) {
	f := newFixture(t)
	ctx := context.Background()

	require.NoError(t, f.sm.Process(ctx, f.channel, inGamePayload(1005, "npc_dota_hero_pudge")))
	require.NoError(t, f.sm.UpdateWinProbability(ctx, f.channel, 1005, 0.625))
	updatesBefore := len(f.emitter.stateUpdates)

	require.NoError(t, f.sm.Process(ctx, f.channel, postGamePayload(1005, gsi.WinTeamRadiant)))

	final, err := f.sm.GetSnapshot(ctx, f.channel)
	require.NoError(t, err)
	require.Equal(t, StateIdle, final.State)
	require.Zero(t, final.WinProbability)
	require.Len(t, f.emitter.stateUpdates, updatesBefore+1)
	require.Zero(t, f.emitter.stateUpdates[updatesBefore].WinProbability)

	restored := New(f.repo, f.emitter, f.kv, slog.Default())
	next, err := restored.GetSnapshot(ctx, f.channel)
	require.NoError(t, err)
	require.Equal(t, StateIdle, next.State)
	require.Zero(t, next.WinProbability)
}
