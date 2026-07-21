package predictions

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"strings"
	"sync"
	"testing"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/nicklaw5/helix/v2"
	"github.com/stretchr/testify/require"
	"github.com/twirapp/twir/apps/dota/internal/match"
	channelsmodel "github.com/twirapp/twir/libs/repositories/channels/model"
	dotarepository "github.com/twirapp/twir/libs/repositories/dota"
	dotamodel "github.com/twirapp/twir/libs/repositories/dota/model"
)

type fakeSettingsRepository struct {
	mu             sync.Mutex
	settings       dotamodel.ChannelDotaSettings
	err            error
	calls          int
	applyErr       error
	applyCalls     []dotarepository.ApplyMatchResultInput
	settledMatches map[int64]struct{}
}

func (r *fakeSettingsRepository) GetByChannelID(
	_ context.Context,
	_ uuid.UUID,
) (dotamodel.ChannelDotaSettings, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.calls++
	if r.err != nil {
		return dotamodel.Nil, r.err
	}
	return r.settings, nil
}

func (r *fakeSettingsRepository) Calls() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.calls
}

func (r *fakeSettingsRepository) ApplyMatchResultOnce(
	_ context.Context,
	input dotarepository.ApplyMatchResultInput,
) (dotamodel.ChannelDotaSettings, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.applyCalls = append(r.applyCalls, input)
	if r.applyErr != nil {
		return dotamodel.Nil, r.applyErr
	}
	if r.settledMatches == nil {
		r.settledMatches = make(map[int64]struct{})
	}
	if _, settled := r.settledMatches[input.MatchID]; !settled {
		r.settledMatches[input.MatchID] = struct{}{}
		r.settings.Mmr += input.MmrDelta
		if input.Won {
			r.settings.SessionWins++
		} else {
			r.settings.SessionLosses++
		}
	}

	return r.settings, nil
}

func (r *fakeSettingsRepository) ApplyCalls() []dotarepository.ApplyMatchResultInput {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]dotarepository.ApplyMatchResultInput(nil), r.applyCalls...)
}

func (r *fakeSettingsRepository) Settlements() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.settledMatches)
}

type fakeChannelsRepository struct {
	mu      sync.Mutex
	channel channelsmodel.Channel
	err     error
	calls   int
}

func (r *fakeChannelsRepository) GetByID(
	_ context.Context,
	_ uuid.UUID,
) (channelsmodel.Channel, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.calls++
	if r.err != nil {
		return channelsmodel.Nil, r.err
	}
	return r.channel, nil
}

type terminalClaim struct {
	token     string
	expiresAt time.Time
}

type matchEndedDeliveryRecord struct {
	token     string
	delivered bool
	expiresAt time.Time
}

type fakePredictionStore struct {
	mu sync.Mutex

	records      map[string]storedPrediction
	reservations map[string]string
	pending      map[string]pendingPredictionIntent
	terminal     map[string]terminalClaim
	deliveries   map[string]matchEndedDeliveryRecord

	reserveErr          error
	commitErr           error
	getErr              error
	releaseErr          error
	claimTerminalErr    error
	renewTerminalErr    error
	completeTerminalErr error
	releaseTerminalErr  error
	claimStarted        chan struct{}
	claimRelease        chan struct{}
	claimStartedOnce    sync.Once
	renewStarted        chan struct{}
	renewRelease        chan struct{}
	renewStartedOnce    sync.Once

	reserveCalls          []string
	commitCalls           []string
	releaseCalls          []string
	claimTerminalCalls    []string
	renewTerminalCalls    []string
	completeTerminalCalls []string
	releaseTerminalCalls  []string
}

func newFakePredictionStore() *fakePredictionStore {
	return &fakePredictionStore{
		records:      make(map[string]storedPrediction),
		reservations: make(map[string]string),
		pending:      make(map[string]pendingPredictionIntent),
		terminal:     make(map[string]terminalClaim),
		deliveries:   make(map[string]matchEndedDeliveryRecord),
	}
}

func (s *fakePredictionStore) Reserve(
	_ context.Context,
	key string,
	intent pendingPredictionIntent,
	_ time.Duration,
) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.reserveCalls = append(s.reserveCalls, key)
	if s.reserveErr != nil {
		return false, s.reserveErr
	}
	if _, exists := s.records[key]; exists {
		return false, nil
	}
	if _, exists := s.reservations[key]; exists {
		return false, nil
	}
	s.reservations[key] = intent.Token
	s.pending[key] = intent
	return true, nil
}

func (s *fakePredictionStore) Commit(
	_ context.Context,
	key string,
	token string,
	record storedPrediction,
	_ time.Duration,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.commitCalls = append(s.commitCalls, key)
	if s.commitErr != nil {
		return s.commitErr
	}
	if existing, exists := s.records[key]; exists {
		if existing == record {
			return nil
		}
		return errPredictionReservationLost
	}
	if s.reservations[key] != token {
		return errPredictionReservationLost
	}
	delete(s.reservations, key)
	delete(s.pending, key)
	s.records[key] = record
	return nil
}

func (s *fakePredictionStore) Get(_ context.Context, key string) (storedPrediction, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.getErr != nil {
		return storedPrediction{}, s.getErr
	}
	if _, pending := s.reservations[key]; pending {
		return storedPrediction{}, errPredictionPending
	}
	record, ok := s.records[key]
	if !ok {
		return storedPrediction{}, errPredictionNotFound
	}
	return record, nil
}

func (s *fakePredictionStore) GetPending(
	_ context.Context,
	key string,
) (pendingPredictionIntent, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	intent, ok := s.pending[key]
	if !ok {
		return pendingPredictionIntent{}, errPredictionIntentNotFound
	}
	return intent, nil
}

func (s *fakePredictionStore) Release(_ context.Context, key string, token string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.releaseCalls = append(s.releaseCalls, key)
	if s.releaseErr != nil {
		return s.releaseErr
	}
	if s.reservations[key] == token {
		delete(s.reservations, key)
		delete(s.pending, key)
	}
	return nil
}

