package processor

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/twirapp/twir/apps/dota/internal/gsi"
	"github.com/twirapp/twir/apps/dota/internal/match"
	busapi "github.com/twirapp/twir/libs/bus-core/api"
	busdota "github.com/twirapp/twir/libs/bus-core/dota"
	dotarepository "github.com/twirapp/twir/libs/repositories/dota"
	"github.com/twirapp/twir/libs/repositories/dota/model"
	"go.uber.org/fx"
)

type fakeWinProbabilityProvider struct {
	probability float64
	err         error
	matchIDs    []int64
	callCh      chan int64
	mu          sync.Mutex
}

type blockingWinProbabilityProvider struct {
	probability float64
	started     chan struct{}
	release     chan struct{}
	releaseOnce sync.Once
	mu          sync.Mutex
	calls       int
}

type cancellationBlockingWinProbabilityProvider struct {
	probability              float64
	returnSuccessAfterCancel bool
	started                  chan struct{}
	canceled                 chan struct{}
	release                  chan struct{}
	releaseOnce              sync.Once
}

type stagedWinProbabilityProvider struct {
	probabilities map[int64]float64
	started       chan int64
	release       map[int64]chan struct{}
}

type fakeLifecycle struct {
	hooks []fx.Hook
}

func (f *fakeLifecycle) Append(hook fx.Hook) {
	f.hooks = append(f.hooks, hook)
}

func (f *blockingWinProbabilityProvider) WinProbability(
	ctx context.Context,
	_ int64,
) (float64, error) {
	f.mu.Lock()
	f.calls++
	f.mu.Unlock()

	select {
	case f.started <- struct{}{}:
	case <-ctx.Done():
		return 0, ctx.Err()
	}

	select {
	case <-f.release:
		return f.probability, nil
	case <-ctx.Done():
		return 0, ctx.Err()
	}
}

func (f *blockingWinProbabilityProvider) Release() {
	f.releaseOnce.Do(func() {
		close(f.release)
	})
}

func (f *blockingWinProbabilityProvider) Calls() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.calls
}

func (f *cancellationBlockingWinProbabilityProvider) WinProbability(
	ctx context.Context,
	_ int64,
) (float64, error) {
	select {
	case f.started <- struct{}{}:
	case <-ctx.Done():
		return 0, ctx.Err()
	}

	<-ctx.Done()
	f.canceled <- struct{}{}
	<-f.release

	if f.returnSuccessAfterCancel {
		return f.probability, nil
	}

	return 0, ctx.Err()
}

func (f *cancellationBlockingWinProbabilityProvider) Release() {
	f.releaseOnce.Do(func() {
		close(f.release)
	})
}

func (f *stagedWinProbabilityProvider) WinProbability(
	ctx context.Context,
	matchID int64,
) (float64, error) {
	select {
	case f.started <- matchID:
	case <-ctx.Done():
		return 0, ctx.Err()
	}

	select {
	case <-f.release[matchID]:
		return f.probabilities[matchID], nil
	case <-ctx.Done():
		return 0, ctx.Err()
	}
}

func (f *fakeWinProbabilityProvider) WinProbability(_ context.Context, matchID int64) (float64, error) {
	f.mu.Lock()
	f.matchIDs = append(f.matchIDs, matchID)
	f.mu.Unlock()

	if f.callCh != nil {
		f.callCh <- matchID
	}

	return f.probability, f.err
}

func (f *fakeWinProbabilityProvider) MatchIDs() []int64 {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]int64(nil), f.matchIDs...)
}

type fakeRepo struct {
	mu       sync.Mutex
	settings model.ChannelDotaSettings
	states   map[uuid.UUID]model.MatchState
}

func (f *fakeRepo) GetByChannelID(
	_ context.Context,
	_ uuid.UUID,
) (model.ChannelDotaSettings, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
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
	_ uuid.UUID,
	_ bool,
	_ int,
) (model.ChannelDotaSettings, error) {
	return f.settings, nil
}

