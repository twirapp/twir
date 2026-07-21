package match

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"math"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/twirapp/twir/apps/dota/internal/gsi"
	busapi "github.com/twirapp/twir/libs/bus-core/api"
	busdota "github.com/twirapp/twir/libs/bus-core/dota"
	dotarepository "github.com/twirapp/twir/libs/repositories/dota"
	"github.com/twirapp/twir/libs/repositories/dota/model"
)

type fakeLifecycleRepository struct {
	mu sync.Mutex

	settings model.ChannelDotaSettings
	states   map[uuid.UUID]model.MatchState
	actions  []model.OutboxActionInput

	updateMatchResultCalls    int
	applyMatchResultOnceCalls int
	applyMatchStateCalls      int
	beforeApply               func()
	applyBarrier              *applyBarrier
}

type applyBarrier struct {
	started chan<- struct{}
	release <-chan struct{}
}

func newFakeLifecycleRepository(channelID uuid.UUID) *fakeLifecycleRepository {
	return &fakeLifecycleRepository{
		settings: model.ChannelDotaSettings{
			ChannelID: channelID,
			Mmr:       3000,
			MmrDelta:  25,
		},
		states: make(map[uuid.UUID]model.MatchState),
	}
}

func (r *fakeLifecycleRepository) GetByChannelID(
	_ context.Context,
	_ uuid.UUID,
) (model.ChannelDotaSettings, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.settings, nil
}

func (r *fakeLifecycleRepository) GetByGsiToken(
	_ context.Context,
	_ string,
) (model.ChannelDotaSettings, error) {
	return model.Nil, dotarepository.ErrNotFound
}

func (r *fakeLifecycleRepository) Create(
	_ context.Context,
	_ dotarepository.CreateInput,
) (model.ChannelDotaSettings, error) {
	return model.Nil, errors.New("not implemented")
}

func (r *fakeLifecycleRepository) Update(
	_ context.Context,
	_ uuid.UUID,
	_ dotarepository.UpdateInput,
) (model.ChannelDotaSettings, error) {
	return model.Nil, errors.New("not implemented")
}

func (r *fakeLifecycleRepository) UpdateMatchResult(
	_ context.Context,
	_ uuid.UUID,
	_ bool,
	_ int,
) (model.ChannelDotaSettings, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.updateMatchResultCalls++
	return r.settings, nil
}

func (r *fakeLifecycleRepository) ApplyMatchResultOnce(
	_ context.Context,
	_ dotarepository.ApplyMatchResultInput,
) (model.ChannelDotaSettings, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.applyMatchResultOnceCalls++
	return r.settings, nil
}

func (r *fakeLifecycleRepository) GetMatchState(
	_ context.Context,
	channelID uuid.UUID,
) (model.MatchState, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	state, ok := r.states[channelID]
	if !ok {
		return model.MatchState{
			ChannelID: channelID,
			Snapshot:  json.RawMessage(`{}`),
		}, nil
	}
	state.Snapshot = append(json.RawMessage(nil), state.Snapshot...)
	return state, nil
}