func (s *fakePredictionStore) ClaimTerminal(
	_ context.Context,
	key string,
	token string,
	ttl time.Duration,
) (bool, error) {
	shouldBlock := false
	if s.claimStarted != nil {
		s.claimStartedOnce.Do(func() {
			shouldBlock = true
			close(s.claimStarted)
		})
	}
	if shouldBlock {
		<-s.claimRelease
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.claimTerminalCalls = append(s.claimTerminalCalls, key)
	if s.claimTerminalErr != nil {
		return false, s.claimTerminalErr
	}
	if _, claimed := s.activeTerminalClaim(key); claimed {
		return false, nil
	}
	s.terminal[key] = terminalClaim{token: token, expiresAt: time.Now().Add(ttl)}
	return true, nil
}

func (s *fakePredictionStore) RenewTerminal(
	ctx context.Context,
	key string,
	token string,
	ttl time.Duration,
) (bool, error) {
	s.mu.Lock()
	s.renewTerminalCalls = append(s.renewTerminalCalls, key)
	if s.renewTerminalErr != nil {
		s.mu.Unlock()
		return false, s.renewTerminalErr
	}
	claim, claimed := s.activeTerminalClaim(key)
	if !claimed || claim.token != token {
		s.mu.Unlock()
		return false, nil
	}
	claim.expiresAt = time.Now().Add(ttl)
	s.terminal[key] = claim
	started := s.renewStarted
	release := s.renewRelease
	s.mu.Unlock()

	if started != nil {
		s.renewStartedOnce.Do(func() { close(started) })
		select {
		case <-release:
		case <-ctx.Done():
			return false, ctx.Err()
		}
	}

	return true, nil
}

func (s *fakePredictionStore) CompleteTerminal(
	_ context.Context,
	key string,
	token string,
) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.completeTerminalCalls = append(s.completeTerminalCalls, key)
	if s.completeTerminalErr != nil {
		return false, s.completeTerminalErr
	}
	claim, claimed := s.activeTerminalClaim(key)
	if !claimed || claim.token != token {
		return false, nil
	}
	delete(s.records, key)
	delete(s.reservations, key)
	delete(s.pending, key)
	delete(s.terminal, key)
	return true, nil
}

func (s *fakePredictionStore) ReleaseTerminal(_ context.Context, key string, token string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.releaseTerminalCalls = append(s.releaseTerminalCalls, key)
	if s.releaseTerminalErr != nil {
		return s.releaseTerminalErr
	}
	claim, claimed := s.activeTerminalClaim(key)
	if claimed && claim.token == token {
		delete(s.terminal, key)
	}
	return nil
}

func (s *fakePredictionStore) ClaimMatchEndedDelivery(
	_ context.Context,
	key string,
	token string,
	ttl time.Duration,
) (matchEndedDeliveryState, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if claim, exists := s.deliveries[key]; exists && claim.expiresAt.After(time.Now()) {
		if claim.delivered {
			return matchEndedDeliveryDelivered, nil
		}
		return matchEndedDeliveryPending, nil
	}
	s.deliveries[key] = matchEndedDeliveryRecord{token: token, expiresAt: time.Now().Add(ttl)}
	return matchEndedDeliveryAcquired, nil
}

func (s *fakePredictionStore) CompleteMatchEndedDelivery(
	_ context.Context,
	key string,
	token string,
	ttl time.Duration,
) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	claim, exists := s.deliveries[key]
	if !exists || !claim.expiresAt.After(time.Now()) || claim.delivered || claim.token != token {
		return false, nil
	}
	claim.delivered = true
	claim.expiresAt = time.Now().Add(ttl)
	s.deliveries[key] = claim
	return true, nil
}

func (s *fakePredictionStore) RenewMatchEndedDelivery(
	_ context.Context,
	key string,
	token string,
	ttl time.Duration,
) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	claim, exists := s.deliveries[key]
	if !exists || !claim.expiresAt.After(time.Now()) || claim.delivered || claim.token != token {
		return false, nil
	}
	claim.expiresAt = time.Now().Add(ttl)
	s.deliveries[key] = claim
	return true, nil
}

func (s *fakePredictionStore) ReleaseMatchEndedDelivery(_ context.Context, key string, token string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if claim, exists := s.deliveries[key]; exists && !claim.delivered && claim.token == token {
		delete(s.deliveries, key)
	}
	return nil
}

func (s *fakePredictionStore) Record(key string) (storedPrediction, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	record, ok := s.records[key]
	return record, ok
}

func (s *fakePredictionStore) Store(key string, record storedPrediction) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.records[key] = record
}

func (s *fakePredictionStore) HasReservation(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.reservations[key]
	return ok
}

func (s *fakePredictionStore) HasTerminalClaim(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.activeTerminalClaim(key)
	return ok
}

func (s *fakePredictionStore) RenewTerminalCallCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.renewTerminalCalls)
}

func (s *fakePredictionStore) activeTerminalClaim(key string) (terminalClaim, bool) {
	claim, ok := s.terminal[key]
	if !ok {
		return terminalClaim{}, false
	}
	if !time.Now().Before(claim.expiresAt) {
		delete(s.terminal, key)
		return terminalClaim{}, false
	}
	return claim, true
}

type fakePredictionClient struct {
	mu sync.Mutex

	requestCtx context.Context

	createResponse *helix.PredictionsResponse
	createErr      error
	getResponse    *helix.PredictionsResponse
	getErr         error
	endResponse    *helix.PredictionsResponse
	endErr         error

	createCalls []*helix.CreatePredictionParams
	getCalls    []*helix.PredictionsParams
	endCalls    []*helix.EndPredictionParams

	createStarted  chan struct{}
	createRelease  chan struct{}
	getStarted     chan struct{}
	getRelease     chan struct{}
	getStartedOnce sync.Once
	endStarted     chan struct{}
	endRelease     chan struct{}
	endStartedOnce sync.Once
}

func (c *fakePredictionClient) CreatePrediction(
	params *helix.CreatePredictionParams,
) (*helix.PredictionsResponse, error) {
	c.mu.Lock()
	c.createCalls = append(c.createCalls, params)
	started := c.createStarted
	release := c.createRelease
	response := c.createResponse
	err := c.createErr
	c.mu.Unlock()

	if started != nil {
		close(started)
		<-release
	}
	return response, err
}

func (c *fakePredictionClient) GetPredictions(
	params *helix.PredictionsParams,
) (*helix.PredictionsResponse, error) {
	c.mu.Lock()
	c.getCalls = append(c.getCalls, params)
	requestCtx := c.requestCtx
	started := c.getStarted
	release := c.getRelease
	response := c.getResponse
	err := c.getErr
	c.mu.Unlock()

	if started != nil {
		c.getStartedOnce.Do(func() { close(started) })
		select {
		case <-release:
		case <-requestCtx.Done():
			return nil, requestCtx.Err()
		}
	}
	if requestCtx != nil {
		select {
		case <-requestCtx.Done():
			return nil, requestCtx.Err()
		default:
		}
	}

	return response, err
}