func (f *fakeRepo) ApplyMatchResultOnce(
	_ context.Context,
	_ dotarepository.ApplyMatchResultInput,
) (model.ChannelDotaSettings, error) {
	return f.settings, nil
}

func (f *fakeRepo) GetMatchState(
	_ context.Context,
	channelID uuid.UUID,
) (model.MatchState, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	state, ok := f.states[channelID]
	if !ok {
		return model.MatchState{ChannelID: channelID, Snapshot: json.RawMessage(`{}`)}, nil
	}
	state.Snapshot = append(json.RawMessage(nil), state.Snapshot...)
	return state, nil
}

func (f *fakeRepo) ApplyMatchStateTransition(
	_ context.Context,
	input dotarepository.ApplyMatchStateTransitionInput,
) (bool, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	current, ok := f.states[input.ChannelID]
	if !ok {
		current = model.MatchState{ChannelID: input.ChannelID, Snapshot: json.RawMessage(`{}`)}
	}
	if current.Revision != input.ExpectedRevision {
		return false, nil
	}

	f.states[input.ChannelID] = model.MatchState{
		ChannelID:         input.ChannelID,
		Revision:          current.Revision + 1,
		ProviderTimestamp: input.ProviderTimestamp,
		Snapshot:          append(json.RawMessage(nil), input.Snapshot...),
	}
	return true, nil
}

func (f *fakeRepo) ClaimPredictionActions(
	_ context.Context,
	_ dotarepository.ClaimPredictionActionsInput,
) ([]model.ClaimedOutboxAction, error) {
	return nil, errors.New("not implemented")
}

func (f *fakeRepo) RenewPredictionAction(
	_ context.Context,
	_ uuid.UUID,
	_ uuid.UUID,
	_ time.Duration,
) error {
	return errors.New("not implemented")
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
	roshanErr     error
	stateUpdates  []busapi.DotaStateUpdateMessage
	stateUpdateCh chan busapi.DotaStateUpdateMessage
}

func (f *fakeEmitter) MatchStarted(_ context.Context, _ busdota.MatchStartedMessage) error {
	return nil
}

func (f *fakeEmitter) MatchEnded(_ context.Context, _ busdota.MatchEndedMessage) error {
	return nil
}

func (f *fakeEmitter) MatchAbandoned(_ context.Context, _ busdota.MatchAbandonedMessage) error {
	return nil
}

func (f *fakeEmitter) RoshanKilled(_ context.Context, _ busdota.RoshanKilledMessage) error {
	return f.roshanErr
}

func (f *fakeEmitter) AegisPickup(_ context.Context, _ busdota.AegisPickupMessage) error {
	return nil
}

func (f *fakeEmitter) StateUpdate(_ context.Context, msg busapi.DotaStateUpdateMessage) error {
	f.stateUpdates = append(f.stateUpdates, msg)
	if f.stateUpdateCh != nil {
		f.stateUpdateCh <- msg
	}
	return nil
}

func newStateMachine() (*match.StateMachine, *fakeEmitter, uuid.UUID) {
	channelID := uuid.New()
	repo := &fakeRepo{
		settings: model.ChannelDotaSettings{
			ChannelID: channelID,
			Mmr:       3000,
			MmrDelta:  25,
		},
		states: make(map[uuid.UUID]model.MatchState),
	}
	emitter := &fakeEmitter{}

	return match.New(repo, emitter, slog.Default()), emitter, channelID
}

func inGamePayload(matchID int64) gsi.Payload {
	return gsi.Payload{
		Provider: gsi.Provider{Timestamp: matchID},
		Map: &gsi.Map{
			MatchID:   matchID,
			GameState: gsi.GameStateInProgress,
		},
		Player: &gsi.Player{
			Activity:  gsi.PlayerActivityPlaying,
			TeamName:  "radiant",
			AccountID: 12345,
		},
	}
}

func waitFor[T any](t *testing.T, ch <-chan T, description string) T {
	t.Helper()

	timer := time.NewTimer(time.Second)
	defer timer.Stop()

	select {
	case value := <-ch:
		return value
	case <-timer.C:
		t.Fatalf("timed out waiting for %s", description)
		var zero T
		return zero
	}
}