func (r *fakeLifecycleRepository) ApplyMatchStateTransition(
	_ context.Context,
	input dotarepository.ApplyMatchStateTransitionInput,
) (bool, error) {
	r.mu.Lock()
	hook := r.beforeApply
	r.beforeApply = nil
	barrier := r.applyBarrier
	r.mu.Unlock()
	if hook != nil {
		hook()
	}
	if barrier != nil {
		barrier.started <- struct{}{}
		<-barrier.release
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.applyMatchStateCalls++

	current, ok := r.states[input.ChannelID]
	if !ok {
		current = model.MatchState{ChannelID: input.ChannelID, Snapshot: json.RawMessage(`{}`)}
	}
	if current.Revision != input.ExpectedRevision {
		return false, nil
	}

	r.states[input.ChannelID] = model.MatchState{
		ChannelID:         input.ChannelID,
		Revision:          current.Revision + 1,
		ProviderTimestamp: input.ProviderTimestamp,
		Snapshot:          append(json.RawMessage(nil), input.Snapshot...),
	}
	for _, action := range input.Actions {
		action.Payload = append(json.RawMessage(nil), action.Payload...)
		r.actions = append(r.actions, action)
	}

	return true, nil
}

func (r *fakeLifecycleRepository) ClaimPredictionActions(
	_ context.Context,
	_ dotarepository.ClaimPredictionActionsInput,
) ([]model.ClaimedOutboxAction, error) {
	return nil, errors.New("not implemented")
}

func (r *fakeLifecycleRepository) RenewPredictionAction(
	_ context.Context,
	_ uuid.UUID,
	_ uuid.UUID,
	_ time.Duration,
) error {
	return errors.New("not implemented")
}

func (r *fakeLifecycleRepository) CompletePredictionAction(
	_ context.Context,
	_ uuid.UUID,
	_ uuid.UUID,
) error {
	return errors.New("not implemented")
}

func (r *fakeLifecycleRepository) RetryPredictionAction(
	_ context.Context,
	_ uuid.UUID,
	_ uuid.UUID,
	_ time.Time,
) error {
	return errors.New("not implemented")
}

func (r *fakeLifecycleRepository) ResetSession(
	_ context.Context,
	_ uuid.UUID,
) (model.ChannelDotaSettings, error) {
	return model.Nil, errors.New("not implemented")
}

func (r *fakeLifecycleRepository) RegenerateGsiToken(
	_ context.Context,
	_ uuid.UUID,
) (model.ChannelDotaSettings, error) {
	return model.Nil, errors.New("not implemented")
}

func (r *fakeLifecycleRepository) snapshot(t testing.TB, channelID uuid.UUID) Snapshot {
	t.Helper()

	state, err := r.GetMatchState(context.Background(), channelID)
	require.NoError(t, err)

	snapshot := Snapshot{ChannelID: channelID, State: StateIdle}
	if len(state.Snapshot) != 0 && string(state.Snapshot) != "{}" {
		require.NoError(t, json.Unmarshal(state.Snapshot, &snapshot))
	}
	if snapshot.ChannelID == uuid.Nil {
		snapshot.ChannelID = channelID
	}
	snapshot.Revision = uint64(state.Revision)
	snapshot.LastProviderTimestamp = state.ProviderTimestamp
	return snapshot
}

func (r *fakeLifecycleRepository) lifecycleActions(t testing.TB) []LifecycleAction {
	t.Helper()

	r.mu.Lock()
	actions := append([]model.OutboxActionInput(nil), r.actions...)
	r.mu.Unlock()

	result := make([]LifecycleAction, 0, len(actions))
	for _, action := range actions {
		var payload LifecycleAction
		require.NoError(t, json.Unmarshal(action.Payload, &payload))
		result = append(result, payload)
	}
	return result
}

func (r *fakeLifecycleRepository) actionKinds() []model.OutboxAction {
	r.mu.Lock()
	defer r.mu.Unlock()

	result := make([]model.OutboxAction, 0, len(r.actions))
	for _, action := range r.actions {
		result = append(result, action.Action)
	}
	return result
}

func (r *fakeLifecycleRepository) actionCount() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.actions)
}

func (r *fakeLifecycleRepository) settlementCalls() (int, int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.updateMatchResultCalls, r.applyMatchResultOnceCalls
}

func (r *fakeLifecycleRepository) setBeforeApply(hook func()) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.beforeApply = hook
}

func (r *fakeLifecycleRepository) setApplyBarrier(barrier *applyBarrier) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.applyBarrier = barrier
}

func (r *fakeLifecycleRepository) setEventClaim(
	t testing.TB,
	channelID uuid.UUID,
	key string,
	token string,
	expiresAt int64,
) {
	t.Helper()

	r.mu.Lock()
	defer r.mu.Unlock()
	state, ok := r.states[channelID]
	require.True(t, ok)

	var snapshot map[string]any
	require.NoError(t, json.Unmarshal(state.Snapshot, &snapshot))
	snapshot["eventClaims"] = map[string]any{
		key: map[string]any{
			"token":     token,
			"expiresAt": expiresAt,
		},
	}
	data, err := json.Marshal(snapshot)
	require.NoError(t, err)
	state.Snapshot = data
	r.states[channelID] = state
}

type fakeEmitter struct {
	mu sync.Mutex

	onMatchStarted func()
	onRoshanKilled func()
	onAegisPickup  func()
	roshanErr      error
	aegisErr       error
	stateUpdateErr error
	roshanStarted  chan<- struct{}
	roshanRelease  <-chan struct{}

	matchStarted        []busdota.MatchStartedMessage
	matchEnded          []busdota.MatchEndedMessage
	matchAbandoned      []busdota.MatchAbandonedMessage
	roshanKilled        []busdota.RoshanKilledMessage
	aegisPickup         []busdota.AegisPickupMessage
	stateUpdateAttempts []busapi.DotaStateUpdateMessage
	stateUpdates        []busapi.DotaStateUpdateMessage
}

func (e *fakeEmitter) MatchStarted(_ context.Context, message busdota.MatchStartedMessage) error {
	e.mu.Lock()
	e.matchStarted = append(e.matchStarted, message)
	hook := e.onMatchStarted
	e.mu.Unlock()
	if hook != nil {
		hook()
	}
	return nil
}

func (e *fakeEmitter) MatchEnded(_ context.Context, message busdota.MatchEndedMessage) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.matchEnded = append(e.matchEnded, message)
	return nil
}

func (e *fakeEmitter) MatchAbandoned(_ context.Context, message busdota.MatchAbandonedMessage) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.matchAbandoned = append(e.matchAbandoned, message)
	return nil
}