func (c *fakePredictionClient) EndPrediction(
	params *helix.EndPredictionParams,
) (*helix.PredictionsResponse, error) {
	c.mu.Lock()
	requestCtx := c.requestCtx
	c.mu.Unlock()
	if requestCtx != nil {
		select {
		case <-requestCtx.Done():
			return nil, requestCtx.Err()
		default:
		}
	}

	c.mu.Lock()
	c.endCalls = append(c.endCalls, params)
	started := c.endStarted
	release := c.endRelease
	response := c.endResponse
	err := c.endErr
	c.mu.Unlock()

	if started != nil {
		c.endStartedOnce.Do(func() { close(started) })
		<-release
	}

	return response, err
}

func (c *fakePredictionClient) CreateCalls() []*helix.CreatePredictionParams {
	c.mu.Lock()
	defer c.mu.Unlock()
	return append([]*helix.CreatePredictionParams(nil), c.createCalls...)
}

func (c *fakePredictionClient) GetCalls() []*helix.PredictionsParams {
	c.mu.Lock()
	defer c.mu.Unlock()
	return append([]*helix.PredictionsParams(nil), c.getCalls...)
}

func (c *fakePredictionClient) EndCalls() []*helix.EndPredictionParams {
	c.mu.Lock()
	defer c.mu.Unlock()
	return append([]*helix.EndPredictionParams(nil), c.endCalls...)
}

func (c *fakePredictionClient) setRequestContext(ctx context.Context) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.requestCtx = ctx
}

type fakeClientFactory struct {
	mu       sync.Mutex
	client   predictionClient
	err      error
	userIDs  []uuid.UUID
	contexts []context.Context
}

