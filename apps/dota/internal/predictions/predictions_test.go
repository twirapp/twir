package predictions

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/nicklaw5/helix/v2"
	"github.com/stretchr/testify/require"
	busdota "github.com/twirapp/twir/libs/bus-core/dota"
	channelsmodel "github.com/twirapp/twir/libs/repositories/channels/model"
	dotamodel "github.com/twirapp/twir/libs/repositories/dota/model"
	"go.uber.org/fx"
)

type fakeSettingsRepository struct {
	mu       sync.Mutex
	settings dotamodel.ChannelDotaSettings
	err      error
	calls    int
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

type fakePredictionStore struct {
	mu sync.Mutex

	records      map[string]storedPrediction
	reservations map[string]string

	reserveErr error
	commitErr  error
	getErr     error
	releaseErr error
	deleteErr  error

	reserveCalls []string
	commitCalls  []string
	releaseCalls []string
	deleteCalls  []string
}

func newFakePredictionStore() *fakePredictionStore {
	return &fakePredictionStore{
		records:      make(map[string]storedPrediction),
		reservations: make(map[string]string),
	}
}

func (s *fakePredictionStore) Reserve(
	_ context.Context,
	key string,
	token string,
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
	s.reservations[key] = token
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
	if s.reservations[key] != token {
		return errPredictionReservationLost
	}
	delete(s.reservations, key)
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

func (s *fakePredictionStore) Release(_ context.Context, key string, token string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.releaseCalls = append(s.releaseCalls, key)
	if s.releaseErr != nil {
		return s.releaseErr
	}
	if s.reservations[key] == token {
		delete(s.reservations, key)
	}
	return nil
}

func (s *fakePredictionStore) Delete(_ context.Context, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.deleteCalls = append(s.deleteCalls, key)
	if s.deleteErr != nil {
		return s.deleteErr
	}
	delete(s.records, key)
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

type fakePredictionClient struct {
	mu sync.Mutex

	createResponse *helix.PredictionsResponse
	createErr      error
	getResponse    *helix.PredictionsResponse
	getErr         error
	endResponse    *helix.PredictionsResponse
	endErr         error

	createCalls []*helix.CreatePredictionParams
	getCalls    []*helix.PredictionsParams
	endCalls    []*helix.EndPredictionParams

	createStarted chan struct{}
	createRelease chan struct{}
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
	defer c.mu.Unlock()
	c.getCalls = append(c.getCalls, params)
	return c.getResponse, c.getErr
}

func (c *fakePredictionClient) EndPrediction(
	params *helix.EndPredictionParams,
) (*helix.PredictionsResponse, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.endCalls = append(c.endCalls, params)
	return c.endResponse, c.endErr
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

type fakeClientFactory struct {
	mu      sync.Mutex
	client  predictionClient
	err     error
	userIDs []uuid.UUID
}

func (f *fakeClientFactory) New(_ context.Context, userID uuid.UUID) (predictionClient, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.userIDs = append(f.userIDs, userID)
	if f.err != nil {
		return nil, f.err
	}
	return f.client, nil
}

type fakeSubscription struct {
	groups       []string
	unsubscribes int
	err          error
}

func (s *fakeSubscription) Subscribe(group string) error {
	s.groups = append(s.groups, group)
	return s.err
}

func (s *fakeSubscription) Unsubscribe() {
	s.unsubscribes++
}

type notifyingSubscription struct {
	unsubscribed chan struct{}
}

func (s *notifyingSubscription) Subscribe(_ string) error {
	return nil
}

func (s *notifyingSubscription) Unsubscribe() {
	close(s.unsubscribed)
}

type fakeLifecycle struct {
	hooks []fx.Hook
}

func (l *fakeLifecycle) Append(hook fx.Hook) {
	l.hooks = append(l.hooks, hook)
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
			nil,
			nil,
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

func startMessage(f *fixture, matchID int64) busdota.MatchStartedMessage {
	return busdota.MatchStartedMessage{
		ChannelID: f.channelID.String(),
		MatchID:   matchID,
		TeamKnown: true,
	}
}

func endMessage(f *fixture, matchID int64, win bool) busdota.MatchEndedMessage {
	return busdota.MatchEndedMessage{
		ChannelID: f.channelID.String(),
		MatchID:   matchID,
		Win:       win,
	}
}

func abandonedMessage(f *fixture, matchID int64) busdota.MatchAbandonedMessage {
	return busdota.MatchAbandonedMessage{
		ChannelID: f.channelID.String(),
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

func TestMatchStartedCreatesAndStoresPrediction(t *testing.T) {
	f := newFixture(t)

	_, err := f.predictions.handleMatchStarted(context.Background(), startMessage(f, 91))

	require.NoError(t, err)
	createCalls := f.client.CreateCalls()
	require.Len(t, createCalls, 1)
	require.Equal(t, &helix.CreatePredictionParams{
		BroadcasterID:    f.broadcaster,
		Title:            "Will the streamer win?",
		PredictionWindow: 300,
		Outcomes: []helix.PredictionChoiceParam{
			{Title: "Yes"},
			{Title: "No"},
		},
	}, createCalls[0])

	record, ok := f.store.Record(predictionKey(f.channelID, 91))
	require.True(t, ok)
	require.Equal(t, storedPrediction{
		PredictionID: "prediction-1",
		YesOutcomeID: "yes-outcome",
		NoOutcomeID:  "no-outcome",
	}, record)
	require.Equal(t, []uuid.UUID{f.twitchUser}, f.clients.userIDs)
}

func TestMatchStartedSkipsIneligibleInputs(t *testing.T) {
	tests := []struct {
		name      string
		configure func(*fixture, *busdota.MatchStartedMessage)
	}{
		{
			name: "disabled dota module",
			configure: func(f *fixture, _ *busdota.MatchStartedMessage) {
				f.settings.settings.Enabled = false
			},
		},
		{
			name: "disabled prediction settings",
			configure: func(f *fixture, _ *busdota.MatchStartedMessage) {
				f.settings.settings.PredictionSettings.Enabled = false
			},
		},
		{
			name: "missing match id",
			configure: func(_ *fixture, msg *busdota.MatchStartedMessage) {
				msg.MatchID = 0
			},
		},
		{
			name: "unknown team",
			configure: func(_ *fixture, msg *busdota.MatchStartedMessage) {
				msg.TeamKnown = false
			},
		},
		{
			name: "blank title",
			configure: func(f *fixture, _ *busdota.MatchStartedMessage) {
				f.settings.settings.PredictionSettings.TitleTemplate = " \t"
			},
		},
		{
			name: "zero window",
			configure: func(f *fixture, _ *busdota.MatchStartedMessage) {
				f.settings.settings.PredictionSettings.WindowSeconds = 0
			},
		},
		{
			name: "window above helix limit",
			configure: func(f *fixture, _ *busdota.MatchStartedMessage) {
				f.settings.settings.PredictionSettings.WindowSeconds = 1_801
			},
		},
		{
			name: "disconnected channel",
			configure: func(f *fixture, _ *busdota.MatchStartedMessage) {
				f.channels.channel = channelsmodel.Channel{}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := newFixture(t)
			message := startMessage(f, 92)
			tt.configure(f, &message)

			_, err := f.predictions.handleMatchStarted(context.Background(), message)

			require.NoError(t, err)
			require.Empty(t, f.client.CreateCalls())
			require.Empty(t, f.store.reserveCalls)
		})
	}
}

func TestMatchStartedEnforcesTwitchPredictionWindowLimits(t *testing.T) {
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

			_, err := f.predictions.handleMatchStarted(context.Background(), startMessage(f, 921))

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

func TestMatchStartedEnforcesTwitchPredictionTitleRuneLimits(t *testing.T) {
	for _, tt := range []struct {
		name    string
		title   string
		creates bool
	}{
		{name: "45 ASCII characters accepted", title: strings.Repeat("a", 45), creates: true},
		{name: "46 ASCII characters rejected", title: strings.Repeat("a", 46), creates: false},
		{name: "45 multibyte characters accepted", title: strings.Repeat("界", 45), creates: true},
		{name: "46 multibyte characters rejected", title: strings.Repeat("界", 46), creates: false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			f := newFixture(t)
			f.settings.settings.PredictionSettings.TitleTemplate = tt.title

			_, err := f.predictions.handleMatchStarted(context.Background(), startMessage(f, 922))

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

func TestMatchStartedDeduplicatesConcurrentAndReplayedEvents(t *testing.T) {
	f := newFixture(t)
	message := startMessage(f, 93)

	var wg sync.WaitGroup
	errs := make(chan error, 2)
	for range 2 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := f.predictions.handleMatchStarted(context.Background(), message)
			errs <- err
		}()
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		require.NoError(t, err)
	}

	_, err := f.predictions.handleMatchStarted(context.Background(), message)
	require.NoError(t, err)
	require.Len(t, f.client.CreateCalls(), 1)
	require.Len(t, f.store.commitCalls, 1)
}

func TestMatchStartedTreatsAlreadyActivePredictionAsNoopAndReleasesReservation(t *testing.T) {
	f := newFixture(t)
	f.client.createResponse = &helix.PredictionsResponse{
		ResponseCommon: helix.ResponseCommon{ErrorMessage: "Prediction is ALREADY ACTIVE"},
	}

	_, err := f.predictions.handleMatchStarted(context.Background(), startMessage(f, 94))

	require.NoError(t, err)
	key := predictionKey(f.channelID, 94)
	require.Equal(t, []string{key}, f.store.releaseCalls)
	require.False(t, f.store.HasReservation(key))
	_, exists := f.store.Record(key)
	require.False(t, exists)
}

func TestMatchStartedTreatsAlreadyActiveCreateErrorAsNoopAndReleasesReservation(t *testing.T) {
	f := newFixture(t)
	f.client.createResponse = nil
	f.client.createErr = errors.New("prediction is already active")

	_, err := f.predictions.handleMatchStarted(context.Background(), startMessage(f, 941))

	require.NoError(t, err)
	key := predictionKey(f.channelID, 941)
	require.Equal(t, []string{key}, f.store.releaseCalls)
	require.False(t, f.store.HasReservation(key))
}

func TestMatchStartedReleasesReservationAfterTransientCreateFailure(t *testing.T) {
	f := newFixture(t)
	f.client.createErr = errors.New("twitch temporarily unavailable")
	message := startMessage(f, 95)

	_, err := f.predictions.handleMatchStarted(context.Background(), message)

	require.Error(t, err)
	key := predictionKey(f.channelID, 95)
	require.False(t, f.store.HasReservation(key))
	require.Equal(t, []string{key}, f.store.releaseCalls)

	f.client.createErr = nil
	_, err = f.predictions.handleMatchStarted(context.Background(), message)
	require.NoError(t, err)
	_, exists := f.store.Record(key)
	require.True(t, exists)
}

func TestMatchStartedRetainsKnownPredictionWhenCommitFailsAndRecoversOnReplay(t *testing.T) {
	f := newFixture(t)
	message := startMessage(f, 951)
	key := predictionKey(f.channelID, message.MatchID)
	f.store.commitErr = errors.New("redis temporarily unavailable")

	_, err := f.predictions.handleMatchStarted(context.Background(), message)

	require.ErrorIs(t, err, f.store.commitErr)
	require.True(t, f.store.HasReservation(key))
	require.Len(t, f.client.CreateCalls(), 1)
	require.Empty(t, f.store.releaseCalls)

	f.store.commitErr = nil
	_, err = f.predictions.handleMatchStarted(context.Background(), message)

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

func TestTerminalEventsRecoverKnownPredictionAfterCommitFailure(t *testing.T) {
	for _, tt := range []struct {
		name   string
		handle func(*fixture) error
		status string
	}{
		{
			name: "resolve",
			handle: func(f *fixture) error {
				_, err := f.predictions.handleMatchEnded(context.Background(), endMessage(f, 952, true))
				return err
			},
			status: "RESOLVED",
		},
		{
			name: "cancel",
			handle: func(f *fixture) error {
				_, err := f.predictions.handleMatchAbandoned(context.Background(), abandonedMessage(f, 952))
				return err
			},
			status: "CANCELED",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			f := newFixture(t)
			message := startMessage(f, 952)
			f.store.commitErr = errors.New("redis temporarily unavailable")

			_, err := f.predictions.handleMatchStarted(context.Background(), message)
			require.ErrorIs(t, err, f.store.commitErr)
			require.Len(t, f.client.CreateCalls(), 1)

			f.store.commitErr = nil
			f.client.getResponse = predictionResponse("prediction-1", "ACTIVE")
			require.NoError(t, tt.handle(f))
			require.Len(t, f.client.CreateCalls(), 1)
			endCalls := f.client.EndCalls()
			require.Len(t, endCalls, 1)
			require.Equal(t, tt.status, endCalls[0].Status)
			_, exists := f.store.Record(predictionKey(f.channelID, message.MatchID))
			require.False(t, exists)
		})
	}
}

func TestMatchStartedReturnsReservationStoreFailureWithoutCallingTwitch(t *testing.T) {
	f := newFixture(t)
	f.store.reserveErr = errors.New("redis unavailable")

	_, err := f.predictions.handleMatchStarted(context.Background(), startMessage(f, 96))

	require.ErrorIs(t, err, f.store.reserveErr)
	require.Empty(t, f.client.CreateCalls())
}

func TestMatchEndedResolvesStoredActivePredictionWithoutRecheckingSettings(t *testing.T) {
	for _, tt := range []struct {
		name           string
		win            bool
		winningOutcome string
	}{
		{name: "win", win: true, winningOutcome: "yes-outcome"},
		{name: "loss", win: false, winningOutcome: "no-outcome"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			f := newFixture(t)
			f.settings.settings.Enabled = false
			f.store.Store(predictionKey(f.channelID, 97), storedPrediction{
				PredictionID: "prediction-97",
				YesOutcomeID: "yes-outcome",
				NoOutcomeID:  "no-outcome",
			})
			f.client.getResponse = predictionResponse("prediction-97", "ACTIVE")

			_, err := f.predictions.handleMatchEnded(context.Background(), endMessage(f, 97, tt.win))

			require.NoError(t, err)
			require.Zero(t, f.settings.Calls())
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

func TestMatchEndedDeletesStoredPredictionThatIsNoLongerActive(t *testing.T) {
	f := newFixture(t)
	key := predictionKey(f.channelID, 98)
	f.store.Store(key, storedPrediction{PredictionID: "prediction-98", YesOutcomeID: "yes", NoOutcomeID: "no"})
	f.client.getResponse = predictionResponse("prediction-98", "RESOLVED")

	_, err := f.predictions.handleMatchEnded(context.Background(), endMessage(f, 98, true))

	require.NoError(t, err)
	require.Empty(t, f.client.EndCalls())
	require.Equal(t, []string{key}, f.store.deleteCalls)
	_, exists := f.store.Record(key)
	require.False(t, exists)
}

func TestMatchEndedRetainsStoredPredictionAfterTransientEndFailure(t *testing.T) {
	f := newFixture(t)
	key := predictionKey(f.channelID, 99)
	f.store.Store(key, storedPrediction{PredictionID: "prediction-99", YesOutcomeID: "yes", NoOutcomeID: "no"})
	f.client.getResponse = predictionResponse("prediction-99", "LOCKED")
	f.client.endErr = errors.New("twitch temporarily unavailable")

	_, err := f.predictions.handleMatchEnded(context.Background(), endMessage(f, 99, true))

	require.ErrorIs(t, err, f.client.endErr)
	require.Empty(t, f.store.deleteCalls)
	_, exists := f.store.Record(key)
	require.True(t, exists)
}

func TestMatchEndedRetainsStoredPredictionAfterTemporaryStorageFailure(t *testing.T) {
	f := newFixture(t)
	key := predictionKey(f.channelID, 100)
	f.store.Store(key, storedPrediction{PredictionID: "prediction-100", YesOutcomeID: "yes", NoOutcomeID: "no"})
	f.store.deleteErr = errors.New("redis temporarily unavailable")
	f.client.getResponse = predictionResponse("prediction-100", "RESOLVED")

	_, err := f.predictions.handleMatchEnded(context.Background(), endMessage(f, 100, true))

	require.ErrorIs(t, err, f.store.deleteErr)
	_, exists := f.store.Record(key)
	require.True(t, exists)
}

func TestMatchAbandonedCancelsOnlyStoredActivePrediction(t *testing.T) {
	f := newFixture(t)
	key := predictionKey(f.channelID, 101)
	f.store.Store(key, storedPrediction{PredictionID: "dota-prediction", YesOutcomeID: "yes", NoOutcomeID: "no"})
	f.client.getResponse = predictionResponse("dota-prediction", "ACTIVE")

	_, err := f.predictions.handleMatchAbandoned(context.Background(), abandonedMessage(f, 101))

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

func TestMatchAbandonedRetainsStoredPredictionAfterTransientCancelFailure(t *testing.T) {
	f := newFixture(t)
	key := predictionKey(f.channelID, 102)
	f.store.Store(key, storedPrediction{PredictionID: "prediction-102", YesOutcomeID: "yes", NoOutcomeID: "no"})
	f.client.getResponse = predictionResponse("prediction-102", "ACTIVE")
	f.client.endErr = errors.New("twitch temporarily unavailable")

	_, err := f.predictions.handleMatchAbandoned(context.Background(), abandonedMessage(f, 102))

	require.ErrorIs(t, err, f.client.endErr)
	require.Empty(t, f.store.deleteCalls)
	_, exists := f.store.Record(key)
	require.True(t, exists)
}

func TestLifecycleSubscribesWithDedicatedGroupAndWaitsForHandlers(t *testing.T) {
	f := newFixture(t)
	lifecycle := &fakeLifecycle{}
	subscriptions := []*fakeSubscription{{}, {}, {}}
	service := newPredictions(
		f.settings,
		f.channels,
		f.clients,
		f.store,
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		[]subscription{subscriptions[0], subscriptions[1], subscriptions[2]},
		lifecycle,
	)
	require.NotNil(t, service)
	require.Len(t, lifecycle.hooks, 1)

	require.NoError(t, lifecycle.hooks[0].OnStart(context.Background()))
	for _, sub := range subscriptions {
		require.Equal(t, []string{"dota-predictions"}, sub.groups)
	}
	require.NoError(t, lifecycle.hooks[0].OnStop(context.Background()))
	for _, sub := range subscriptions {
		require.Equal(t, 1, sub.unsubscribes)
	}
}

func TestLifecycleStopWaitsForInFlightHandler(t *testing.T) {
	f := newFixture(t)
	lifecycle := &fakeLifecycle{}
	sub := &notifyingSubscription{unsubscribed: make(chan struct{})}
	f.client.createStarted = make(chan struct{})
	f.client.createRelease = make(chan struct{})
	service := newPredictions(
		f.settings,
		f.channels,
		f.clients,
		f.store,
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		[]subscription{sub},
		lifecycle,
	)

	require.NoError(t, lifecycle.hooks[0].OnStart(context.Background()))
	handlerResult := make(chan error, 1)
	go func() {
		_, err := service.handleMatchStarted(context.Background(), startMessage(f, 103))
		handlerResult <- err
	}()
	<-f.client.createStarted

	stopResult := make(chan error, 1)
	go func() {
		stopResult <- lifecycle.hooks[0].OnStop(context.Background())
	}()
	<-sub.unsubscribed
	select {
	case err := <-stopResult:
		t.Fatalf("OnStop returned while the handler was blocked: %v", err)
	default:
	}

	close(f.client.createRelease)
	require.NoError(t, <-handlerResult)
	require.NoError(t, <-stopResult)
}