func (e *fakeEmitter) RoshanKilled(_ context.Context, message busdota.RoshanKilledMessage) error {
	e.mu.Lock()
	err := e.roshanErr
	var started chan<- struct{}
	var release <-chan struct{}
	if err == nil {
		e.roshanKilled = append(e.roshanKilled, message)
		started = e.roshanStarted
		release = e.roshanRelease
		e.roshanStarted = nil
		e.roshanRelease = nil
	}
	hook := e.onRoshanKilled
	e.mu.Unlock()
	if err != nil {
		return err
	}
	if started != nil {
		started <- struct{}{}
		<-release
	}
	if hook != nil {
		hook()
	}
	return nil
}

func (e *fakeEmitter) blockNextRoshan(started chan<- struct{}, release <-chan struct{}) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.roshanStarted = started
	e.roshanRelease = release
}

func (e *fakeEmitter) AegisPickup(_ context.Context, message busdota.AegisPickupMessage) error {
	e.mu.Lock()
	err := e.aegisErr
	if err == nil {
		e.aegisPickup = append(e.aegisPickup, message)
	}
	hook := e.onAegisPickup
	e.mu.Unlock()
	if err != nil {
		return err
	}
	if hook != nil {
		hook()
	}
	return nil
}

func (e *fakeEmitter) StateUpdate(_ context.Context, message busapi.DotaStateUpdateMessage) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.stateUpdateAttempts = append(e.stateUpdateAttempts, message)
	if e.stateUpdateErr != nil {
		err := e.stateUpdateErr
		e.stateUpdateErr = nil
		return err
	}
	e.stateUpdates = append(e.stateUpdates, message)
	return nil
}

func newFixture() (*StateMachine, *fakeLifecycleRepository, *fakeEmitter, uuid.UUID) {
	channelID := uuid.New()
	repository := newFakeLifecycleRepository(channelID)
	emitter := &fakeEmitter{}
	return New(repository, emitter, slog.Default()), repository, emitter, channelID
}

func inGamePayloadAt(matchID int64, providerTimestamp int64) gsi.Payload {
	return inGamePayloadWithGameTime(matchID, providerTimestamp, 100)
}

func inGamePayloadWithGameTime(matchID int64, providerTimestamp int64, gameTime int) gsi.Payload {
	return gsi.Payload{
		Provider: gsi.Provider{Timestamp: providerTimestamp},
		Map: &gsi.Map{
			MatchID:   matchID,
			GameState: gsi.GameStateInProgress,
			GameTime:  gameTime,
		},
		Player: &gsi.Player{
			Activity:  gsi.PlayerActivityPlaying,
			TeamName:  "radiant",
			AccountID: 12345,
		},
		Hero: &gsi.Hero{Name: "npc_dota_hero_axe"},
	}
}

func postGamePayloadAt(
	matchID int64,
	providerTimestamp int64,
	winTeam gsi.WinTeam,
) gsi.Payload {
	return gsi.Payload{
		Provider: gsi.Provider{Timestamp: providerTimestamp},
		Map: &gsi.Map{
			MatchID:   matchID,
			GameState: gsi.GameStatePostGame,
			GameTime:  300,
			WinTeam:   winTeam,
		},
		Player: &gsi.Player{
			Activity:  gsi.PlayerActivityPlaying,
			TeamName:  "radiant",
			AccountID: 12345,
		},
	}
}

func TestInGameStartCommitsOutboxBeforeMatchStarted(t *testing.T) {
	sm, repository, emitter, channelID := newFixture()
	startedAfterCommit := false
	emitter.onMatchStarted = func() {
		startedAfterCommit = repository.actionCount() == 1
	}

	require.NoError(t, sm.Process(context.Background(), channelID, inGamePayloadAt(101, 100)))

	snapshot := repository.snapshot(t, channelID)
	require.Equal(t, StateInGame, snapshot.State)
	require.Equal(t, int64(101), snapshot.MatchID)
	require.Equal(t, uint64(1), snapshot.Revision)
	require.True(t, startedAfterCommit)
	require.Equal(t, []model.OutboxAction{model.OutboxActionCreate}, repository.actionKinds())
	require.Equal(t, []LifecycleAction{{
		Kind:           ActionCreate,
		ChannelID:      channelID,
		MatchID:        101,
		SteamAccountID: "12345",
		HeroName:       "axe",
		TeamKnown:      true,
	}}, repository.lifecycleActions(t))
	require.Equal(t, []busdota.MatchStartedMessage{{
		ChannelID:      channelID.String(),
		SteamAccountID: "12345",
		HeroName:       "axe",
		MatchID:        101,
		TeamKnown:      true,
	}}, emitter.matchStarted)
	require.Len(t, emitter.stateUpdates, 1)
}