func (f *fakeClientFactory) New(ctx context.Context, userID uuid.UUID) (predictionClient, error) {
	f.mu.Lock()
	f.userIDs = append(f.userIDs, userID)
	f.contexts = append(f.contexts, ctx)
	client := f.client
	err := f.err
	f.mu.Unlock()
	if predictionClient, ok := client.(*fakePredictionClient); ok {
		predictionClient.setRequestContext(ctx)
	}
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (f *fakeClientFactory) LastContext() context.Context {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.contexts[len(f.contexts)-1]
}

type fixture struct {
	predictions *Predictions
	settings    *fakeSettingsRepository
	channels    *fakeChannelsRepository
	clients     *fakeClientFactory
	client      *fakePredictionClient
	store       *fakePredictionStore
	channelID   uuid.UUID
	twitchUser  uuid.UUID
	broadcaster string
}

func newFixture(t *testing.T) *fixture {
	return newFixtureWithTerminalLease(t, defaultTerminalLease())
}

func newFixtureWithTerminalLease(t *testing.T, terminalLease terminalLease) *fixture {
	t.Helper()

	channelID := uuid.New()
	twitchUser := uuid.New()
	broadcaster := "123456"
	client := &fakePredictionClient{
		createResponse: createdPredictionResponse("prediction-1", "yes-outcome", "no-outcome"),
		endResponse:    &helix.PredictionsResponse{},
	}
	settings := &fakeSettingsRepository{settings: dotamodel.ChannelDotaSettings{
		Enabled: true,
		PredictionSettings: dotamodel.PredictionSettings{
			Enabled:       true,
			TitleTemplate: "Will the streamer win?",
			WindowSeconds: 300,
		},
	}}
	channels := &fakeChannelsRepository{channel: channelsmodel.Channel{
		ID:               channelID,
		TwitchUserID:     &twitchUser,
		TwitchPlatformID: &broadcaster,
	}}
	clients := &fakeClientFactory{client: client}
	store := newFakePredictionStore()

	return &fixture{
		predictions: newPredictions(
			settings,
			channels,
			clients,
			store,
			slog.New(slog.NewTextHandler(io.Discard, nil)),
			terminalLease,
		),
		settings:    settings,
		channels:    channels,
		clients:     clients,
		client:      client,
		store:       store,
		channelID:   channelID,
		twitchUser:  twitchUser,
		broadcaster: broadcaster,
	}
}

func createAction(f *fixture, matchID int64) match.LifecycleAction {
	return match.LifecycleAction{
		Kind:      match.ActionCreate,
		ChannelID: f.channelID,
		MatchID:   matchID,
		TeamKnown: true,
	}
}

func resolveAction(f *fixture, matchID int64, win bool) match.LifecycleAction {
	return match.LifecycleAction{
		Kind:      match.ActionResolve,
		ChannelID: f.channelID,
		MatchID:   matchID,
		Win:       win,
	}
}

func cancelAction(f *fixture, matchID int64) match.LifecycleAction {
	return match.LifecycleAction{
		Kind:      match.ActionCancel,
		ChannelID: f.channelID,
		MatchID:   matchID,
	}
}

func createdPredictionResponse(id string, yesOutcomeID string, noOutcomeID string) *helix.PredictionsResponse {
	return &helix.PredictionsResponse{
		Data: helix.ManyPredictions{Predictions: []helix.Prediction{{
			ID: id,
			Outcomes: []helix.Outcomes{
				{ID: yesOutcomeID, Title: "Yes"},
				{ID: noOutcomeID, Title: "No"},
			},
		}}},
	}
}

func predictionResponse(id string, status string) *helix.PredictionsResponse {
	return &helix.PredictionsResponse{
		Data: helix.ManyPredictions{Predictions: []helix.Prediction{{
			ID:     id,
			Status: status,
		}}},
	}
}

func TestCreateCreatesAndStoresPrediction(t *testing.T) {
	f := newFixture(t)
	action := match.LifecycleAction{
		Kind:      match.ActionCreate,
		ChannelID: f.channelID,
		MatchID:   901,
		TeamKnown: true,
	}

	err := f.predictions.Create(context.Background(), action)

	require.NoError(t, err)
	require.Len(t, f.client.CreateCalls(), 1)
	record, ok := f.store.Record(predictionKey(f.channelID, action.MatchID))
	require.True(t, ok)
	require.Equal(t, storedPrediction{
		PredictionID: "prediction-1",
		YesOutcomeID: "yes-outcome",
		NoOutcomeID:  "no-outcome",
	}, record)
}

func TestResolveSettlesOnceAndReturnsSettingsOnReplay(t *testing.T) {
	f := newFixture(t)
	f.settings.settings.Mmr = 1_500
	f.settings.settings.MmrDelta = 25
	action := match.LifecycleAction{
		Kind:      match.ActionResolve,
		ChannelID: f.channelID,
		MatchID:   902,
		Win:       false,
	}
	key := predictionKey(f.channelID, action.MatchID)
	f.store.Store(key, storedPrediction{
		PredictionID: "prediction-902",
		YesOutcomeID: "yes-outcome",
		NoOutcomeID:  "no-outcome",
	})
	f.client.getResponse = predictionResponse("prediction-902", "ACTIVE")

	settings, err := f.predictions.Resolve(context.Background(), action)

	require.NoError(t, err)
	require.Equal(t, 1_475, settings.Mmr)
	require.Equal(t, 1, settings.SessionLosses)
	require.Equal(t, []dotarepository.ApplyMatchResultInput{{
		ChannelID: f.channelID,
		MatchID:   action.MatchID,
		Won:       false,
		MmrDelta:  -25,
	}}, f.settings.ApplyCalls())
	require.Len(t, f.client.EndCalls(), 1)

	replayedSettings, err := f.predictions.Resolve(context.Background(), action)

	require.NoError(t, err)
	require.Equal(t, settings, replayedSettings)
	require.Equal(t, 1, f.settings.Settlements())
	require.Len(t, f.settings.ApplyCalls(), 2)
	require.Len(t, f.client.EndCalls(), 1)
}

func TestResolveAndCancelRaceEndPredictionOnce(t *testing.T) {
	f := newFixture(t)
	key := predictionKey(f.channelID, 903)
	f.store.Store(key, storedPrediction{
		PredictionID: "prediction-903",
		YesOutcomeID: "yes-outcome",
		NoOutcomeID:  "no-outcome",
	})
	f.client.getResponse = predictionResponse("prediction-903", "ACTIVE")
	f.client.endStarted = make(chan struct{})
	f.client.endRelease = make(chan struct{})

	var releaseOnce sync.Once
	releaseEnd := func() {
		releaseOnce.Do(func() { close(f.client.endRelease) })
	}
	defer releaseEnd()

	results := make(chan error, 2)
	go func() {
		_, err := f.predictions.Resolve(context.Background(), match.LifecycleAction{
			Kind:      match.ActionResolve,
			ChannelID: f.channelID,
			MatchID:   903,
			Win:       true,
		})
		results <- err
	}()
	go func() {
		results <- f.predictions.Cancel(context.Background(), match.LifecycleAction{
			Kind:      match.ActionCancel,
			ChannelID: f.channelID,
			MatchID:   903,
		})
	}()

	<-f.client.endStarted
	select {
	case err := <-results:
		require.ErrorIs(t, err, ErrTerminalInProgress)
	case <-time.After(time.Second):
		t.Fatal("terminal loser did not return while the winner held EndPrediction")
	}

	releaseEnd()
	require.NoError(t, <-results)
	require.Len(t, f.client.EndCalls(), 1)
}

func TestTerminalClaimPrecedesRecordReadAfterConcurrentCompletion(t *testing.T) {
	f := newFixture(t)
	key := predictionKey(f.channelID, 909)
	f.store.Store(key, storedPrediction{
		PredictionID: "prediction-909",
		YesOutcomeID: "yes-outcome",
		NoOutcomeID:  "no-outcome",
	})
	f.store.claimStarted = make(chan struct{})
	f.store.claimRelease = make(chan struct{})
	f.client.getResponse = predictionResponse("prediction-909", "ACTIVE")

	var releaseClaimOnce sync.Once
	releaseClaim := func() {
		releaseClaimOnce.Do(func() { close(f.store.claimRelease) })
	}
	defer releaseClaim()

	firstDone := make(chan error, 1)
	go func() {
		firstDone <- f.predictions.Cancel(context.Background(), cancelAction(f, 909))
	}()
	<-f.store.claimStarted

	secondDone := make(chan error, 1)
	go func() {
		secondDone <- f.predictions.Cancel(context.Background(), cancelAction(f, 909))
	}()
	require.NoError(t, <-secondDone)

	releaseClaim()
	require.NoError(t, <-firstDone)
	require.Len(t, f.client.EndCalls(), 1)
	_, exists := f.store.Record(key)
	require.False(t, exists)
}

func TestTerminalHeartbeatFailureCancelsClientOperationBeforeEnd(t *testing.T) {
	lease := terminalLease{
		ttl:           100 * time.Millisecond,
		renewInterval: 10 * time.Millisecond,
	}
	f := newFixtureWithTerminalLease(t, lease)
	key := predictionKey(f.channelID, 910)
	f.store.Store(key, storedPrediction{
		PredictionID: "prediction-910",
		YesOutcomeID: "yes-outcome",
		NoOutcomeID:  "no-outcome",
	})
	f.store.renewTerminalErr = errors.New("redis temporarily unavailable")
	f.client.getResponse = predictionResponse("prediction-910", "ACTIVE")
	f.client.getStarted = make(chan struct{})
	f.client.getRelease = make(chan struct{})

	var releaseGetOnce sync.Once
	releaseGet := func() {
		releaseGetOnce.Do(func() { close(f.client.getRelease) })
	}
	defer releaseGet()

	result := make(chan error, 1)
	go func() {
		result <- f.predictions.Cancel(context.Background(), cancelAction(f, 910))
	}()
	<-f.client.getStarted

	operationCtx := f.clients.LastContext()
	select {
	case <-operationCtx.Done():
	case <-time.After(time.Second):
		t.Fatal("heartbeat failure did not cancel the client operation context")
	}

	err := <-result
	require.ErrorIs(t, err, context.Canceled)
	require.Empty(t, f.client.EndCalls())
	_, exists := f.store.Record(key)
	require.True(t, exists)
}

func TestResolveReleasesTerminalClaimAfterTransientFailureAndRetries(t *testing.T) {
	f := newFixture(t)
	f.settings.settings.MmrDelta = 25
	action := match.LifecycleAction{
		Kind:      match.ActionResolve,
		ChannelID: f.channelID,
		MatchID:   904,
		Win:       true,
	}
	key := predictionKey(f.channelID, action.MatchID)
	f.store.Store(key, storedPrediction{
		PredictionID: "prediction-904",
		YesOutcomeID: "yes-outcome",
		NoOutcomeID:  "no-outcome",
	})
	f.client.getResponse = predictionResponse("prediction-904", "ACTIVE")
	f.client.endErr = errors.New("twitch temporarily unavailable")

	settings, err := f.predictions.Resolve(context.Background(), action)

	require.ErrorIs(t, err, f.client.endErr)
	require.Equal(t, 25, settings.Mmr)
	require.True(t, f.store.HasTerminalClaim(key) == false)
	_, exists := f.store.Record(key)
	require.True(t, exists)

	f.client.endErr = nil
	settings, err = f.predictions.Resolve(context.Background(), action)

	require.NoError(t, err)
	require.Equal(t, 25, settings.Mmr)
	require.Equal(t, 1, f.settings.Settlements())
	require.Len(t, f.client.EndCalls(), 2)
	_, exists = f.store.Record(key)
	require.False(t, exists)
}

func TestResolveReleasesTerminalClaimAfterTransientLookupFailureAndRetries(t *testing.T) {
	f := newFixture(t)
	action := resolveAction(f, 905, true)
	key := predictionKey(f.channelID, action.MatchID)
	f.store.Store(key, storedPrediction{
		PredictionID: "prediction-905",
		YesOutcomeID: "yes-outcome",
		NoOutcomeID:  "no-outcome",
	})
	f.client.getResponse = predictionResponse("prediction-905", "ACTIVE")
	f.client.getErr = errors.New("twitch temporarily unavailable")

	_, err := f.predictions.Resolve(context.Background(), action)

	require.ErrorIs(t, err, f.client.getErr)
	require.Equal(t, []string{key}, f.store.releaseTerminalCalls)
	require.Empty(t, f.client.EndCalls())
	_, exists := f.store.Record(key)
	require.True(t, exists)

	f.client.getErr = nil
	_, err = f.predictions.Resolve(context.Background(), action)

	require.NoError(t, err)
	require.Len(t, f.client.GetCalls(), 2)
	require.Len(t, f.client.EndCalls(), 1)
	_, exists = f.store.Record(key)
	require.False(t, exists)
}

func TestTerminalLeaseRenewalPreventsConcurrentEndAfterClaimExpiry(t *testing.T) {
	lease := terminalLease{
		ttl:           100 * time.Millisecond,
		renewInterval: 10 * time.Millisecond,
	}
	f := newFixtureWithTerminalLease(t, lease)
	key := predictionKey(f.channelID, 906)
	f.store.Store(key, storedPrediction{
		PredictionID: "prediction-906",
		YesOutcomeID: "yes-outcome",
		NoOutcomeID:  "no-outcome",
	})
	f.client.getResponse = predictionResponse("prediction-906", "ACTIVE")
	f.client.endStarted = make(chan struct{})
	f.client.endRelease = make(chan struct{})

	var releaseOnce sync.Once
	releaseEnd := func() {
		releaseOnce.Do(func() { close(f.client.endRelease) })
	}
	defer releaseEnd()

	firstDone := make(chan error, 1)
	go func() {
		_, err := f.predictions.Resolve(context.Background(), resolveAction(f, 906, true))
		firstDone <- err
	}()
	<-f.client.endStarted

	deadline := time.After(time.Second)
	for f.store.RenewTerminalCallCount() == 0 {
		select {
		case <-deadline:
			t.Fatal("terminal lease was not renewed while EndPrediction was blocked")
		case <-time.After(lease.renewInterval):
		}
	}
	time.Sleep(lease.ttl * 2)
	require.True(t, f.store.HasTerminalClaim(key))

	secondDone := make(chan error, 1)
	go func() {
		secondDone <- f.predictions.Cancel(context.Background(), cancelAction(f, 906))
	}()
	select {
	case err := <-secondDone:
		require.ErrorIs(t, err, ErrTerminalInProgress)
	case <-time.After(time.Second):
		t.Fatal("terminal caller did not yield while another caller held the renewed lease")
	}
	require.Len(t, f.client.EndCalls(), 1)

	releaseEnd()
	require.NoError(t, <-firstDone)
}

func TestTerminalHeartbeatShutdownWaitsForRenewalWithoutFailure(t *testing.T) {
	lease := terminalLease{
		ttl:           100 * time.Millisecond,
		renewInterval: 10 * time.Millisecond,
	}
	f := newFixtureWithTerminalLease(t, lease)
	key := predictionKey(f.channelID, 907)
	f.store.Store(key, storedPrediction{
		PredictionID: "prediction-907",
		YesOutcomeID: "yes-outcome",
		NoOutcomeID:  "no-outcome",
	})
	f.client.getResponse = predictionResponse("prediction-907", "ACTIVE")
	f.client.endStarted = make(chan struct{})
	f.client.endRelease = make(chan struct{})
	f.store.renewStarted = make(chan struct{})
	f.store.renewRelease = make(chan struct{})

	var releaseEndOnce sync.Once
	releaseEnd := func() {
		releaseEndOnce.Do(func() { close(f.client.endRelease) })
	}
	defer releaseEnd()
	var releaseRenewOnce sync.Once
	releaseRenew := func() {
		releaseRenewOnce.Do(func() { close(f.store.renewRelease) })
	}
	defer releaseRenew()

	result := make(chan error, 1)
	go func() {
		_, err := f.predictions.Resolve(context.Background(), resolveAction(f, 907, true))
		result <- err
	}()
	<-f.client.endStarted
	<-f.store.renewStarted

	releaseEnd()
	select {
	case err := <-result:
		t.Fatalf("terminal completion returned before renewal stopped: %v", err)
	case <-time.After(50 * time.Millisecond):
	}

	releaseRenew()
	require.NoError(t, <-result)
	_, exists := f.store.Record(key)
	require.False(t, exists)
}

func TestTerminalHeartbeatHonorsOperationDeadline(t *testing.T) {
	f := newFixture(t)
	key := predictionKey(f.channelID, 908)
	f.store.Store(key, storedPrediction{
		PredictionID: "prediction-908",
		YesOutcomeID: "yes-outcome",
		NoOutcomeID:  "no-outcome",
	})
	f.client.getResponse = predictionResponse("prediction-908", "ACTIVE")
	f.client.endStarted = make(chan struct{})
	f.client.endRelease = make(chan struct{})

	var releaseEndOnce sync.Once
	releaseEnd := func() {
		releaseEndOnce.Do(func() { close(f.client.endRelease) })
	}
	defer releaseEnd()

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	result := make(chan error, 1)
	go func() {
		_, err := f.predictions.Resolve(ctx, resolveAction(f, 908, true))
		result <- err
	}()
	<-f.client.endStarted
	<-ctx.Done()

	releaseEnd()
	err := <-result
	require.ErrorIs(t, err, context.DeadlineExceeded)
	_, exists := f.store.Record(key)
	require.True(t, exists)
}

func predictionCorrelationFromTitle(t *testing.T, template string, title string) string {
	t.Helper()

	const markerPrefix = " [d:"
	require.True(t, strings.HasPrefix(title, template+markerPrefix))
	require.True(t, strings.HasSuffix(title, "]"))

	correlation := strings.TrimSuffix(strings.TrimPrefix(title, template+markerPrefix), "]")
	decoded, err := base64.RawURLEncoding.DecodeString(correlation)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(decoded), 8)

	return correlation
}