func requireNoReceive[T any](t *testing.T, ch <-chan T, description string) {
	t.Helper()

	timer := time.NewTimer(time.Second)
	defer timer.Stop()

	select {
	case value := <-ch:
		t.Fatalf("unexpected %s: %v", description, value)
	case <-timer.C:
	}
}

func TestProcessReturnsBeforeBlockedWinProbabilityAndCommitsResult(t *testing.T) {
	sm, emitter, channelID := newStateMachine()
	emitter.stateUpdateCh = make(chan busapi.DotaStateUpdateMessage, 2)
	provider := &blockingWinProbabilityProvider{
		probability: 0.625,
		started:     make(chan struct{}, 1),
		release:     make(chan struct{}),
	}
	t.Cleanup(provider.Release)
	processor := New(sm, provider, slog.Default(), &fakeLifecycle{})

	processDone := make(chan error, 1)
	go func() {
		processDone <- processor.Process(context.Background(), channelID, inGamePayload(2004))
	}()

	waitFor(t, provider.started, "win probability request to start")
	initial := waitFor(t, emitter.stateUpdateCh, "initial state update")
	require.Zero(t, initial.WinProbability)
	require.NoError(t, waitFor(t, processDone, "Process to return"))

	snapshot, err := sm.GetSnapshot(context.Background(), channelID)
	require.NoError(t, err)
	require.Zero(t, snapshot.WinProbability)

	provider.Release()
	update := waitFor(t, emitter.stateUpdateCh, "win probability state update")
	require.Equal(t, 0.625, update.WinProbability)

	snapshot, err = sm.GetSnapshot(context.Background(), channelID)
	require.NoError(t, err)
	require.Equal(t, int64(2004), snapshot.MatchID)
	require.Equal(t, 0.625, snapshot.WinProbability)
}

func TestProcessCoalescesBlockedWinProbabilityFetches(t *testing.T) {
	sm, emitter, channelID := newStateMachine()
	emitter.stateUpdateCh = make(chan busapi.DotaStateUpdateMessage, 2)
	provider := &blockingWinProbabilityProvider{
		probability: 0.625,
		started:     make(chan struct{}, 2),
		release:     make(chan struct{}),
	}
	t.Cleanup(provider.Release)
	processor := New(sm, provider, slog.Default(), &fakeLifecycle{})
	payload := inGamePayload(2005)

	require.NoError(t, processor.Process(context.Background(), channelID, payload))
	waitFor(t, provider.started, "first win probability request to start")
	waitFor(t, emitter.stateUpdateCh, "initial state update")
	require.NoError(t, processor.Process(context.Background(), channelID, payload))

	timer := time.NewTimer(time.Second)
	defer timer.Stop()
	select {
	case <-provider.started:
		t.Fatal("started a duplicate win probability request")
	case <-timer.C:
	}
	require.Equal(t, 1, provider.Calls())

	provider.Release()
	update := waitFor(t, emitter.stateUpdateCh, "win probability state update")
	require.Equal(t, 0.625, update.WinProbability)
}

func TestProcessRefreshesWinProbabilityAfterProcessingLiveMatch(t *testing.T) {
	sm, emitter, channelID := newStateMachine()
	emitter.stateUpdateCh = make(chan busapi.DotaStateUpdateMessage, 2)
	provider := &fakeWinProbabilityProvider{
		probability: 0.625,
		callCh:      make(chan int64, 1),
	}
	processor := New(sm, provider, slog.Default(), &fakeLifecycle{})
	ctx := context.Background()

	require.NoError(t, processor.Process(ctx, channelID, inGamePayload(2001)))
	initial := waitFor(t, emitter.stateUpdateCh, "initial state update")
	require.Zero(t, initial.WinProbability)
	require.Equal(t, int64(2001), waitFor(t, provider.callCh, "win probability request"))
	update := waitFor(t, emitter.stateUpdateCh, "win probability state update")
	require.Equal(t, 0.625, update.WinProbability)
	require.Equal(t, []int64{2001}, provider.MatchIDs())

	snapshot, err := sm.GetSnapshot(ctx, channelID)
	require.NoError(t, err)
	require.True(t, snapshot.InGame)
	require.Equal(t, int64(2001), snapshot.MatchID)
	require.Equal(t, 0.625, snapshot.WinProbability)
	require.Len(t, emitter.stateUpdates, 2)
	require.Equal(t, 0.625, emitter.stateUpdates[1].WinProbability)
}