func TestInGameStartPersistsUnknownTeamForCreateAction(t *testing.T) {
	sm, repository, _, channelID := newFixture()
	payload := inGamePayloadAt(1_001, 100)
	payload.Player.TeamName = ""

	require.NoError(t, sm.Process(context.Background(), channelID, payload))

	actions := repository.lifecycleActions(t)
	require.Len(t, actions, 1)
	require.False(t, actions[0].TeamKnown)

	repository.mu.Lock()
	storedPayload := append([]byte(nil), repository.actions[0].Payload...)
	repository.mu.Unlock()
	var rawPayload map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(storedPayload, &rawPayload))
	require.Equal(t, json.RawMessage("false"), rawPayload["teamKnown"])
}

func TestPostGamePersistsResolveWithoutSettlementOrMatchEnded(t *testing.T) {
	sm, repository, emitter, channelID := newFixture()

	require.NoError(t, sm.Process(context.Background(), channelID, inGamePayloadAt(102, 100)))
	require.NoError(
		t,
		sm.Process(context.Background(), channelID, postGamePayloadAt(102, 200, gsi.WinTeamRadiant)),
	)

	snapshot := repository.snapshot(t, channelID)
	require.Equal(t, StateIdle, snapshot.State)
	require.Zero(t, snapshot.MatchID)
	require.Equal(t, int64(200), snapshot.LastProviderTimestamp)
	require.Equal(t, []model.OutboxAction{
		model.OutboxActionCreate,
		model.OutboxActionResolve,
	}, repository.actionKinds())
	actions := repository.lifecycleActions(t)
	require.True(t, actions[1].Win)
	require.Equal(t, "axe", actions[1].HeroName)
	require.Equal(t, "12345", actions[1].SteamAccountID)
	updated, applied := repository.settlementCalls()
	require.Zero(t, updated)
	require.Zero(t, applied)
	require.Empty(t, emitter.matchEnded)
}

func TestDelayedPostGameCannotSettleNewerMatch(t *testing.T) {
	sm, repository, emitter, channelID := newFixture()

	require.NoError(t, sm.Process(context.Background(), channelID, inGamePayloadAt(2006, 100)))
	require.NoError(t, sm.Process(context.Background(), channelID, inGamePayloadAt(2007, 200)))
	require.NoError(
		t,
		sm.Process(context.Background(), channelID, postGamePayloadAt(2006, 150, gsi.WinTeamRadiant)),
	)

	snapshot := repository.snapshot(t, channelID)
	require.Equal(t, int64(2007), snapshot.MatchID)
	require.Equal(t, []model.OutboxAction{
		model.OutboxActionCreate,
		model.OutboxActionCancel,
		model.OutboxActionCreate,
	}, repository.actionKinds())
	for _, action := range repository.lifecycleActions(t) {
		require.NotEqual(t, ActionResolve, action.Kind)
	}
	require.Empty(t, emitter.matchEnded)
}

func TestDelayedInactivePostGameCannotCancelNewerMatch(t *testing.T) {
	for _, test := range []struct {
		name   string
		player *gsi.Player
	}{
		{
			name: "missing player",
		},
		{
			name: "non-playing player",
			player: &gsi.Player{
				Activity: gsi.PlayerActivity("spectating"),
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			sm, repository, emitter, channelID := newFixture()

			require.NoError(t, sm.Process(context.Background(), channelID, inGamePayloadAt(2010, 100)))
			require.NoError(t, sm.Process(context.Background(), channelID, inGamePayloadAt(2011, 200)))
			delayed := postGamePayloadAt(2010, 300, gsi.WinTeamRadiant)
			delayed.Player = test.player
			require.NoError(t, sm.Process(context.Background(), channelID, delayed))

			snapshot := repository.snapshot(t, channelID)
			require.Equal(t, int64(2011), snapshot.MatchID)
			require.Equal(t, []model.OutboxAction{
				model.OutboxActionCreate,
				model.OutboxActionCancel,
				model.OutboxActionCreate,
			}, repository.actionKinds())
			require.Equal(t, []busdota.MatchAbandonedMessage{{
				ChannelID: channelID.String(),
				MatchID:   2010,
			}}, emitter.matchAbandoned)
			updated, applied := repository.settlementCalls()
			require.Zero(t, updated)
			require.Zero(t, applied)
			require.Empty(t, emitter.matchEnded)
		})
	}
}

func TestConcurrentTerminalInputsPersistOneAction(t *testing.T) {
	first, repository, emitter, channelID := newFixture()
	second := New(repository, emitter, slog.Default())

	require.NoError(t, first.Process(context.Background(), channelID, inGamePayloadAt(2007, 100)))
	started := make(chan struct{}, 2)
	release := make(chan struct{})
	repository.setApplyBarrier(&applyBarrier{started: started, release: release})

	var wg sync.WaitGroup
	errs := make(chan error, 2)
	for _, machine := range []*StateMachine{first, second} {
		wg.Add(1)
		go func(machine *StateMachine) {
			defer wg.Done()
			errs <- machine.Process(
				context.Background(),
				channelID,
				postGamePayloadAt(2007, 200, gsi.WinTeamRadiant),
			)
		}(machine)
	}
	<-started
	<-started
	close(release)
	wg.Wait()
	close(errs)
	for err := range errs {
		require.NoError(t, err)
	}

	resolveActions := 0
	for _, action := range repository.lifecycleActions(t) {
		if action.Kind == ActionResolve {
			resolveActions++
		}
	}
	require.Equal(t, 1, resolveActions)
}