func TestCreateUsesConfiguredTwitchParameters(t *testing.T) {
	f := newFixture(t)

	err := f.predictions.Create(context.Background(), createAction(f, 91))

	require.NoError(t, err)
	createCalls := f.client.CreateCalls()
	require.Len(t, createCalls, 1)
	require.Equal(t, f.broadcaster, createCalls[0].BroadcasterID)
	require.Equal(t, 300, createCalls[0].PredictionWindow)
	require.Equal(t, []helix.PredictionChoiceParam{
		{Title: "Yes"},
		{Title: "No"},
	}, createCalls[0].Outcomes)
	predictionCorrelationFromTitle(t, f.settings.settings.PredictionSettings.TitleTemplate, createCalls[0].Title)
	require.LessOrEqual(t, utf8.RuneCountInString(createCalls[0].Title), maxPredictionTitleRunes)

	record, ok := f.store.Record(predictionKey(f.channelID, 91))
	require.True(t, ok)
	require.Equal(t, storedPrediction{
		PredictionID: "prediction-1",
		YesOutcomeID: "yes-outcome",
		NoOutcomeID:  "no-outcome",
	}, record)
	require.Equal(t, []uuid.UUID{f.twitchUser}, f.clients.userIDs)
}

func TestCreatePersistsMarkedTitleAndCorrelationBeforeCreate(t *testing.T) {
	f := newFixture(t)
	f.store.commitErr = errors.New("redis temporarily unavailable")
	action := createAction(f, 911)
	key := predictionKey(f.channelID, action.MatchID)

	err := f.predictions.Create(context.Background(), action)

	require.ErrorIs(t, err, f.store.commitErr)
	createCalls := f.client.CreateCalls()
	require.Len(t, createCalls, 1)
	correlation := predictionCorrelationFromTitle(
		t,
		f.settings.settings.PredictionSettings.TitleTemplate,
		createCalls[0].Title,
	)
	intent, pendingErr := f.store.GetPending(context.Background(), key)
	require.NoError(t, pendingErr)
	require.Equal(t, createCalls[0].Title, intent.Title)
	require.Equal(t, correlation, intent.Correlation)

	encodedIntent, marshalErr := json.Marshal(intent)
	require.NoError(t, marshalErr)
	var persisted map[string]any
	require.NoError(t, json.Unmarshal(encodedIntent, &persisted))
	require.Equal(t, correlation, persisted["correlation"])
}