func TestProcessKeepsGsiAvailableWhenWinProbabilityFails(t *testing.T) {
	sm, emitter, channelID := newStateMachine()
	provider := &fakeWinProbabilityProvider{
		err:    errors.New("stratz unavailable"),
		callCh: make(chan int64, 1),
	}
	processor := New(sm, provider, slog.Default(), &fakeLifecycle{})

	require.NoError(t, processor.Process(context.Background(), channelID, inGamePayload(2002)))
	require.Equal(t, int64(2002), waitFor(t, provider.callCh, "win probability request"))
	require.Equal(t, []int64{2002}, provider.MatchIDs())

	snapshot, err := sm.GetSnapshot(context.Background(), channelID)
	require.NoError(t, err)
	require.True(t, snapshot.InGame)
	require.Equal(t, int64(2002), snapshot.MatchID)
	require.Zero(t, snapshot.WinProbability)
	require.Len(t, emitter.stateUpdates, 1)
}

func TestProcessKeepsStateCommittedWhenRoshanPublicationFails(t *testing.T) {
	sm, emitter, channelID := newStateMachine()
	emitter.roshanErr = errors.New("publish failed")
	provider := &fakeWinProbabilityProvider{
		probability: 0.625,
		callCh:      make(chan int64, 1),
	}
	processor := New(sm, provider, slog.Default(), &fakeLifecycle{})
	payload := inGamePayload(2003)
	payload.Events = []gsi.Event{{EventType: "roshan_killed", GameTime: 300}}

	require.NoError(t, processor.Process(context.Background(), channelID, payload))
	require.Equal(t, int64(2003), waitFor(t, provider.callCh, "win probability request"))
	require.Equal(t, []int64{2003}, provider.MatchIDs())

	snapshot, err := sm.GetSnapshot(context.Background(), channelID)
	require.NoError(t, err)
	require.Equal(t, int64(2003), snapshot.MatchID)
}

func TestOnStopCancelsAndWaitsForWinProbabilityJob(t *testing.T) {
	sm, _, channelID := newStateMachine()
	lifecycle := &fakeLifecycle{}
	provider := &cancellationBlockingWinProbabilityProvider{
		started:  make(chan struct{}, 1),
		canceled: make(chan struct{}, 1),
		release:  make(chan struct{}),
	}
	t.Cleanup(provider.Release)
	processor := New(sm, provider, slog.Default(), lifecycle)

	require.NoError(t, processor.Process(context.Background(), channelID, inGamePayload(2006)))
	waitFor(t, provider.started, "win probability request to start")

	stopResult := make(chan error, 1)
	go func() {
		stopResult <- lifecycle.hooks[0].OnStop(context.Background())
	}()

	waitFor(t, provider.canceled, "win probability request cancellation")
	requireNoReceive(t, stopResult, "stop result before win probability request exits")
	provider.Release()
	require.NoError(t, waitFor(t, stopResult, "stop to finish"))
}