func TestConcurrentResolveAndInactiveInputPersistOneTerminalAction(t *testing.T) {
	resolveMachine, repository, emitter, channelID := newFixture()
	cancelMachine := New(repository, emitter, slog.Default())

	require.NoError(t, resolveMachine.Process(context.Background(), channelID, inGamePayloadAt(2012, 100)))
	started := make(chan struct{}, 3)
	release := make(chan struct{})
	repository.setApplyBarrier(&applyBarrier{started: started, release: release})

	var wg sync.WaitGroup
	errs := make(chan error, 2)
	wg.Add(2)
	go func() {
		defer wg.Done()
		errs <- resolveMachine.Process(
			context.Background(),
			channelID,
			postGamePayloadAt(2012, 200, gsi.WinTeamRadiant),
		)
	}()
	go func() {
		defer wg.Done()
		errs <- cancelMachine.Process(context.Background(), channelID, gsi.Payload{
			Provider: gsi.Provider{Timestamp: 201},
		})
	}()
	<-started
	<-started
	close(release)
	wg.Wait()
	close(errs)
	for err := range errs {
		require.NoError(t, err)
	}

	terminalActions := make([]LifecycleAction, 0, 1)
	for _, action := range repository.lifecycleActions(t) {
		if action.MatchID != 2012 {
			continue
		}
		if action.Kind == ActionResolve || action.Kind == ActionCancel {
			terminalActions = append(terminalActions, action)
		}
	}
	require.Len(t, terminalActions, 1)
	require.Contains(t, []ActionKind{ActionResolve, ActionCancel}, terminalActions[0].Kind)
	updated, applied := repository.settlementCalls()
	require.Zero(t, updated)
	require.Zero(t, applied)
	require.Empty(t, emitter.matchEnded)
}

func TestMaplessInputCancelsOnlyWhenProviderTimestampIsStrictlyNewer(t *testing.T) {
	sm, repository, _, channelID := newFixture()

	require.NoError(t, sm.Process(context.Background(), channelID, inGamePayloadAt(103, 100)))
	require.NoError(t, sm.Process(context.Background(), channelID, gsi.Payload{
		Provider: gsi.Provider{Timestamp: 100},
	}))
	require.Equal(t, int64(103), repository.snapshot(t, channelID).MatchID)
	require.Equal(t, []model.OutboxAction{model.OutboxActionCreate}, repository.actionKinds())

	require.NoError(t, sm.Process(context.Background(), channelID, gsi.Payload{
		Provider: gsi.Provider{Timestamp: 101},
	}))
	require.Equal(t, StateIdle, repository.snapshot(t, channelID).State)
	require.Equal(t, []model.OutboxAction{
		model.OutboxActionCreate,
		model.OutboxActionCancel,
	}, repository.actionKinds())
}

func TestSourceOrderingRejectsStaleAndEqualReplacementPayloads(t *testing.T) {
	sm, repository, _, channelID := newFixture()

	require.NoError(t, sm.Process(context.Background(), channelID, inGamePayloadWithGameTime(104, 100, 100)))
	require.NoError(t, sm.Process(context.Background(), channelID, inGamePayloadWithGameTime(104, 99, 999)))
	require.Equal(t, 100, repository.snapshot(t, channelID).GameTime)

	require.NoError(t, sm.Process(context.Background(), channelID, inGamePayloadWithGameTime(104, 100, 100)))
	require.Equal(t, 100, repository.snapshot(t, channelID).GameTime)

	require.NoError(t, sm.Process(context.Background(), channelID, inGamePayloadWithGameTime(104, 100, 101)))
	require.Equal(t, 101, repository.snapshot(t, channelID).GameTime)

	require.NoError(t, sm.Process(context.Background(), channelID, inGamePayloadWithGameTime(105, 100, 200)))
	require.Equal(t, int64(104), repository.snapshot(t, channelID).MatchID)

	require.NoError(t, sm.Process(context.Background(), channelID, inGamePayloadWithGameTime(105, 101, 200)))
	require.Equal(t, int64(105), repository.snapshot(t, channelID).MatchID)
	require.Equal(t, []model.OutboxAction{
		model.OutboxActionCreate,
		model.OutboxActionCancel,
		model.OutboxActionCreate,
	}, repository.actionKinds())
}