func TestCreateSkipsIneligibleInputs(t *testing.T) {
	tests := []struct {
		name      string
		configure func(*fixture, *match.LifecycleAction)
	}{
		{
			name: "disabled dota module",
			configure: func(f *fixture, _ *match.LifecycleAction) {
				f.settings.settings.Enabled = false
			},
		},
		{
			name: "disabled prediction settings",
			configure: func(f *fixture, _ *match.LifecycleAction) {
				f.settings.settings.PredictionSettings.Enabled = false
			},
		},
		{
			name: "missing match id",
			configure: func(_ *fixture, action *match.LifecycleAction) {
				action.MatchID = 0
			},
		},
		{
			name: "missing channel id",
			configure: func(_ *fixture, action *match.LifecycleAction) {
				action.ChannelID = uuid.Nil
			},
		},
		{
			name: "unknown team",
			configure: func(_ *fixture, action *match.LifecycleAction) {
				action.TeamKnown = false
			},
		},
		{
			name: "blank title",
			configure: func(f *fixture, _ *match.LifecycleAction) {
				f.settings.settings.PredictionSettings.TitleTemplate = " \t"
			},
		},
		{
			name: "zero window",
			configure: func(f *fixture, _ *match.LifecycleAction) {
				f.settings.settings.PredictionSettings.WindowSeconds = 0
			},
		},
		{
			name: "window above helix limit",
			configure: func(f *fixture, _ *match.LifecycleAction) {
				f.settings.settings.PredictionSettings.WindowSeconds = 1_801
			},
		},
		{
			name: "disconnected channel",
			configure: func(f *fixture, _ *match.LifecycleAction) {
				f.channels.channel = channelsmodel.Channel{}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := newFixture(t)
			action := createAction(f, 92)
			tt.configure(f, &action)

			err := f.predictions.Create(context.Background(), action)

			require.NoError(t, err)
			require.Empty(t, f.client.CreateCalls())
			require.Empty(t, f.store.reserveCalls)
		})
	}
}

func TestCreateEnforcesTwitchPredictionWindowLimits(t *testing.T) {
	for _, tt := range []struct {
		name          string
		windowSeconds int
		creates       bool
	}{
		{name: "minimum accepted", windowSeconds: 30, creates: true},
		{name: "below minimum rejected", windowSeconds: 29, creates: false},
		{name: "maximum accepted", windowSeconds: 1_800, creates: true},
		{name: "above maximum rejected", windowSeconds: 1_801, creates: false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			f := newFixture(t)
			f.settings.settings.PredictionSettings.WindowSeconds = tt.windowSeconds

			err := f.predictions.Create(context.Background(), createAction(f, 921))

			require.NoError(t, err)
			if tt.creates {
				require.Len(t, f.client.CreateCalls(), 1)
			} else {
				require.Empty(t, f.client.CreateCalls())
				require.Empty(t, f.store.reserveCalls)
			}
		})
	}
}

func TestCreateEnforcesTwitchPredictionTitleRuneLimits(t *testing.T) {
	const correlationMarkerRunes = 16
	longestTemplateRunes := maxPredictionTitleRunes - correlationMarkerRunes

	for _, tt := range []struct {
		name    string
		title   string
		creates bool
	}{
		{
			name:    "longest template that fits the marker is accepted",
			title:   strings.Repeat("界", longestTemplateRunes),
			creates: true,
		},
		{
			name:    "one rune over the marker boundary is rejected",
			title:   strings.Repeat("界", longestTemplateRunes+1),
			creates: false,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			f := newFixture(t)
			f.settings.settings.PredictionSettings.TitleTemplate = tt.title

			err := f.predictions.Create(context.Background(), createAction(f, 922))

			require.NoError(t, err)
			if tt.creates {
				require.Len(t, f.client.CreateCalls(), 1)
			} else {
				require.Empty(t, f.client.CreateCalls())
				require.Empty(t, f.store.reserveCalls)
			}
		})
	}
}

