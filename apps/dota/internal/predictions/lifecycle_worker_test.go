package predictions

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/twirapp/twir/apps/dota/internal/match"
	busapi "github.com/twirapp/twir/libs/bus-core/api"
	busdota "github.com/twirapp/twir/libs/bus-core/dota"
	dotarepository "github.com/twirapp/twir/libs/repositories/dota"
	"github.com/twirapp/twir/libs/repositories/dota/model"
	"go.uber.org/fx"
)

type workerCallRecorder struct {
	mu    sync.Mutex
	calls []string
}

func (r *workerCallRecorder) add(call string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.calls = append(r.calls, call)
}

func (r *workerCallRecorder) all() []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]string(nil), r.calls...)
}

type workerCompletion struct {
	actionID  uuid.UUID
	lockToken uuid.UUID
}

type workerRetry struct {
	actionID    uuid.UUID
	lockToken   uuid.UUID
	availableAt time.Time
}

type workerRenewal struct {
	actionID  uuid.UUID
	lockToken uuid.UUID
	lease     time.Duration
}

type workerDeliveryClaim struct {
	key   string
	token string
	ttl   time.Duration
}

type workerDeliveryRelease struct {
	key   string
	token string
}

type workerDeliveryRenewal struct {
	key   string
	token string
	ttl   time.Duration
}

type workerDeliveryCompletion struct {
	key   string
	token string
	ttl   time.Duration
}

type workerDeliveryRecord struct {
	token     string
	delivered bool
}

type fakeWorkerMatchEndedDeliveryStore struct {
	mu sync.Mutex

	claims     map[string]workerDeliveryRecord
	claimErr   error
	renewErr   error
	releaseErr error
	renewFn    func(context.Context, string, string, time.Duration) (bool, error)

	claimCalls      []workerDeliveryClaim
	renewCalls      []workerDeliveryRenewal
	completionCalls []workerDeliveryCompletion
	releaseCalls    []workerDeliveryRelease
}

func newFakeWorkerMatchEndedDeliveryStore() *fakeWorkerMatchEndedDeliveryStore {
	return &fakeWorkerMatchEndedDeliveryStore{claims: make(map[string]workerDeliveryRecord)}
}

func (s *fakeWorkerMatchEndedDeliveryStore) ClaimMatchEndedDelivery(
	_ context.Context,
	key string,
	token string,
	ttl time.Duration,
) (matchEndedDeliveryState, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.claimCalls = append(s.claimCalls, workerDeliveryClaim{key: key, token: token, ttl: ttl})
	if s.claimErr != nil {
		return matchEndedDeliveryPending, s.claimErr
	}
	if claim, exists := s.claims[key]; exists {
		if claim.delivered {
			return matchEndedDeliveryDelivered, nil
		}
		return matchEndedDeliveryPending, nil
	}
	s.claims[key] = workerDeliveryRecord{token: token}
	return matchEndedDeliveryAcquired, nil
}

func (s *fakeWorkerMatchEndedDeliveryStore) CompleteMatchEndedDelivery(
	_ context.Context,
	key string,
	token string,
	ttl time.Duration,
) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.completionCalls = append(s.completionCalls, workerDeliveryCompletion{key: key, token: token, ttl: ttl})
	claim, exists := s.claims[key]
	if !exists || claim.delivered || claim.token != token {
		return false, nil
	}
	claim.delivered = true
	s.claims[key] = claim
	return true, nil
}

func (s *fakeWorkerMatchEndedDeliveryStore) RenewMatchEndedDelivery(
	ctx context.Context,
	key string,
	token string,
	ttl time.Duration,
) (bool, error) {
	s.mu.Lock()
	s.renewCalls = append(s.renewCalls, workerDeliveryRenewal{key: key, token: token, ttl: ttl})
	renewFn := s.renewFn
	renewErr := s.renewErr
	claim, exists := s.claims[key]
	s.mu.Unlock()
	if renewFn != nil {
		return renewFn(ctx, key, token, ttl)
	}
	if renewErr != nil {
		return false, renewErr
	}
	if !exists || claim.delivered || claim.token != token {
		return false, nil
	}

	return true, nil
}

func (s *fakeWorkerMatchEndedDeliveryStore) ReleaseMatchEndedDelivery(
	_ context.Context,
	key string,
	token string,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.releaseCalls = append(s.releaseCalls, workerDeliveryRelease{key: key, token: token})
	if s.releaseErr != nil {
		return s.releaseErr
	}
	if claim, exists := s.claims[key]; exists && !claim.delivered && claim.token == token {
		delete(s.claims, key)
	}
	return nil
}

func (s *fakeWorkerMatchEndedDeliveryStore) claimHistory() []workerDeliveryClaim {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]workerDeliveryClaim(nil), s.claimCalls...)
}

func (s *fakeWorkerMatchEndedDeliveryStore) releaseHistory() []workerDeliveryRelease {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]workerDeliveryRelease(nil), s.releaseCalls...)
}

func (s *fakeWorkerMatchEndedDeliveryStore) renewalHistory() []workerDeliveryRenewal {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]workerDeliveryRenewal(nil), s.renewCalls...)
}

func (s *fakeWorkerMatchEndedDeliveryStore) completionHistory() []workerDeliveryCompletion {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]workerDeliveryCompletion(nil), s.completionCalls...)
}

func (s *fakeWorkerMatchEndedDeliveryStore) hasClaim(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, exists := s.claims[key]
	return exists
}

type fakeWorkerRepository struct {
	mu sync.Mutex

	claimed     []model.ClaimedOutboxAction
	states      map[uuid.UUID]model.MatchState
	claimErr    error
	completeErr error
	retryErr    error
	renewErr    error
	claimFn     func(context.Context, dotarepository.ClaimPredictionActionsInput) ([]model.ClaimedOutboxAction, error)
	renewFn     func(context.Context, uuid.UUID, uuid.UUID, time.Duration) error

	claimInputs []dotarepository.ClaimPredictionActionsInput
	completions []workerCompletion
	retries     []workerRetry
	renewals    []workerRenewal
	calls       *workerCallRecorder
}

func newFakeWorkerRepository() *fakeWorkerRepository {
	return &fakeWorkerRepository{states: make(map[uuid.UUID]model.MatchState)}
}

func (r *fakeWorkerRepository) GetMatchState(
	_ context.Context,
	channelID uuid.UUID,
) (model.MatchState, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	state, ok := r.states[channelID]
	if !ok {
		return model.MatchState{ChannelID: channelID, Snapshot: json.RawMessage(`{}`)}, nil
	}
	state.Snapshot = append(json.RawMessage(nil), state.Snapshot...)
	return state, nil
}

func (r *fakeWorkerRepository) ClaimPredictionActions(
	ctx context.Context,
	input dotarepository.ClaimPredictionActionsInput,
) ([]model.ClaimedOutboxAction, error) {
	r.mu.Lock()
	r.claimInputs = append(r.claimInputs, input)
	claimFn := r.claimFn
	claimed := append([]model.ClaimedOutboxAction(nil), r.claimed...)
	err := r.claimErr
	r.mu.Unlock()

	if claimFn != nil {
		return claimFn(ctx, input)
	}
	return claimed, err
}