func TestReplacementTransitionUsesDistinctChannelGlobalActionSequences(t *testing.T) {
	sm, repository, _, channelID := newFixture()

	require.NoError(t, sm.Process(context.Background(), channelID, inGamePayloadAt(119, 100)))
	require.NoError(t, sm.Process(context.Background(), channelID, inGamePayloadAt(120, 200)))

	repository.mu.Lock()
	actions := append([]model.OutboxActionInput(nil), repository.actions...)
	repository.mu.Unlock()
	require.Len(t, actions, 3)
	require.Equal(t, []model.OutboxAction{
		model.OutboxActionCreate,
		model.OutboxActionCancel,
		model.OutboxActionCreate,
	}, []model.OutboxAction{actions[0].Action, actions[1].Action, actions[2].Action})
	require.Equal(t, []int64{1, 3, 4}, []int64{actions[0].Sequence, actions[1].Sequence, actions[2].Sequence})
}

func TestTransitionInputRejectsOutboxSequenceOverflow(t *testing.T) {
	channelID := uuid.New()
	maxInt64 := uint64(^uint64(0) >> 1)
	overflowRevision := maxInt64/2 + 2
	state := model.MatchState{
		ChannelID: channelID,
		Revision:  int64(overflowRevision - 1),
		Snapshot:  json.RawMessage(`{}`),
	}
	next := Snapshot{
		ChannelID:             channelID,
		Revision:              overflowRevision,
		State:                 StateInGame,
		InGame:                true,
		MatchID:               122,
		LastProviderTimestamp: 100,
	}

	_, err := transitionInput(state, next, []LifecycleAction{{
		Kind:      ActionCreate,
		ChannelID: channelID,
		MatchID:   122,
	}})

	require.Error(t, err)
}

func TestUpdateWinProbabilityPreservesCurrentMatchSourceFields(t *testing.T) {
	sm, repository, emitter, channelID := newFixture()

	require.NoError(t, sm.Process(context.Background(), channelID, inGamePayloadAt(106, 100)))
	updatesBefore := len(emitter.stateUpdates)
	require.NoError(t, sm.UpdateWinProbability(context.Background(), channelID, 106, 0.625))

	snapshot := repository.snapshot(t, channelID)
	require.Equal(t, int64(106), snapshot.MatchID)
	require.Equal(t, int64(100), snapshot.LastProviderTimestamp)
	require.Equal(t, 100, snapshot.LastGameTime)
	require.Equal(t, 0.625, snapshot.WinProbability)
	require.Len(t, emitter.stateUpdates, updatesBefore+1)
}

func TestStaleAsyncProbabilityPreservesNewerMatch(t *testing.T) {
	sm, repository, _, channelID := newFixture()

	require.NoError(t, sm.Process(context.Background(), channelID, inGamePayloadAt(107, 100)))
	repository.setBeforeApply(func() {
		require.NoError(t, sm.Process(context.Background(), channelID, inGamePayloadAt(108, 200)))
	})

	require.NoError(t, sm.UpdateWinProbability(context.Background(), channelID, 107, 0.625))

	snapshot := repository.snapshot(t, channelID)
	require.Equal(t, int64(108), snapshot.MatchID)
	require.Zero(t, snapshot.WinProbability)
	require.Equal(t, int64(200), snapshot.LastProviderTimestamp)
}

func TestUpdateStatsPreservesNewerMatchAndSourceFields(t *testing.T) {
	sm, repository, _, channelID := newFixture()

	require.NoError(t, sm.Process(context.Background(), channelID, inGamePayloadAt(120, 100)))
	repository.setBeforeApply(func() {
		require.NoError(t, sm.Process(context.Background(), channelID, inGamePayloadAt(121, 200)))
	})

	settings := model.ChannelDotaSettings{
		ChannelID:     channelID,
		Mmr:           3_025,
		SessionWins:   4,
		SessionLosses: 2,
	}
	require.NoError(t, sm.UpdateStats(context.Background(), channelID, settings))

	snapshot := repository.snapshot(t, channelID)
	require.Equal(t, int64(121), snapshot.MatchID)
	require.Equal(t, StateInGame, snapshot.State)
	require.Equal(t, int64(200), snapshot.LastProviderTimestamp)
	require.Equal(t, 100, snapshot.LastGameTime)
	require.Equal(t, 3_025, snapshot.Mmr)
	require.Equal(t, 4, snapshot.SessionWins)
	require.Equal(t, 2, snapshot.SessionLosses)
}