func TestCreateRecoversConcurrentPendingActionAndDeduplicatesReplay(t *testing.T) {
	f := newFixture(t)
	action := createAction(f, 93)
	f.client.createStarted = make(chan struct{})
	f.client.createRelease = make(chan struct{})

	var releaseOnce sync.Once
	releaseCreate := func() {
		releaseOnce.Do(func() { close(f.client.createRelease) })
	}
	defer releaseCreate()

	firstResult := make(chan error, 1)
	go func() {
		firstResult <- f.predictions.Create(context.Background(), action)
	}()
	<-f.client.createStarted

	intent, err := f.store.GetPending(context.Background(), predictionKey(f.channelID, action.MatchID))
	require.NoError(t, err)
	f.client.getResponse = predictionsResponse(activePredictionForIntent(intent, "prediction-1"))

	secondResult := make(chan error, 1)
	go func() {
		secondResult <- f.predictions.Create(context.Background(), action)
	}()
	require.NoError(t, <-secondResult)

	releaseCreate()
	require.NoError(t, <-firstResult)

	err = f.predictions.Create(context.Background(), action)
	require.NoError(t, err)
	require.Len(t, f.client.CreateCalls(), 1)
	require.Len(t, f.store.commitCalls, 2)
}

func TestCreateTreatsAlreadyActivePredictionAsNoopAndReleasesReservation(t *testing.T) {
	f := newFixture(t)
	f.client.createResponse = &helix.PredictionsResponse{
		ResponseCommon: helix.ResponseCommon{ErrorMessage: "Prediction is ALREADY ACTIVE"},
	}

	err := f.predictions.Create(context.Background(), createAction(f, 94))

	require.NoError(t, err)
	key := predictionKey(f.channelID, 94)
	require.Equal(t, []string{key}, f.store.releaseCalls)
	require.False(t, f.store.HasReservation(key))
	_, exists := f.store.Record(key)
	require.False(t, exists)
}

func TestCreateTreatsAlreadyActiveCreateErrorAsNoopAndReleasesReservation(t *testing.T) {
	f := newFixture(t)
	f.client.createResponse = nil
	f.client.createErr = errors.New("prediction is already active")

	err := f.predictions.Create(context.Background(), createAction(f, 941))

	require.NoError(t, err)
	key := predictionKey(f.channelID, 941)
	require.Equal(t, []string{key}, f.store.releaseCalls)
	require.False(t, f.store.HasReservation(key))
}

func TestCreateReleasesReservationAfterTransientCreateFailure(t *testing.T) {
	f := newFixture(t)
	f.client.createErr = errors.New("twitch temporarily unavailable")
	action := createAction(f, 95)

	err := f.predictions.Create(context.Background(), action)

	require.Error(t, err)
	key := predictionKey(f.channelID, 95)
	require.False(t, f.store.HasReservation(key))
	require.Equal(t, []string{key}, f.store.releaseCalls)

	f.client.createErr = nil
	err = f.predictions.Create(context.Background(), action)
	require.NoError(t, err)
	_, exists := f.store.Record(key)
	require.True(t, exists)
}

func TestCreateRetainsPendingIntentWhenCommitFailsAndRecoversOnReplay(t *testing.T) {
	f := newFixture(t)
	action := createAction(f, 951)
	key := predictionKey(f.channelID, action.MatchID)
	f.store.commitErr = errors.New("redis temporarily unavailable")

	err := f.predictions.Create(context.Background(), action)

	require.ErrorIs(t, err, f.store.commitErr)
	require.True(t, f.store.HasReservation(key))
	require.Len(t, f.client.CreateCalls(), 1)
	require.Empty(t, f.store.releaseCalls)

	f.store.commitErr = nil
	intent, pendingErr := f.store.GetPending(context.Background(), key)
	require.NoError(t, pendingErr)
	f.client.getResponse = predictionsResponse(activePredictionForIntent(intent, "prediction-1"))
	err = f.predictions.Create(context.Background(), action)

	require.NoError(t, err)
	require.Len(t, f.client.CreateCalls(), 1)
	record, exists := f.store.Record(key)
	require.True(t, exists)
	require.Equal(t, storedPrediction{
		PredictionID: "prediction-1",
		YesOutcomeID: "yes-outcome",
		NoOutcomeID:  "no-outcome",
	}, record)
}

func TestTerminalActionsRecoverKnownPredictionAfterCommitFailure(t *testing.T) {
	for _, tt := range []struct {
		name   string
		handle func(*fixture) error
		status string
	}{
		{
			name: "resolve",
			handle: func(f *fixture) error {
				_, err := f.predictions.Resolve(context.Background(), resolveAction(f, 952, true))
				return err
			},
			status: "RESOLVED",
		},
		{
			name: "cancel",
			handle: func(f *fixture) error {
				return f.predictions.Cancel(context.Background(), cancelAction(f, 952))
			},
			status: "CANCELED",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			f := newFixture(t)
			action := createAction(f, 952)
			f.store.commitErr = errors.New("redis temporarily unavailable")

			err := f.predictions.Create(context.Background(), action)
			require.ErrorIs(t, err, f.store.commitErr)
			require.Len(t, f.client.CreateCalls(), 1)

			f.store.commitErr = nil
			intent, pendingErr := f.store.GetPending(context.Background(), predictionKey(f.channelID, action.MatchID))
			require.NoError(t, pendingErr)
			f.client.getResponse = predictionsResponse(activePredictionForIntent(intent, "prediction-1"))
			require.NoError(t, tt.handle(f))
			require.Len(t, f.client.CreateCalls(), 1)
			endCalls := f.client.EndCalls()
			require.Len(t, endCalls, 1)
			require.Equal(t, tt.status, endCalls[0].Status)
			_, exists := f.store.Record(predictionKey(f.channelID, action.MatchID))
			require.False(t, exists)
		})
	}
}

func TestTerminalActionRecoversPredictionOnAnotherReplicaAfterCommitFailure(t *testing.T) {
	f := newFixture(t)
	action := createAction(f, 953)
	f.store.commitErr = errors.New("redis temporarily unavailable")

	err := f.predictions.Create(context.Background(), action)
	require.ErrorIs(t, err, f.store.commitErr)
	require.Len(t, f.client.CreateCalls(), 1)

	replicaClient := newReplicaClient()
	replica := newReplica(f, replicaClient)
	f.store.commitErr = nil
	intent, pendingErr := f.store.GetPending(context.Background(), predictionKey(f.channelID, action.MatchID))
	require.NoError(t, pendingErr)
	replicaClient.getResponse = predictionsResponse(activePredictionForIntent(intent, "prediction-1"))

	_, err = replica.Resolve(context.Background(), resolveAction(f, action.MatchID, true))

	require.NoError(t, err)
	require.Len(t, f.client.CreateCalls(), 1)
	require.Empty(t, replicaClient.CreateCalls())
	endCalls := replicaClient.EndCalls()
	require.Len(t, endCalls, 1)
	require.Equal(t, "prediction-1", endCalls[0].ID)
	require.Equal(t, "RESOLVED", endCalls[0].Status)
}