func (r *fakeWorkerRepository) CompletePredictionAction(
	_ context.Context,
	actionID uuid.UUID,
	lockToken uuid.UUID,
) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.completions = append(r.completions, workerCompletion{actionID: actionID, lockToken: lockToken})
	if r.calls != nil {
		r.calls.add("complete")
	}
	return r.completeErr
}

func (r *fakeWorkerRepository) RetryPredictionAction(
	_ context.Context,
	actionID uuid.UUID,
	lockToken uuid.UUID,
	availableAt time.Time,
) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.retries = append(r.retries, workerRetry{
		actionID:    actionID,
		lockToken:   lockToken,
		availableAt: availableAt,
	})
	if r.calls != nil {
		r.calls.add("retry")
	}
	return r.retryErr
}

func (r *fakeWorkerRepository) RenewPredictionAction(
	ctx context.Context,
	actionID uuid.UUID,
	lockToken uuid.UUID,
	lease time.Duration,
) error {
	r.mu.Lock()
	r.renewals = append(r.renewals, workerRenewal{actionID: actionID, lockToken: lockToken, lease: lease})
	renewFn := r.renewFn
	err := r.renewErr
	calls := r.calls
	r.mu.Unlock()
	if calls != nil {
		calls.add("renew")
	}
	if renewFn != nil {
		return renewFn(ctx, actionID, lockToken, lease)
	}
	return err
}

func (r *fakeWorkerRepository) claimedInputs() []dotarepository.ClaimPredictionActionsInput {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]dotarepository.ClaimPredictionActionsInput(nil), r.claimInputs...)
}

func (r *fakeWorkerRepository) completed() []workerCompletion {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]workerCompletion(nil), r.completions...)
}

func (r *fakeWorkerRepository) retried() []workerRetry {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]workerRetry(nil), r.retries...)
}

func (r *fakeWorkerRepository) renewed() []workerRenewal {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]workerRenewal(nil), r.renewals...)
}

type fakeWorkerPredictions struct {
	mu sync.Mutex

	createErr  error
	resolveErr error
	cancelErr  error
	settings   model.ChannelDotaSettings

	resolveStarted       chan<- struct{}
	resolveCanceled      chan<- struct{}
	resolveRelease       <-chan struct{}
	resolveWaitForCancel bool

	creates  []match.LifecycleAction
	resolves []match.LifecycleAction
	cancels  []match.LifecycleAction
	calls    *workerCallRecorder
}

func (p *fakeWorkerPredictions) Create(_ context.Context, action match.LifecycleAction) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.creates = append(p.creates, action)
	if p.calls != nil {
		p.calls.add("create")
	}
	return p.createErr
}

func (p *fakeWorkerPredictions) Resolve(
	ctx context.Context,
	action match.LifecycleAction,
) (model.ChannelDotaSettings, error) {
	p.mu.Lock()
	p.resolves = append(p.resolves, action)
	if p.calls != nil {
		p.calls.add("resolve")
	}
	settings := p.settings
	err := p.resolveErr
	started := p.resolveStarted
	canceled := p.resolveCanceled
	release := p.resolveRelease
	waitForCancel := p.resolveWaitForCancel
	p.mu.Unlock()
	if started != nil {
		select {
		case started <- struct{}{}:
		case <-ctx.Done():
			return settings, ctx.Err()
		}
	}
	if waitForCancel {
		<-ctx.Done()
		if canceled != nil {
			canceled <- struct{}{}
		}
		if release != nil {
			<-release
		}
		return settings, err
	}
	if release != nil {
		select {
		case <-release:
		case <-ctx.Done():
			return settings, ctx.Err()
		}
	}
	return settings, err
}

func (p *fakeWorkerPredictions) Cancel(_ context.Context, action match.LifecycleAction) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.cancels = append(p.cancels, action)
	if p.calls != nil {
		p.calls.add("cancel")
	}
	return p.cancelErr
}

func (p *fakeWorkerPredictions) createCalls() []match.LifecycleAction {
	p.mu.Lock()
	defer p.mu.Unlock()
	return append([]match.LifecycleAction(nil), p.creates...)
}

func (p *fakeWorkerPredictions) resolveCalls() []match.LifecycleAction {
	p.mu.Lock()
	defer p.mu.Unlock()
	return append([]match.LifecycleAction(nil), p.resolves...)
}

func (p *fakeWorkerPredictions) cancelCalls() []match.LifecycleAction {
	p.mu.Lock()
	defer p.mu.Unlock()
	return append([]match.LifecycleAction(nil), p.cancels...)
}

type workerStatsUpdate struct {
	channelID uuid.UUID
	settings  model.ChannelDotaSettings
}

type fakeWorkerStatsUpdater struct {
	mu           sync.Mutex
	err          error
	updates      []workerStatsUpdate
	calls        *workerCallRecorder
	stateEmitter *fakeWorkerStateUpdateEmitter
}

func (u *fakeWorkerStatsUpdater) UpdateStats(
	ctx context.Context,
	channelID uuid.UUID,
	settings model.ChannelDotaSettings,
) error {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.updates = append(u.updates, workerStatsUpdate{channelID: channelID, settings: settings})
	if u.calls != nil {
		u.calls.add("update_stats")
	}
	if u.stateEmitter != nil {
		return u.stateEmitter.StateUpdate(ctx, busapi.DotaStateUpdateMessage{
			ChannelID:     channelID.String(),
			Mmr:           settings.Mmr,
			SessionWins:   settings.SessionWins,
			SessionLosses: settings.SessionLosses,
		})
	}
	return u.err
}

func (u *fakeWorkerStatsUpdater) all() []workerStatsUpdate {
	u.mu.Lock()
	defer u.mu.Unlock()
	return append([]workerStatsUpdate(nil), u.updates...)
}

type fakeWorkerStateUpdateEmitter struct {
	mu       sync.Mutex
	failNext error
	attempts []busapi.DotaStateUpdateMessage
	messages []busapi.DotaStateUpdateMessage
}

func (e *fakeWorkerStateUpdateEmitter) StateUpdate(
	_ context.Context,
	message busapi.DotaStateUpdateMessage,
) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.attempts = append(e.attempts, message)
	if e.failNext != nil {
		err := e.failNext
		e.failNext = nil
		return err
	}
	e.messages = append(e.messages, message)
	return nil
}

func (e *fakeWorkerStateUpdateEmitter) all() []busapi.DotaStateUpdateMessage {
	e.mu.Lock()
	defer e.mu.Unlock()
	return append([]busapi.DotaStateUpdateMessage(nil), e.messages...)
}

func (e *fakeWorkerStateUpdateEmitter) allAttempts() []busapi.DotaStateUpdateMessage {
	e.mu.Lock()
	defer e.mu.Unlock()
	return append([]busapi.DotaStateUpdateMessage(nil), e.attempts...)
}