func TestOnStopSkipsSuccessfulWinProbabilityResultAfterCancellation(t *testing.T) {
	sm, emitter, channelID := newStateMachine()
	emitter.stateUpdateCh = make(chan busapi.DotaStateUpdateMessage, 2)
	lifecycle := &fakeLifecycle{}
	provider := &cancellationBlockingWinProbabilityProvider{
		probability:              0.625,
		returnSuccessAfterCancel: true,
		started:                  make(chan struct{}, 1),
		canceled:                 make(chan struct{}, 1),
		release:                  make(chan struct{}),
	}
	t.Cleanup(provider.Release)
	processor := New(sm, provider, slog.Default(), lifecycle)

	require.NoError(t, processor.Process(context.Background(), channelID, inGamePayload(2007)))
	waitFor(t, emitter.stateUpdateCh, "initial state update")
	waitFor(t, provider.started, "win probability request to start")

	stopResult := make(chan error, 1)
	go func() {
		stopResult <- lifecycle.hooks[0].OnStop(context.Background())
	}()

	waitFor(t, provider.canceled, "win probability request cancellation")
	provider.Release()
	require.NoError(t, waitFor(t, stopResult, "stop to finish"))
	requireNoReceive(t, emitter.stateUpdateCh, "win probability state update after stop")

	snapshot, err := sm.GetSnapshot(context.Background(), channelID)
	require.NoError(t, err)
	require.Zero(t, snapshot.WinProbability)
}

func TestProcessDoesNotScheduleWinProbabilityAfterStop(t *testing.T) {
	sm, emitter, channelID := newStateMachine()
	emitter.stateUpdateCh = make(chan busapi.DotaStateUpdateMessage, 3)
	lifecycle := &fakeLifecycle{}
	provider := &fakeWinProbabilityProvider{
		probability: 0.625,
		callCh:      make(chan int64, 2),
	}
	processor := New(sm, provider, slog.Default(), lifecycle)

	require.NoError(t, processor.Process(context.Background(), channelID, inGamePayload(2008)))
	waitFor(t, emitter.stateUpdateCh, "initial state update")
	require.Equal(t, int64(2008), waitFor(t, provider.callCh, "initial win probability request"))
	waitFor(t, emitter.stateUpdateCh, "initial win probability state update")
	callsBefore := len(provider.MatchIDs())

	require.NoError(t, lifecycle.hooks[0].OnStop(context.Background()))
	require.NoError(t, processor.Process(context.Background(), channelID, inGamePayload(2009)))
	requireNoReceive(t, provider.callCh, "post-stop win probability request")
	require.Equal(t, callsBefore, len(provider.MatchIDs()))

	snapshot, err := sm.GetSnapshot(context.Background(), channelID)
	require.NoError(t, err)
	require.Equal(t, int64(2009), snapshot.MatchID)
}

func TestProcessDiscardsDelayedWinProbabilityForPreviousMatch(t *testing.T) {
	sm, emitter, channelID := newStateMachine()
	emitter.stateUpdateCh = make(chan busapi.DotaStateUpdateMessage, 4)
	releaseA := make(chan struct{})
	releaseB := make(chan struct{})
	provider := &stagedWinProbabilityProvider{
		probabilities: map[int64]float64{2010: 0.625, 2011: 0.700},
		started:       make(chan int64, 2),
		release:       map[int64]chan struct{}{2010: releaseA, 2011: releaseB},
	}
	processor := New(sm, provider, slog.Default(), &fakeLifecycle{})

	require.NoError(t, processor.Process(context.Background(), channelID, inGamePayload(2010)))
	waitFor(t, emitter.stateUpdateCh, "match A initial state update")
	require.Equal(t, int64(2010), waitFor(t, provider.started, "match A win probability request"))

	require.NoError(t, processor.Process(context.Background(), channelID, inGamePayload(2011)))
	require.Equal(t, int64(2011), waitFor(t, provider.started, "match B win probability request"))

	close(releaseB)
	update := waitFor(t, emitter.stateUpdateCh, "match B win probability state update")
	require.Equal(t, 0.700, update.WinProbability)
	close(releaseA)

	jobsDone := make(chan struct{})
	go func() {
		processor.jobs.Wait()
		close(jobsDone)
	}()
	waitFor(t, jobsDone, "win probability jobs to finish")

	snapshot, err := sm.GetSnapshot(context.Background(), channelID)
	require.NoError(t, err)
	require.Equal(t, int64(2011), snapshot.MatchID)
	require.Equal(t, 0.700, snapshot.WinProbability)
}