func TestCreateReturnsReservationStoreFailureWithoutCallingTwitch(t *testing.T) {
	f := newFixture(t)
	f.store.reserveErr = errors.New("redis unavailable")

	err := f.predictions.Create(context.Background(), createAction(f, 96))

	require.ErrorIs(t, err, f.store.reserveErr)
	require.Empty(t, f.client.CreateCalls())
}

func TestResolveSettlesAndResolvesStoredActivePrediction(t *testing.T) {
	for _, tt := range []struct {
		name           string
		win            bool
		winningOutcome string
		mmrDelta       int
	}{
		{name: "win", win: true, winningOutcome: "yes-outcome", mmrDelta: 25},
		{name: "loss", win: false, winningOutcome: "no-outcome", mmrDelta: -25},
	} {
		t.Run(tt.name, func(t *testing.T) {
			f := newFixture(t)
			f.settings.settings.Enabled = false
			f.settings.settings.MmrDelta = 25
			f.store.Store(predictionKey(f.channelID, 97), storedPrediction{
				PredictionID: "prediction-97",
				YesOutcomeID: "yes-outcome",
				NoOutcomeID:  "no-outcome",
			})
			f.client.getResponse = predictionResponse("prediction-97", "ACTIVE")

			_, err := f.predictions.Resolve(context.Background(), resolveAction(f, 97, tt.win))

			require.NoError(t, err)
			require.Equal(t, 1, f.settings.Calls())
			require.Equal(t, []dotarepository.ApplyMatchResultInput{{
				ChannelID: f.channelID,
				MatchID:   97,
				Won:       tt.win,
				MmrDelta:  tt.mmrDelta,
			}}, f.settings.ApplyCalls())
			require.Equal(t, []*helix.PredictionsParams{{
				BroadcasterID: f.broadcaster,
				ID:            "prediction-97",
			}}, f.client.GetCalls())
			require.Equal(t, []*helix.EndPredictionParams{{
				BroadcasterID:    f.broadcaster,
				ID:               "prediction-97",
				Status:           "RESOLVED",
				WinningOutcomeID: tt.winningOutcome,
			}}, f.client.EndCalls())
			_, exists := f.store.Record(predictionKey(f.channelID, 97))
			require.False(t, exists)
		})
	}
}

func TestResolveCompletesStoredPredictionThatIsNoLongerActive(t *testing.T) {
	f := newFixture(t)
	key := predictionKey(f.channelID, 98)
	f.store.Store(key, storedPrediction{PredictionID: "prediction-98", YesOutcomeID: "yes", NoOutcomeID: "no"})
	f.client.getResponse = predictionResponse("prediction-98", "RESOLVED")

	_, err := f.predictions.Resolve(context.Background(), resolveAction(f, 98, true))

	require.NoError(t, err)
	require.Empty(t, f.client.EndCalls())
	require.Equal(t, []string{key}, f.store.completeTerminalCalls)
	_, exists := f.store.Record(key)
	require.False(t, exists)
}

func TestResolveRetainsStoredPredictionAfterTransientEndFailure(t *testing.T) {
	f := newFixture(t)
	key := predictionKey(f.channelID, 99)
	f.store.Store(key, storedPrediction{PredictionID: "prediction-99", YesOutcomeID: "yes", NoOutcomeID: "no"})
	f.client.getResponse = predictionResponse("prediction-99", "LOCKED")
	f.client.endErr = errors.New("twitch temporarily unavailable")

	_, err := f.predictions.Resolve(context.Background(), resolveAction(f, 99, true))

	require.ErrorIs(t, err, f.client.endErr)
	require.Empty(t, f.store.completeTerminalCalls)
	require.Equal(t, []string{key}, f.store.releaseTerminalCalls)
	_, exists := f.store.Record(key)
	require.True(t, exists)
}

func TestResolveRetainsStoredPredictionAfterTemporaryStorageFailure(t *testing.T) {
	f := newFixture(t)
	key := predictionKey(f.channelID, 100)
	f.store.Store(key, storedPrediction{PredictionID: "prediction-100", YesOutcomeID: "yes", NoOutcomeID: "no"})
	f.store.completeTerminalErr = errors.New("redis temporarily unavailable")
	f.client.getResponse = predictionResponse("prediction-100", "RESOLVED")

	_, err := f.predictions.Resolve(context.Background(), resolveAction(f, 100, true))

	require.ErrorIs(t, err, f.store.completeTerminalErr)
	_, exists := f.store.Record(key)
	require.True(t, exists)
}

func TestCancelCancelsOnlyStoredActivePrediction(t *testing.T) {
	f := newFixture(t)
	key := predictionKey(f.channelID, 101)
	f.store.Store(key, storedPrediction{PredictionID: "dota-prediction", YesOutcomeID: "yes", NoOutcomeID: "no"})
	f.client.getResponse = predictionResponse("dota-prediction", "ACTIVE")

	err := f.predictions.Cancel(context.Background(), cancelAction(f, 101))

	require.NoError(t, err)
	require.Equal(t, []*helix.PredictionsParams{{
		BroadcasterID: f.broadcaster,
		ID:            "dota-prediction",
	}}, f.client.GetCalls())
	require.Equal(t, []*helix.EndPredictionParams{{
		BroadcasterID: f.broadcaster,
		ID:            "dota-prediction",
		Status:        "CANCELED",
	}}, f.client.EndCalls())
	_, exists := f.store.Record(key)
	require.False(t, exists)
}

func TestCancelRetainsStoredPredictionAfterTransientCancelFailure(t *testing.T) {
	f := newFixture(t)
	key := predictionKey(f.channelID, 102)
	f.store.Store(key, storedPrediction{PredictionID: "prediction-102", YesOutcomeID: "yes", NoOutcomeID: "no"})
	f.client.getResponse = predictionResponse("prediction-102", "ACTIVE")
	f.client.endErr = errors.New("twitch temporarily unavailable")

	err := f.predictions.Cancel(context.Background(), cancelAction(f, 102))

	require.ErrorIs(t, err, f.client.endErr)
	require.Empty(t, f.store.completeTerminalCalls)
	require.Equal(t, []string{key}, f.store.releaseTerminalCalls)
	_, exists := f.store.Record(key)
	require.True(t, exists)
}