type fakeWorkerMatchEndedEmitter struct {
	mu       sync.Mutex
	err      error
	messages []busdota.MatchEndedMessage
	calls    *workerCallRecorder
	started  chan<- struct{}
	release  <-chan struct{}
}

func (e *fakeWorkerMatchEndedEmitter) MatchEnded(ctx context.Context, message busdota.MatchEndedMessage) error {
	e.mu.Lock()
	err := e.err
	started := e.started
	release := e.release
	calls := e.calls
	e.mu.Unlock()

	if calls != nil {
		calls.add("match_ended")
	}
	if started != nil {
		select {
		case started <- struct{}{}:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	if release != nil {
		select {
		case <-release:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	if err != nil {
		return err
	}

	e.mu.Lock()
	defer e.mu.Unlock()
	e.messages = append(e.messages, message)
	return nil
}

func (e *fakeWorkerMatchEndedEmitter) all() []busdota.MatchEndedMessage {
	e.mu.Lock()
	defer e.mu.Unlock()
	return append([]busdota.MatchEndedMessage(nil), e.messages...)
}

type fakeWorkerLifecycle struct {
	hooks []fx.Hook
}

func (l *fakeWorkerLifecycle) Append(hook fx.Hook) {
	l.hooks = append(l.hooks, hook)
}

func newWorkerForTest(
	repository *fakeWorkerRepository,
	predictionActions *fakeWorkerPredictions,
	stats *fakeWorkerStatsUpdater,
	emitter *fakeWorkerMatchEndedEmitter,
	lifecycle fx.Lifecycle,
) *LifecycleWorker {
	return newLifecycleWorker(
		repository,
		predictionActions,
		stats,
		newFakeWorkerMatchEndedDeliveryStore(),
		emitter,
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		lifecycle,
		time.Hour,
	)
}

func lifecycleActionClaim(
	t testing.TB,
	rowAction model.OutboxAction,
	action match.LifecycleAction,
) model.ClaimedOutboxAction {
	t.Helper()
	payload, err := json.Marshal(action)
	require.NoError(t, err)

	return model.ClaimedOutboxAction{
		ID: uuid.New(),
		OutboxActionInput: model.OutboxActionInput{
			ChannelID: action.ChannelID,
			MatchID:   action.MatchID,
			Action:    rowAction,
			Sequence:  1,
			Payload:   payload,
		},
		LockToken: uuid.New(),
		Attempts:  1,
	}
}

func workerState(t testing.TB, snapshot match.Snapshot) model.MatchState {
	t.Helper()
	payload, err := json.Marshal(snapshot)
	require.NoError(t, err)
	return model.MatchState{
		ChannelID:         snapshot.ChannelID,
		Revision:          int64(snapshot.Revision),
		ProviderTimestamp: snapshot.LastProviderTimestamp,
		Snapshot:          payload,
	}
}

func TestLifecycleWorkerProcessOnceCreatesThenCompletesClaim(t *testing.T) {
	channelID := uuid.New()
	action := match.LifecycleAction{
		Kind:           match.ActionCreate,
		ChannelID:      channelID,
		MatchID:        101,
		SteamAccountID: "12345",
		HeroName:       "axe",
		TeamKnown:      true,
	}
	recorder := &workerCallRecorder{}
	repository := newFakeWorkerRepository()
	repository.calls = recorder
	repository.claimed = []model.ClaimedOutboxAction{
		lifecycleActionClaim(t, model.OutboxActionCreate, action),
	}
	repository.states[channelID] = workerState(t, match.Snapshot{
		ChannelID: channelID,
		State:     match.StateInGame,
		InGame:    true,
		MatchID:   action.MatchID,
	})
	predictionActions := &fakeWorkerPredictions{calls: recorder}
	stats := &fakeWorkerStatsUpdater{calls: recorder}
	emitter := &fakeWorkerMatchEndedEmitter{calls: recorder}
	worker := newWorkerForTest(repository, predictionActions, stats, emitter, nil)

	require.NoError(t, worker.ProcessOnce(context.Background()))

	require.Equal(t, []dotarepository.ClaimPredictionActionsInput{{
		Limit: 10,
		Lease: 2 * time.Minute,
	}}, repository.claimedInputs())
	require.Equal(t, []match.LifecycleAction{action}, predictionActions.createCalls())
	require.Equal(t, []string{"create", "complete"}, recorder.all())
	completions := repository.completed()
	require.Len(t, completions, 1)
	require.Equal(t, repository.claimed[0].ID, completions[0].actionID)
	require.Equal(t, repository.claimed[0].LockToken, completions[0].lockToken)
	require.Empty(t, repository.retried())
}

func TestLifecycleWorkerProcessOnceRetriesDispatchFailureWithClaimToken(t *testing.T) {
	channelID := uuid.New()
	action := match.LifecycleAction{
		Kind:      match.ActionCreate,
		ChannelID: channelID,
		MatchID:   102,
		TeamKnown: true,
	}
	recorder := &workerCallRecorder{}
	repository := newFakeWorkerRepository()
	repository.calls = recorder
	repository.claimed = []model.ClaimedOutboxAction{
		lifecycleActionClaim(t, model.OutboxActionCreate, action),
	}
	repository.states[channelID] = workerState(t, match.Snapshot{
		ChannelID: channelID,
		State:     match.StateInGame,
		InGame:    true,
		MatchID:   action.MatchID,
	})
	predictionActions := &fakeWorkerPredictions{
		createErr: errors.New("Twitch temporarily unavailable"),
		calls:     recorder,
	}
	worker := newWorkerForTest(repository, predictionActions, &fakeWorkerStatsUpdater{}, &fakeWorkerMatchEndedEmitter{}, nil)
	before := time.Now()

	require.NoError(t, worker.ProcessOnce(context.Background()))

	require.Equal(t, []string{"create", "retry"}, recorder.all())
	require.Empty(t, repository.completed())
	retries := repository.retried()
	require.Len(t, retries, 1)
	require.Equal(t, repository.claimed[0].ID, retries[0].actionID)
	require.Equal(t, repository.claimed[0].LockToken, retries[0].lockToken)
	require.GreaterOrEqual(t, retries[0].availableAt, before.Add(time.Second))
}

func TestLifecycleWorkerProcessOnceCompletesStaleCreateWithoutDispatch(t *testing.T) {
	channelID := uuid.New()
	action := match.LifecycleAction{
		Kind:      match.ActionCreate,
		ChannelID: channelID,
		MatchID:   103,
		TeamKnown: true,
	}

	for _, test := range []struct {
		name     string
		snapshot match.Snapshot
	}{
		{
			name: "terminal current state",
			snapshot: match.Snapshot{
				ChannelID: channelID,
				State:     match.StateIdle,
			},
		},
		{
			name: "replaced current match",
			snapshot: match.Snapshot{
				ChannelID: channelID,
				State:     match.StateInGame,
				InGame:    true,
				MatchID:   104,
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			recorder := &workerCallRecorder{}
			repository := newFakeWorkerRepository()
			repository.calls = recorder
			repository.claimed = []model.ClaimedOutboxAction{
				lifecycleActionClaim(t, model.OutboxActionCreate, action),
			}
			repository.states[channelID] = workerState(t, test.snapshot)
			predictionActions := &fakeWorkerPredictions{calls: recorder}
			worker := newWorkerForTest(repository, predictionActions, &fakeWorkerStatsUpdater{}, &fakeWorkerMatchEndedEmitter{}, nil)

			require.NoError(t, worker.ProcessOnce(context.Background()))

			require.Empty(t, predictionActions.createCalls())
			require.Equal(t, []string{"complete"}, recorder.all())
			require.Len(t, repository.completed(), 1)
			require.Empty(t, repository.retried())
		})
	}
}

func TestLifecycleWorkerProcessOnceResolvesUpdatesStatsAndPublishesMatchEnded(t *testing.T) {
	channelID := uuid.New()
	action := match.LifecycleAction{
		Kind:           match.ActionResolve,
		ChannelID:      channelID,
		MatchID:        105,
		SteamAccountID: "12345",
		HeroName:       "axe",
		Win:            true,
	}
	settings := model.ChannelDotaSettings{
		ChannelID:     channelID,
		Mmr:           3_025,
		SessionWins:   4,
		SessionLosses: 2,
	}
	recorder := &workerCallRecorder{}
	repository := newFakeWorkerRepository()
	repository.calls = recorder
	repository.claimed = []model.ClaimedOutboxAction{
		lifecycleActionClaim(t, model.OutboxActionResolve, action),
	}
	predictionActions := &fakeWorkerPredictions{settings: settings, calls: recorder}
	stats := &fakeWorkerStatsUpdater{calls: recorder}
	emitter := &fakeWorkerMatchEndedEmitter{calls: recorder}
	worker := newWorkerForTest(repository, predictionActions, stats, emitter, nil)

	require.NoError(t, worker.ProcessOnce(context.Background()))

	require.Equal(t, []match.LifecycleAction{action}, predictionActions.resolveCalls())
	require.Equal(t, []workerStatsUpdate{{channelID: channelID, settings: settings}}, stats.all())
	require.Equal(t, []busdota.MatchEndedMessage{{
		ChannelID:      channelID.String(),
		SteamAccountID: "12345",
		Win:            true,
		HeroName:       "axe",
		Mmr:            3_025,
		SessionWins:    4,
		SessionLosses:  2,
		MatchID:        105,
	}}, emitter.all())
	require.Equal(t, []string{"resolve", "update_stats", "match_ended", "complete"}, recorder.all())
	require.Len(t, repository.completed(), 1)
	require.Empty(t, repository.retried())
}

func TestLifecycleWorkerRetriesStateUpdateBeforeMatchEndedAndCompletion(t *testing.T) {
	channelID := uuid.New()
	action := match.LifecycleAction{
		Kind:           match.ActionResolve,
		ChannelID:      channelID,
		MatchID:        110,
		SteamAccountID: "12345",
		HeroName:       "axe",
		Win:            true,
	}
	settings := model.ChannelDotaSettings{
		ChannelID:     channelID,
		Mmr:           3_025,
		SessionWins:   4,
		SessionLosses: 2,
	}
	stateUpdateErr := errors.New("state update unavailable")
	recorder := &workerCallRecorder{}
	repository := newFakeWorkerRepository()
	repository.calls = recorder
	firstClaim := lifecycleActionClaim(t, model.OutboxActionResolve, action)
	repository.claimed = []model.ClaimedOutboxAction{firstClaim}
	predictionActions := &fakeWorkerPredictions{settings: settings, calls: recorder}
	stateEmitter := &fakeWorkerStateUpdateEmitter{failNext: stateUpdateErr}
	stats := &fakeWorkerStatsUpdater{calls: recorder, stateEmitter: stateEmitter}
	emitter := &fakeWorkerMatchEndedEmitter{calls: recorder}
	worker := newWorkerForTest(repository, predictionActions, stats, emitter, nil)

	require.NoError(t, worker.ProcessOnce(context.Background()))

	require.Equal(t, []string{"resolve", "update_stats", "retry"}, recorder.all())
	require.Empty(t, emitter.all())
	require.Empty(t, repository.completed())
	retries := repository.retried()
	require.Len(t, retries, 1)
	require.Equal(t, firstClaim.LockToken, retries[0].lockToken)
	require.Empty(t, stateEmitter.all())

	replayClaim := firstClaim
	replayClaim.LockToken = uuid.New()
	replayClaim.Attempts++
	repository.claimed = []model.ClaimedOutboxAction{replayClaim}
	require.NoError(t, worker.ProcessOnce(context.Background()))

	require.Equal(t, []string{
		"resolve",
		"update_stats",
		"retry",
		"resolve",
		"update_stats",
		"match_ended",
		"complete",
	}, recorder.all())
	require.Len(t, predictionActions.resolveCalls(), 2)
	require.Len(t, stats.all(), 2)
	require.Len(t, stateEmitter.allAttempts(), 2)
	require.Equal(t, []busapi.DotaStateUpdateMessage{{
		ChannelID:     channelID.String(),
		Mmr:           3_025,
		SessionWins:   4,
		SessionLosses: 2,
	}}, stateEmitter.all())
	require.Len(t, emitter.all(), 1)
	completions := repository.completed()
	require.Len(t, completions, 1)
	require.Equal(t, replayClaim.LockToken, completions[0].lockToken)
}

func TestLifecycleWorkerRenewsLeaseWhileResolveIsInFlight(t *testing.T) {
	channelID := uuid.New()
	action := match.LifecycleAction{Kind: match.ActionResolve, ChannelID: channelID, MatchID: 111}
	settings := model.ChannelDotaSettings{ChannelID: channelID}
	repository := newFakeWorkerRepository()
	claim := lifecycleActionClaim(t, model.OutboxActionResolve, action)
	repository.claimed = []model.ClaimedOutboxAction{claim}
	renewed := make(chan struct{}, 16)
	repository.renewFn = func(context.Context, uuid.UUID, uuid.UUID, time.Duration) error {
		renewed <- struct{}{}
		return nil
	}
	resolveStarted := make(chan struct{}, 1)
	resolveRelease := make(chan struct{})
	predictionActions := &fakeWorkerPredictions{
		settings:       settings,
		resolveStarted: resolveStarted,
		resolveRelease: resolveRelease,
	}
	stats := &fakeWorkerStatsUpdater{}
	emitter := &fakeWorkerMatchEndedEmitter{}
	worker := newWorkerForTest(repository, predictionActions, stats, emitter, nil)
	worker.actionLease = 30 * time.Millisecond
	worker.actionRenewEvery = 5 * time.Millisecond

	done := make(chan error, 1)
	go func() {
		done <- worker.ProcessOnce(context.Background())
	}()
	<-resolveStarted
	for range int(worker.actionLease/worker.actionRenewEvery) + 1 {
		select {
		case <-renewed:
		case <-time.After(time.Second):
			t.Fatal("worker did not renew the in-flight action lease")
		}
	}
	close(resolveRelease)
	require.NoError(t, <-done)

	renewals := repository.renewed()
	require.GreaterOrEqual(t, len(renewals), int(worker.actionLease/worker.actionRenewEvery)+1)
	for _, renewal := range renewals {
		require.Equal(t, claim.ID, renewal.actionID)
		require.Equal(t, claim.LockToken, renewal.lockToken)
		require.Equal(t, worker.actionLease, renewal.lease)
	}
	require.Len(t, emitter.all(), 1)
	require.Len(t, repository.completed(), 1)
}

func TestLifecycleWorkerCancelsStaleResolveAfterRenewalOwnershipLoss(t *testing.T) {
	channelID := uuid.New()
	action := match.LifecycleAction{Kind: match.ActionResolve, ChannelID: channelID, MatchID: 112}
	repository := newFakeWorkerRepository()
	claim := lifecycleActionClaim(t, model.OutboxActionResolve, action)
	repository.claimed = []model.ClaimedOutboxAction{claim}
	repository.renewErr = dotarepository.ErrPredictionActionOwnershipLost
	repository.retryErr = dotarepository.ErrPredictionActionOwnershipLost
	resolveStarted := make(chan struct{}, 1)
	resolveCanceled := make(chan struct{}, 1)
	resolveRelease := make(chan struct{})
	predictionActions := &fakeWorkerPredictions{
		settings:             model.ChannelDotaSettings{ChannelID: channelID},
		resolveStarted:       resolveStarted,
		resolveCanceled:      resolveCanceled,
		resolveRelease:       resolveRelease,
		resolveWaitForCancel: true,
	}
	stats := &fakeWorkerStatsUpdater{}
	emitter := &fakeWorkerMatchEndedEmitter{}
	worker := newWorkerForTest(repository, predictionActions, stats, emitter, nil)
	worker.actionLease = 30 * time.Millisecond
	worker.actionRenewEvery = 5 * time.Millisecond

	done := make(chan error, 1)
	go func() {
		done <- worker.ProcessOnce(context.Background())
	}()
	<-resolveStarted
	select {
	case <-resolveCanceled:
	case <-time.After(time.Second):
		t.Fatal("renewal ownership loss did not cancel resolve")
	}
	close(resolveRelease)
	require.NoError(t, <-done)

	require.Len(t, repository.renewed(), 1)
	require.Empty(t, stats.all())
	require.Empty(t, emitter.all())
	require.Empty(t, repository.completed())
	require.Len(t, repository.retried(), 1)
}

func TestLifecycleWorkerRenewsLaterClaimBeforeSerialDispatch(t *testing.T) {
	firstAction := match.LifecycleAction{Kind: match.ActionResolve, ChannelID: uuid.New(), MatchID: 113}
	secondAction := match.LifecycleAction{Kind: match.ActionCancel, ChannelID: uuid.New(), MatchID: 114}
	firstClaim := lifecycleActionClaim(t, model.OutboxActionResolve, firstAction)
	secondClaim := lifecycleActionClaim(t, model.OutboxActionCancel, secondAction)
	repository := newFakeWorkerRepository()
	repository.claimed = []model.ClaimedOutboxAction{firstClaim, secondClaim}
	renewed := make(chan uuid.UUID, 32)
	repository.renewFn = func(_ context.Context, actionID, _ uuid.UUID, _ time.Duration) error {
		renewed <- actionID
		return nil
	}
	resolveStarted := make(chan struct{}, 1)
	resolveRelease := make(chan struct{})
	releaseResolve := func() {
		select {
		case <-resolveRelease:
		default:
			close(resolveRelease)
		}
	}
	defer releaseResolve()
	predictionActions := &fakeWorkerPredictions{
		resolveStarted: resolveStarted,
		resolveRelease: resolveRelease,
	}
	worker := newWorkerForTest(repository, predictionActions, &fakeWorkerStatsUpdater{}, &fakeWorkerMatchEndedEmitter{}, nil)
	worker.actionLease = 30 * time.Millisecond
	worker.actionRenewEvery = 5 * time.Millisecond

	done := make(chan error, 1)
	go func() {
		done <- worker.ProcessOnce(context.Background())
	}()
	<-resolveStarted

	deadline := time.NewTimer(time.Second)
	defer deadline.Stop()
	for {
		select {
		case actionID := <-renewed:
			if actionID == secondClaim.ID {
				require.Empty(t, predictionActions.cancelCalls())
				releaseResolve()
				require.NoError(t, <-done)
				return
			}
		case <-deadline.C:
			t.Fatal("later claimed action was not renewed before serial dispatch")
		}
	}
}

func TestLifecycleWorkerCancelsResolveWhenRenewalTimesOut(t *testing.T) {
	action := match.LifecycleAction{Kind: match.ActionResolve, ChannelID: uuid.New(), MatchID: 115}
	repository := newFakeWorkerRepository()
	claim := lifecycleActionClaim(t, model.OutboxActionResolve, action)
	repository.claimed = []model.ClaimedOutboxAction{claim}
	renewStarted := make(chan struct{}, 1)
	renewHasDeadline := make(chan bool, 1)
	renewRelease := make(chan struct{})
	defer close(renewRelease)
	repository.renewFn = func(ctx context.Context, _ uuid.UUID, _ uuid.UUID, _ time.Duration) error {
		_, hasDeadline := ctx.Deadline()
		renewHasDeadline <- hasDeadline
		renewStarted <- struct{}{}
		<-renewRelease
		return errors.New("renew request released")
	}
	resolveStarted := make(chan struct{}, 1)
	resolveCanceled := make(chan struct{}, 1)
	predictionActions := &fakeWorkerPredictions{
		resolveStarted:       resolveStarted,
		resolveCanceled:      resolveCanceled,
		resolveWaitForCancel: true,
	}
	stats := &fakeWorkerStatsUpdater{}
	emitter := &fakeWorkerMatchEndedEmitter{}
	worker := newWorkerForTest(repository, predictionActions, stats, emitter, nil)
	worker.actionLease = 30 * time.Millisecond
	worker.actionRenewEvery = 5 * time.Millisecond

	done := make(chan error, 1)
	go func() {
		done <- worker.ProcessOnce(context.Background())
	}()
	<-resolveStarted
	select {
	case <-renewStarted:
	case <-time.After(time.Second):
		t.Fatal("worker did not start renewal")
	}
	select {
	case hasDeadline := <-renewHasDeadline:
		if !hasDeadline {
			t.Fatal("renewal request did not receive a deadline")
		}
	case <-time.After(time.Second):
		t.Fatal("renewal request did not report its deadline")
	}
	select {
	case <-resolveCanceled:
	case <-time.After(time.Second):
		t.Fatal("renewal timeout did not cancel resolve")
	}
	select {
	case err := <-done:
		require.NoError(t, err)
	case <-time.After(time.Second):
		t.Fatal("stuck renewal blocked worker shutdown")
	}

	require.Empty(t, stats.all())
	require.Empty(t, emitter.all())
	require.Empty(t, repository.completed())
}

func TestLifecycleWorkerProcessOnceRetriesInvalidOrMismatchedPayloadWithoutDispatch(t *testing.T) {
	channelID := uuid.New()
	validAction := match.LifecycleAction{
		Kind:      match.ActionCreate,
		ChannelID: channelID,
		MatchID:   106,
		TeamKnown: true,
	}

	for _, test := range []struct {
		name   string
		mutate func(*model.ClaimedOutboxAction)
	}{
		{
			name: "invalid JSON",
			mutate: func(claim *model.ClaimedOutboxAction) {
				claim.Payload = json.RawMessage(`{`)
			},
		},
		{
			name: "mismatched action kind",
			mutate: func(claim *model.ClaimedOutboxAction) {
				payload, err := json.Marshal(match.LifecycleAction{
					Kind:      match.ActionResolve,
					ChannelID: claim.ChannelID,
					MatchID:   claim.MatchID,
				})
				require.NoError(t, err)
				claim.Payload = payload
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			recorder := &workerCallRecorder{}
			repository := newFakeWorkerRepository()
			repository.calls = recorder
			claim := lifecycleActionClaim(t, model.OutboxActionCreate, validAction)
			test.mutate(&claim)
			repository.claimed = []model.ClaimedOutboxAction{claim}
			predictionActions := &fakeWorkerPredictions{calls: recorder}
			worker := newWorkerForTest(repository, predictionActions, &fakeWorkerStatsUpdater{}, &fakeWorkerMatchEndedEmitter{}, nil)

			require.NoError(t, worker.ProcessOnce(context.Background()))

			require.Empty(t, predictionActions.createCalls())
			require.Empty(t, predictionActions.resolveCalls())
			require.Empty(t, predictionActions.cancelCalls())
			require.Equal(t, []string{"retry"}, recorder.all())
			require.Empty(t, repository.completed())
			retries := repository.retried()
			require.Len(t, retries, 1)
			require.Equal(t, claim.LockToken, retries[0].lockToken)
		})
	}
}

func TestLifecycleWorkerProcessOnceReschedulesTerminalInProgress(t *testing.T) {
	channelID := uuid.New()
	action := match.LifecycleAction{
		Kind:      match.ActionCancel,
		ChannelID: channelID,
		MatchID:   107,
	}
	recorder := &workerCallRecorder{}
	repository := newFakeWorkerRepository()
	repository.calls = recorder
	repository.claimed = []model.ClaimedOutboxAction{
		lifecycleActionClaim(t, model.OutboxActionCancel, action),
	}
	predictionActions := &fakeWorkerPredictions{
		cancelErr: ErrTerminalInProgress,
		calls:     recorder,
	}
	worker := newWorkerForTest(repository, predictionActions, &fakeWorkerStatsUpdater{}, &fakeWorkerMatchEndedEmitter{}, nil)

	require.NoError(t, worker.ProcessOnce(context.Background()))

	require.Equal(t, []match.LifecycleAction{action}, predictionActions.cancelCalls())
	require.Equal(t, []string{"cancel", "retry"}, recorder.all())
	require.Empty(t, repository.completed())
	retries := repository.retried()
	require.Len(t, retries, 1)
	require.Equal(t, repository.claimed[0].LockToken, retries[0].lockToken)
}

func TestLifecycleWorkerProcessOnceRetriesMatchEndedPublishFailure(t *testing.T) {
	channelID := uuid.New()
	action := match.LifecycleAction{
		Kind:      match.ActionResolve,
		ChannelID: channelID,
		MatchID:   108,
	}
	recorder := &workerCallRecorder{}
	repository := newFakeWorkerRepository()
	repository.calls = recorder
	repository.claimed = []model.ClaimedOutboxAction{
		lifecycleActionClaim(t, model.OutboxActionResolve, action),
	}
	predictionActions := &fakeWorkerPredictions{calls: recorder}
	stats := &fakeWorkerStatsUpdater{calls: recorder}
	emitter := &fakeWorkerMatchEndedEmitter{
		err:   errors.New("bus unavailable"),
		calls: recorder,
	}
	worker := newWorkerForTest(repository, predictionActions, stats, emitter, nil)

	require.NoError(t, worker.ProcessOnce(context.Background()))

	require.Equal(t, []string{"resolve", "update_stats", "match_ended", "retry"}, recorder.all())
	require.Empty(t, repository.completed())
	require.Len(t, repository.retried(), 1)
}

func TestLifecycleWorkerDeliversMatchEndedOnceAcrossLeaseLossReplay(t *testing.T) {
	action := match.LifecycleAction{
		Kind:           match.ActionResolve,
		ChannelID:      uuid.New(),
		MatchID:        116,
		SteamAccountID: "12345",
		HeroName:       "axe",
		Win:            true,
	}
	settings := model.ChannelDotaSettings{ChannelID: action.ChannelID}
	deliveryStore := newFakeWorkerMatchEndedDeliveryStore()
	emitter := &fakeWorkerMatchEndedEmitter{}

	firstClaim := lifecycleActionClaim(t, model.OutboxActionResolve, action)
	firstRepository := newFakeWorkerRepository()
	firstRepository.claimed = []model.ClaimedOutboxAction{firstClaim}
	firstRepository.completeErr = dotarepository.ErrPredictionActionOwnershipLost
	firstWorker := newWorkerForTest(
		firstRepository,
		&fakeWorkerPredictions{settings: settings},
		&fakeWorkerStatsUpdater{},
		emitter,
		nil,
	)
	firstWorker.deliveryStore = deliveryStore

	require.NoError(t, firstWorker.ProcessOnce(context.Background()))

	replayClaim := firstClaim
	replayClaim.LockToken = uuid.New()
	replayClaim.Attempts++
	replayRepository := newFakeWorkerRepository()
	replayRepository.claimed = []model.ClaimedOutboxAction{replayClaim}
	replayWorker := newWorkerForTest(
		replayRepository,
		&fakeWorkerPredictions{settings: settings},
		&fakeWorkerStatsUpdater{},
		emitter,
		nil,
	)
	replayWorker.deliveryStore = deliveryStore

	require.NoError(t, replayWorker.ProcessOnce(context.Background()))
	require.Len(t, emitter.all(), 1)
	require.Len(t, replayRepository.completed(), 1)
	require.Empty(t, replayRepository.retried())
	claims := deliveryStore.claimHistory()
	require.Len(t, claims, 2)
	require.Equal(t, matchEndedDeliveryKey(action.ChannelID, action.MatchID), claims[0].key)
	require.NotEqual(t, claims[0].token, claims[1].token)
	require.Equal(t, firstWorker.actionLease, claims[0].ttl)
}

func TestLifecycleWorkerReleasesMatchEndedDeliveryAfterPublishFailure(t *testing.T) {
	action := match.LifecycleAction{Kind: match.ActionResolve, ChannelID: uuid.New(), MatchID: 117}
	claim := lifecycleActionClaim(t, model.OutboxActionResolve, action)
	repository := newFakeWorkerRepository()
	repository.claimed = []model.ClaimedOutboxAction{claim}
	deliveryStore := newFakeWorkerMatchEndedDeliveryStore()
	emitter := &fakeWorkerMatchEndedEmitter{err: errors.New("Core NATS unavailable")}
	worker := newWorkerForTest(
		repository,
		&fakeWorkerPredictions{settings: model.ChannelDotaSettings{ChannelID: action.ChannelID}},
		&fakeWorkerStatsUpdater{},
		emitter,
		nil,
	)
	worker.deliveryStore = deliveryStore

	require.NoError(t, worker.ProcessOnce(context.Background()))
	require.Len(t, repository.retried(), 1)
	require.Empty(t, repository.completed())
	require.False(t, deliveryStore.hasClaim(matchEndedDeliveryKey(action.ChannelID, action.MatchID)))
	releases := deliveryStore.releaseHistory()
	require.Len(t, releases, 1)
	require.Equal(t, matchEndedDeliveryKey(action.ChannelID, action.MatchID), releases[0].key)
	claims := deliveryStore.claimHistory()
	require.Len(t, claims, 1)
	require.Equal(t, claims[0].token, releases[0].token)
}

func TestLifecycleWorkerRetriesPendingMatchEndedDeliveryUntilDelivered(t *testing.T) {
	action := match.LifecycleAction{
		Kind:      match.ActionResolve,
		ChannelID: uuid.New(),
		MatchID:   118,
	}
	settings := model.ChannelDotaSettings{ChannelID: action.ChannelID}
	deliveryStore := newFakeWorkerMatchEndedDeliveryStore()
	initialClaim := lifecycleActionClaim(t, model.OutboxActionResolve, action)
	firstRepository := newFakeWorkerRepository()
	firstRepository.claimed = []model.ClaimedOutboxAction{initialClaim}
	firstStarted := make(chan struct{}, 1)
	firstRelease := make(chan struct{})
	var releaseFirst sync.Once
	releaseFirstPublish := func() {
		releaseFirst.Do(func() { close(firstRelease) })
	}
	defer releaseFirstPublish()
	firstWorker := newWorkerForTest(
		firstRepository,
		&fakeWorkerPredictions{settings: settings},
		&fakeWorkerStatsUpdater{},
		&fakeWorkerMatchEndedEmitter{
			err:     errors.New("Core NATS unavailable"),
			started: firstStarted,
			release: firstRelease,
		},
		nil,
	)
	firstWorker.deliveryStore = deliveryStore
	firstDone := make(chan error, 1)
	go func() {
		firstDone <- firstWorker.ProcessOnce(context.Background())
	}()
	select {
	case <-firstStarted:
	case <-time.After(time.Second):
		t.Fatal("first worker did not reach MatchEnded publication")
	}

	key := matchEndedDeliveryKey(action.ChannelID, action.MatchID)
	require.True(t, deliveryStore.hasClaim(key))
	pendingClaim := initialClaim
	pendingClaim.LockToken = uuid.New()
	pendingClaim.Attempts++
	pendingRepository := newFakeWorkerRepository()
	pendingRepository.claimed = []model.ClaimedOutboxAction{pendingClaim}
	pendingEmitter := &fakeWorkerMatchEndedEmitter{}
	pendingWorker := newWorkerForTest(
		pendingRepository,
		&fakeWorkerPredictions{settings: settings},
		&fakeWorkerStatsUpdater{},
		pendingEmitter,
		nil,
	)
	pendingWorker.deliveryStore = deliveryStore

	require.NoError(t, pendingWorker.ProcessOnce(context.Background()))
	require.Len(t, pendingRepository.retried(), 1)
	require.Empty(t, pendingRepository.completed())
	require.Empty(t, pendingEmitter.all())

	releaseFirstPublish()
	require.NoError(t, <-firstDone)
	require.Len(t, firstRepository.retried(), 1)
	require.False(t, deliveryStore.hasClaim(key))

	retryClaim := pendingClaim
	retryClaim.LockToken = uuid.New()
	retryClaim.Attempts++
	retryRepository := newFakeWorkerRepository()
	retryRepository.claimed = []model.ClaimedOutboxAction{retryClaim}
	successEmitter := &fakeWorkerMatchEndedEmitter{}
	retryWorker := newWorkerForTest(
		retryRepository,
		&fakeWorkerPredictions{settings: settings},
		&fakeWorkerStatsUpdater{},
		successEmitter,
		nil,
	)
	retryWorker.deliveryStore = deliveryStore

	require.NoError(t, retryWorker.ProcessOnce(context.Background()))
	require.Len(t, successEmitter.all(), 1)
	require.Len(t, retryRepository.completed(), 1)
	require.True(t, deliveryStore.hasClaim(key))

	laterClaim := retryClaim
	laterClaim.LockToken = uuid.New()
	laterClaim.Attempts++
	laterRepository := newFakeWorkerRepository()
	laterRepository.claimed = []model.ClaimedOutboxAction{laterClaim}
	laterWorker := newWorkerForTest(
		laterRepository,
		&fakeWorkerPredictions{settings: settings},
		&fakeWorkerStatsUpdater{},
		successEmitter,
		nil,
	)
	laterWorker.deliveryStore = deliveryStore

	require.NoError(t, laterWorker.ProcessOnce(context.Background()))
	require.Len(t, successEmitter.all(), 1)
	require.Len(t, laterRepository.completed(), 1)
	require.Empty(t, laterRepository.retried())
}

func TestLifecycleWorkerRenewsActionAndPendingDeliveryLeasesWhileMatchEndedPublishIsInFlight(t *testing.T) {
	action := match.LifecycleAction{Kind: match.ActionResolve, ChannelID: uuid.New(), MatchID: 119}
	claim := lifecycleActionClaim(t, model.OutboxActionResolve, action)
	repository := newFakeWorkerRepository()
	repository.claimed = []model.ClaimedOutboxAction{claim}
	actionRenewed := make(chan struct{}, 16)
	repository.renewFn = func(context.Context, uuid.UUID, uuid.UUID, time.Duration) error {
		actionRenewed <- struct{}{}
		return nil
	}
	deliveryStore := newFakeWorkerMatchEndedDeliveryStore()
	deliveryRenewed := make(chan struct{}, 16)
	deliveryStore.renewFn = func(context.Context, string, string, time.Duration) (bool, error) {
		deliveryRenewed <- struct{}{}
		return true, nil
	}
	publishStarted := make(chan struct{}, 1)
	publishRelease := make(chan struct{})
	var releasePublishOnce sync.Once
	releasePublish := func() {
		releasePublishOnce.Do(func() { close(publishRelease) })
	}
	defer releasePublish()
	worker := newWorkerForTest(
		repository,
		&fakeWorkerPredictions{settings: model.ChannelDotaSettings{ChannelID: action.ChannelID}},
		&fakeWorkerStatsUpdater{},
		&fakeWorkerMatchEndedEmitter{started: publishStarted, release: publishRelease},
		nil,
	)
	worker.deliveryStore = deliveryStore
	worker.actionLease = 30 * time.Millisecond
	worker.actionRenewEvery = 5 * time.Millisecond

	done := make(chan error, 1)
	go func() {
		done <- worker.ProcessOnce(context.Background())
	}()
	select {
	case <-publishStarted:
	case <-time.After(time.Second):
		t.Fatal("worker did not start MatchEnded publication")
	}
	for range 2 {
		select {
		case <-actionRenewed:
		case <-time.After(time.Second):
			t.Fatal("worker did not renew the in-flight action lease")
		}
		select {
		case <-deliveryRenewed:
		case <-time.After(time.Second):
			t.Fatal("worker did not renew the pending delivery lease")
		}
	}

	releasePublish()
	require.NoError(t, <-done)
	claims := deliveryStore.claimHistory()
	require.Len(t, claims, 1)
	require.Equal(t, worker.actionLease, claims[0].ttl)
	renewals := deliveryStore.renewalHistory()
	require.GreaterOrEqual(t, len(renewals), 2)
	for _, renewal := range renewals {
		require.Equal(t, claims[0].key, renewal.key)
		require.Equal(t, claims[0].token, renewal.token)
		require.Equal(t, worker.actionLease, renewal.ttl)
	}
	completions := deliveryStore.completionHistory()
	require.Len(t, completions, 1)
	require.Equal(t, claims[0].key, completions[0].key)
	require.Equal(t, claims[0].token, completions[0].token)
	require.Equal(t, matchEndedDeliveryTTL, completions[0].ttl)
}

func TestLifecycleWorkerCancelsBlockedMatchEndedPublishAfterDeliveryRenewalOwnershipLoss(t *testing.T) {
	action := match.LifecycleAction{Kind: match.ActionResolve, ChannelID: uuid.New(), MatchID: 120}
	claim := lifecycleActionClaim(t, model.OutboxActionResolve, action)
	repository := newFakeWorkerRepository()
	repository.claimed = []model.ClaimedOutboxAction{claim}
	deliveryStore := newFakeWorkerMatchEndedDeliveryStore()
	deliveryRenewed := make(chan struct{}, 1)
	deliveryStore.renewFn = func(context.Context, string, string, time.Duration) (bool, error) {
		deliveryRenewed <- struct{}{}
		return false, nil
	}
	publishStarted := make(chan struct{}, 1)
	publishRelease := make(chan struct{})
	var releasePublishOnce sync.Once
	releasePublish := func() {
		releasePublishOnce.Do(func() { close(publishRelease) })
	}
	defer releasePublish()
	worker := newWorkerForTest(
		repository,
		&fakeWorkerPredictions{settings: model.ChannelDotaSettings{ChannelID: action.ChannelID}},
		&fakeWorkerStatsUpdater{},
		&fakeWorkerMatchEndedEmitter{started: publishStarted, release: publishRelease},
		nil,
	)
	worker.deliveryStore = deliveryStore
	worker.actionLease = 30 * time.Millisecond
	worker.actionRenewEvery = 5 * time.Millisecond

	done := make(chan error, 1)
	go func() {
		done <- worker.ProcessOnce(context.Background())
	}()
	select {
	case <-publishStarted:
	case <-time.After(time.Second):
		t.Fatal("worker did not start MatchEnded publication")
	}
	select {
	case <-deliveryRenewed:
	case <-time.After(time.Second):
		t.Fatal("worker did not renew the pending delivery lease")
	}
	select {
	case err := <-done:
		require.NoError(t, err)
	case <-time.After(time.Second):
		t.Fatal("delivery lease ownership loss did not cancel MatchEnded publication")
	}

	require.Empty(t, repository.completed())
	require.Len(t, repository.retried(), 1)
	require.Len(t, deliveryStore.releaseHistory(), 1)
	require.False(t, deliveryStore.hasClaim(matchEndedDeliveryKey(action.ChannelID, action.MatchID)))
}

func TestLifecycleWorkerProcessOnceIgnoresOwnershipLoss(t *testing.T) {
	channelID := uuid.New()
	action := match.LifecycleAction{
		Kind:      match.ActionCreate,
		ChannelID: channelID,
		MatchID:   109,
		TeamKnown: true,
	}

	for _, test := range []struct {
		name        string
		createErr   error
		completeErr error
		retryErr    error
		wantCalls   []string
	}{
		{
			name:        "while completing",
			completeErr: dotarepository.ErrPredictionActionOwnershipLost,
			wantCalls:   []string{"create", "complete"},
		},
		{
			name:      "while retrying",
			createErr: errors.New("Twitch unavailable"),
			retryErr:  dotarepository.ErrPredictionActionOwnershipLost,
			wantCalls: []string{"create", "retry"},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			recorder := &workerCallRecorder{}
			repository := newFakeWorkerRepository()
			repository.calls = recorder
			repository.completeErr = test.completeErr
			repository.retryErr = test.retryErr
			repository.claimed = []model.ClaimedOutboxAction{
				lifecycleActionClaim(t, model.OutboxActionCreate, action),
			}
			repository.states[channelID] = workerState(t, match.Snapshot{
				ChannelID: channelID,
				State:     match.StateInGame,
				InGame:    true,
				MatchID:   action.MatchID,
			})
			predictionActions := &fakeWorkerPredictions{createErr: test.createErr, calls: recorder}
			worker := newWorkerForTest(repository, predictionActions, &fakeWorkerStatsUpdater{}, &fakeWorkerMatchEndedEmitter{}, nil)

			require.NoError(t, worker.ProcessOnce(context.Background()))

			require.Equal(t, test.wantCalls, recorder.all())
		})
	}
}

func TestLifecycleWorkerRetryDelayCapsExponentialBackoff(t *testing.T) {
	for _, test := range []struct {
		attempts int
		want     time.Duration
	}{
		{attempts: 1, want: time.Second},
		{attempts: 2, want: 2 * time.Second},
		{attempts: 3, want: 4 * time.Second},
		{attempts: 7, want: time.Minute},
		{attempts: 100, want: time.Minute},
	} {
		t.Run(time.Duration(test.attempts).String(), func(t *testing.T) {
			require.Equal(t, test.want, retryDelay(test.attempts))
		})
	}
}

func TestLifecycleWorkerLifecycleStopsCancellablePollLoop(t *testing.T) {
	started := make(chan struct{}, 1)
	repository := newFakeWorkerRepository()
	repository.claimFn = func(ctx context.Context, _ dotarepository.ClaimPredictionActionsInput) ([]model.ClaimedOutboxAction, error) {
		started <- struct{}{}
		<-ctx.Done()
		return nil, ctx.Err()
	}
	lifecycle := &fakeWorkerLifecycle{}
	worker := newWorkerForTest(
		repository,
		&fakeWorkerPredictions{},
		&fakeWorkerStatsUpdater{},
		&fakeWorkerMatchEndedEmitter{},
		lifecycle,
	)

	require.NotNil(t, worker)
	require.Len(t, lifecycle.hooks, 1)
	require.NoError(t, lifecycle.hooks[0].OnStart(context.Background()))
	select {
	case <-started:
	case <-time.After(time.Second):
		t.Fatal("worker did not begin polling")
	}

	stopCtx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	require.NoError(t, lifecycle.hooks[0].OnStop(stopCtx))
}