func TestUpdateStatsRetriesStateUpdateWithoutRewritingCurrentStats(t *testing.T) {
	sm, repository, emitter, channelID := newFixture()
	stateUpdateErr := errors.New("state update unavailable")
	emitter.stateUpdateErr = stateUpdateErr
	settings := model.ChannelDotaSettings{
		ChannelID:     channelID,
		Mmr:           3_025,
		SessionWins:   4,
		SessionLosses: 2,
	}

	err := sm.UpdateStats(context.Background(), channelID, settings)

	require.ErrorIs(t, err, stateUpdateErr)
	firstSnapshot := repository.snapshot(t, channelID)
	require.Equal(t, 3_025, firstSnapshot.Mmr)
	require.Equal(t, 4, firstSnapshot.SessionWins)
	require.Equal(t, 2, firstSnapshot.SessionLosses)
	require.Len(t, emitter.stateUpdateAttempts, 1)
	require.Empty(t, emitter.stateUpdates)

	require.NoError(t, sm.UpdateStats(context.Background(), channelID, settings))
	secondSnapshot := repository.snapshot(t, channelID)
	require.Equal(t, firstSnapshot.Revision, secondSnapshot.Revision)
	require.Len(t, emitter.stateUpdateAttempts, 2)
	require.Len(t, emitter.stateUpdates, 1)
}

func TestUpdateWinProbabilityKeepsThresholdAndInputValidation(t *testing.T) {
	sm, repository, emitter, channelID := newFixture()

	require.Error(t, sm.UpdateWinProbability(context.Background(), channelID, 109, math.NaN()))
	require.NoError(t, sm.Process(context.Background(), channelID, inGamePayloadAt(109, 100)))
	require.NoError(t, sm.UpdateWinProbability(context.Background(), channelID, 109, 0.30))
	updatesBefore := len(emitter.stateUpdates)
	revisionBefore := repository.snapshot(t, channelID).Revision

	require.NoError(t, sm.UpdateWinProbability(context.Background(), channelID, 109, 0.349))
	require.Equal(t, revisionBefore, repository.snapshot(t, channelID).Revision)
	require.Len(t, emitter.stateUpdates, updatesBefore)

	require.NoError(t, sm.UpdateWinProbability(context.Background(), channelID, 109, 0.35))
	require.Equal(t, 0.35, repository.snapshot(t, channelID).WinProbability)
}

func TestRoshanEventIsAcknowledgedAfterPostCommitEmission(t *testing.T) {
	first, repository, emitter, channelID := newFixture()
	payload := inGamePayloadAt(110, 100)
	payload.Events = []gsi.Event{{EventType: eventRoshanKilled, KillerTeam: "dire", GameTime: 500}}

	require.NoError(t, first.Process(context.Background(), channelID, payload))
	second := New(repository, emitter, slog.Default())
	payload.Provider.Timestamp = 101
	require.NoError(t, second.Process(context.Background(), channelID, payload))

	require.Len(t, emitter.roshanKilled, 1)
	require.Contains(t, repository.snapshot(t, channelID).SeenEvents, "roshan_killed:500")
}

func TestRoshanPublishFailureRetriesOnReplayedPayload(t *testing.T) {
	sm, repository, emitter, channelID := newFixture()
	payload := inGamePayloadAt(111, 100)
	payload.Events = []gsi.Event{{EventType: eventRoshanKilled, KillerTeam: "dire", GameTime: 500}}
	emitter.roshanErr = errors.New("temporary publish failure")

	require.NoError(t, sm.Process(context.Background(), channelID, payload))
	require.Empty(t, emitter.roshanKilled)
	require.NotContains(t, repository.snapshot(t, channelID).SeenEvents, "roshan_killed:500")

	emitter.roshanErr = nil
	require.NoError(t, sm.Process(context.Background(), channelID, payload))
	require.Len(t, emitter.roshanKilled, 1)
	require.Contains(t, repository.snapshot(t, channelID).SeenEvents, "roshan_killed:500")

	payload.Provider.Timestamp = 101
	require.NoError(t, sm.Process(context.Background(), channelID, payload))
	require.Len(t, emitter.roshanKilled, 1)
}

func TestRoshanDeliveryAcknowledgementDoesNotMutateNewerMatch(t *testing.T) {
	sm, repository, emitter, channelID := newFixture()
	emitter.onRoshanKilled = func() {
		require.NoError(t, sm.Process(context.Background(), channelID, inGamePayloadAt(113, 200)))
	}
	payload := inGamePayloadAt(112, 100)
	payload.Events = []gsi.Event{{EventType: eventRoshanKilled, KillerTeam: "dire", GameTime: 500}}

	require.NoError(t, sm.Process(context.Background(), channelID, payload))

	snapshot := repository.snapshot(t, channelID)
	require.Equal(t, int64(113), snapshot.MatchID)
	require.NotContains(t, snapshot.SeenEvents, "roshan_killed:500")
	require.Equal(t, []model.OutboxAction{
		model.OutboxActionCreate,
		model.OutboxActionCancel,
		model.OutboxActionCreate,
	}, repository.actionKinds())
}

func TestConcurrentRoshanDeliveryClaimsOnePublisher(t *testing.T) {
	first, repository, emitter, channelID := newFixture()
	second := New(repository, emitter, slog.Default())
	payload := inGamePayloadAt(117, 100)
	payload.Events = []gsi.Event{{EventType: eventRoshanKilled, KillerTeam: "dire", GameTime: 500}}
	started := make(chan struct{}, 1)
	release := make(chan struct{})
	emitter.blockNextRoshan(started, release)

	firstResult := make(chan error, 1)
	go func() {
		firstResult <- first.Process(context.Background(), channelID, payload)
	}()
	<-started

	var releaseOnce sync.Once
	releasePublisher := func() {
		releaseOnce.Do(func() { close(release) })
	}
	defer releasePublisher()

	require.NoError(t, second.Process(context.Background(), channelID, payload))
	require.Len(t, emitter.roshanKilled, 1)

	releasePublisher()
	require.NoError(t, <-firstResult)
	require.Len(t, emitter.roshanKilled, 1)
	require.Contains(t, repository.snapshot(t, channelID).SeenEvents, "roshan_killed:500")
}

func TestRoshanEventClaimCanBeReclaimedAfterLeaseExpiry(t *testing.T) {
	sm, repository, emitter, channelID := newFixture()
	require.NoError(t, sm.Process(context.Background(), channelID, inGamePayloadAt(118, 100)))
	payload := inGamePayloadAt(118, 100)
	payload.Events = []gsi.Event{{EventType: eventRoshanKilled, KillerTeam: "dire", GameTime: 500}}

	repository.setEventClaim(
		t,
		channelID,
		"roshan_killed:500",
		"another-replica",
		time.Now().Add(time.Minute).UnixNano(),
	)
	require.NoError(t, sm.Process(context.Background(), channelID, payload))
	require.Empty(t, emitter.roshanKilled)

	repository.setEventClaim(
		t,
		channelID,
		"roshan_killed:500",
		"crashed-replica",
		time.Now().Add(-time.Second).UnixNano(),
	)
	require.NoError(t, sm.Process(context.Background(), channelID, payload))
	require.Len(t, emitter.roshanKilled, 1)
	require.Contains(t, repository.snapshot(t, channelID).SeenEvents, "roshan_killed:500")
}

func TestAegisPublishFailureRetriesOnReplayedOrNewerPayload(t *testing.T) {
	for _, test := range []struct {
		name              string
		recoveryTimestamp int64
		laterTimestamp    int64
	}{
		{
			name:              "same replay",
			recoveryTimestamp: 100,
			laterTimestamp:    101,
		},
		{
			name:              "newer payload",
			recoveryTimestamp: 101,
			laterTimestamp:    102,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			sm, repository, emitter, channelID := newFixture()
			playerID := 3
			payload := inGamePayloadAt(114, 100)
			payload.Events = []gsi.Event{{
				EventType: eventAegisPicked,
				PlayerID:  &playerID,
				GameTime:  600,
			}}
			emitter.aegisErr = errors.New("temporary publish failure")

			require.NoError(t, sm.Process(context.Background(), channelID, payload))
			require.Empty(t, emitter.aegisPickup)
			require.NotContains(t, repository.snapshot(t, channelID).SeenEvents, "aegis_picked_up:600")

			emitter.aegisErr = nil
			payload.Provider.Timestamp = test.recoveryTimestamp
			require.NoError(t, sm.Process(context.Background(), channelID, payload))
			require.Len(t, emitter.aegisPickup, 1)
			require.NotNil(t, emitter.aegisPickup[0].PlayerID)
			require.Equal(t, playerID, *emitter.aegisPickup[0].PlayerID)
			require.Contains(t, repository.snapshot(t, channelID).SeenEvents, "aegis_picked_up:600")

			payload.Provider.Timestamp = test.laterTimestamp
			require.NoError(t, sm.Process(context.Background(), channelID, payload))
			require.Len(t, emitter.aegisPickup, 1)
		})
	}
}

func TestAegisDeliveryAcknowledgementDoesNotMutateNewerMatch(t *testing.T) {
	sm, repository, emitter, channelID := newFixture()
	emitter.onAegisPickup = func() {
		require.NoError(t, sm.Process(context.Background(), channelID, inGamePayloadAt(116, 200)))
	}
	playerID := 3
	payload := inGamePayloadAt(115, 100)
	payload.Events = []gsi.Event{{
		EventType: eventAegisPicked,
		PlayerID:  &playerID,
		GameTime:  600,
	}}

	require.NoError(t, sm.Process(context.Background(), channelID, payload))

	snapshot := repository.snapshot(t, channelID)
	require.Equal(t, int64(116), snapshot.MatchID)
	require.NotContains(t, snapshot.SeenEvents, "aegis_picked_up:600")
	require.Equal(t, []model.OutboxAction{
		model.OutboxActionCreate,
		model.OutboxActionCancel,
		model.OutboxActionCreate,
	}, repository.actionKinds())
}
