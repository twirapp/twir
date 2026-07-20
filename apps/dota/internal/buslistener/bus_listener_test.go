package buslistener

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strconv"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/twirapp/twir/apps/dota/internal/match"
	"github.com/twirapp/twir/apps/dota/internal/stats"
	buscore "github.com/twirapp/twir/libs/bus-core"
	busdota "github.com/twirapp/twir/libs/bus-core/dota"
	dotarepository "github.com/twirapp/twir/libs/repositories/dota"
	"github.com/twirapp/twir/libs/repositories/dota/model"
	"go.uber.org/fx"
)

var errNotImplemented = errors.New("not implemented")

type testLogRecord struct {
	message string
	attrs   []slog.Attr
}

type testLogHandler struct {
	mu       sync.Mutex
	messages []string
	records  []testLogRecord
}

func (h *testLogHandler) Enabled(context.Context, slog.Level) bool {
	return true
}

func (h *testLogHandler) Handle(_ context.Context, record slog.Record) error {
	attrs := make([]slog.Attr, 0, record.NumAttrs())
	record.Attrs(func(attr slog.Attr) bool {
		attrs = append(attrs, attr)
		return true
	})

	h.mu.Lock()
	defer h.mu.Unlock()

	h.messages = append(h.messages, record.Message)
	h.records = append(h.records, testLogRecord{message: record.Message, attrs: attrs})
	return nil
}

func (h *testLogHandler) WithAttrs([]slog.Attr) slog.Handler {
	return h
}

func (h *testLogHandler) WithGroup(string) slog.Handler {
	return h
}

func (h *testLogHandler) Messages() []string {
	h.mu.Lock()
	defer h.mu.Unlock()

	return append([]string(nil), h.messages...)
}

func (h *testLogHandler) Records() []testLogRecord {
	h.mu.Lock()
	defer h.mu.Unlock()

	records := make([]testLogRecord, len(h.records))
	for i, record := range h.records {
		records[i] = testLogRecord{
			message: record.message,
			attrs:   append([]slog.Attr(nil), record.attrs...),
		}
	}

	return records
}

func (r testLogRecord) Attr(key string) (slog.Attr, bool) {
	for _, attr := range r.attrs {
		if attr.Key == key {
			return attr, true
		}
	}

	return slog.Attr{}, false
}

func logAttrText(attr slog.Attr) string {
	if attr.Value.Kind() == slog.KindAny {
		return fmt.Sprint(attr.Value.Any())
	}

	return attr.Value.String()
}

type fakeRepository struct {
	mu         sync.Mutex
	settings   model.ChannelDotaSettings
	err        error
	channelIDs []uuid.UUID
}

var _ dotarepository.Repository = (*fakeRepository)(nil)

func (f *fakeRepository) GetByChannelID(
	_ context.Context,
	channelID uuid.UUID,
) (model.ChannelDotaSettings, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.channelIDs = append(f.channelIDs, channelID)
	if f.err != nil {
		return model.Nil, f.err
	}

	return f.settings, nil
}

func (f *fakeRepository) GetByGsiToken(
	_ context.Context,
	_ string,
) (model.ChannelDotaSettings, error) {
	return model.Nil, errNotImplemented
}

func (f *fakeRepository) Create(
	_ context.Context,
	_ dotarepository.CreateInput,
) (model.ChannelDotaSettings, error) {
	return model.Nil, errNotImplemented
}

func (f *fakeRepository) Update(
	_ context.Context,
	_ uuid.UUID,
	_ dotarepository.UpdateInput,
) (model.ChannelDotaSettings, error) {
	return model.Nil, errNotImplemented
}

func (f *fakeRepository) UpdateMatchResult(
	_ context.Context,
	_ uuid.UUID,
	_ bool,
	_ int,
) (model.ChannelDotaSettings, error) {
	return model.Nil, errNotImplemented
}

func (f *fakeRepository) ResetSession(
	_ context.Context,
	_ uuid.UUID,
) (model.ChannelDotaSettings, error) {
	return model.Nil, errNotImplemented
}

func (f *fakeRepository) RegenerateGsiToken(
	_ context.Context,
	_ uuid.UUID,
) (model.ChannelDotaSettings, error) {
	return model.Nil, errNotImplemented
}

type fakeStateProvider struct {
	mu         sync.Mutex
	snapshot   match.Snapshot
	err        error
	channelIDs []uuid.UUID
}

func (f *fakeStateProvider) GetSnapshot(
	_ context.Context,
	channelID uuid.UUID,
) (match.Snapshot, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.channelIDs = append(f.channelIDs, channelID)
	return f.snapshot, f.err
}

type notablePlayersCall struct {
	matchID           int64
	streamerAccountID string
}

type fakeStatsProvider struct {
	mu                  sync.Mutex
	winProbability      float64
	winProbabilityErr   error
	winProbabilityStart chan struct{}
	winProbabilityDone  <-chan struct{}
	notablePlayers      []string
	notablePlayersErr   error
	lastGame            *stats.LastGame
	lastGameErr         error
	winProbabilityCalls []int64
	notablePlayersCalls []notablePlayersCall
	lastGameCalls       []int64
}

var _ StatsProvider = (*fakeStatsProvider)(nil)

func (f *fakeStatsProvider) WinProbability(ctx context.Context, matchID int64) (float64, error) {
	f.mu.Lock()
	f.winProbabilityCalls = append(f.winProbabilityCalls, matchID)
	probability := f.winProbability
	err := f.winProbabilityErr
	start := f.winProbabilityStart
	done := f.winProbabilityDone
	f.mu.Unlock()

	if start != nil {
		start <- struct{}{}
	}
	if done != nil {
		select {
		case <-done:
		case <-ctx.Done():
			return 0, ctx.Err()
		}
	}

	return probability, err
}

func (f *fakeStatsProvider) NotablePlayers(
	_ context.Context,
	matchID int64,
	streamerAccountID string,
) ([]string, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.notablePlayersCalls = append(
		f.notablePlayersCalls,
		notablePlayersCall{matchID: matchID, streamerAccountID: streamerAccountID},
	)
	return f.notablePlayers, f.notablePlayersErr
}

func (f *fakeStatsProvider) LastGame(_ context.Context, accountID int64) (*stats.LastGame, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.lastGameCalls = append(f.lastGameCalls, accountID)
	return f.lastGame, f.lastGameErr
}

func (f *fakeStatsProvider) callCounts() (winProbability, notablePlayers, lastGame int) {
	f.mu.Lock()
	defer f.mu.Unlock()

	return len(f.winProbabilityCalls), len(f.notablePlayersCalls), len(f.lastGameCalls)
}

type fakeLifecycle struct {
	hooks []fx.Hook
}

func (f *fakeLifecycle) Append(hook fx.Hook) {
	f.hooks = append(f.hooks, hook)
}

type fakeGetDataQueue struct {
	mu               sync.Mutex
	group            string
	callback         buscore.QueueSubscribeCallback[busdota.GetDataRequest, busdota.GetDataResponse]
	subscribeErr     error
	unsubscribes     int
	unsubscribed     chan struct{}
	allowUnsubscribe <-chan struct{}
	unsubscribeOnce  sync.Once
}

func (f *fakeGetDataQueue) SubscribeGroup(
	group string,
	callback buscore.QueueSubscribeCallback[busdota.GetDataRequest, busdota.GetDataResponse],
) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.group = group
	f.callback = callback
	return f.subscribeErr
}

func (f *fakeGetDataQueue) Unsubscribe() {
	f.mu.Lock()
	f.unsubscribes++
	unsubscribed := f.unsubscribed
	allowUnsubscribe := f.allowUnsubscribe
	f.mu.Unlock()

	if unsubscribed != nil {
		f.unsubscribeOnce.Do(func() { close(unsubscribed) })
	}
	if allowUnsubscribe != nil {
		<-allowUnsubscribe
	}
}

func (f *fakeGetDataQueue) Group() string {
	f.mu.Lock()
	defer f.mu.Unlock()

	return f.group
}

func (f *fakeGetDataQueue) Callback() buscore.QueueSubscribeCallback[busdota.GetDataRequest, busdota.GetDataResponse] {
	f.mu.Lock()
	defer f.mu.Unlock()

	return f.callback
}

func (f *fakeGetDataQueue) UnsubscribeCount() int {
	f.mu.Lock()
	defer f.mu.Unlock()

	return f.unsubscribes
}

func TestGetDataRejectsInvalidChannelID(t *testing.T) {
	repo := &fakeRepository{}
	state := &fakeStateProvider{}
	statsProvider := &fakeStatsProvider{}

	response, err := newTestListener(repo, state, statsProvider).GetData(
		context.Background(),
		busdota.GetDataRequest{ChannelID: "not-a-uuid"},
	)

	require.Error(t, err)
	require.Zero(t, response)
	require.Empty(t, repo.channelIDs)
	require.Empty(t, state.channelIDs)
	requireNoStatsCalls(t, statsProvider)
}

func TestGetDataReturnsDisabledUnlinkedWhenSettingsNotFound(t *testing.T) {
	channelID := uuid.New()
	repo := &fakeRepository{err: dotarepository.ErrNotFound}
	state := &fakeStateProvider{}
	statsProvider := &fakeStatsProvider{}

	response, err := newTestListener(repo, state, statsProvider).GetData(
		context.Background(),
		busdota.GetDataRequest{ChannelID: channelID.String()},
	)

	require.NoError(t, err)
	require.Equal(t, busdota.GetDataResponse{Enabled: false, Linked: false}, response)
	require.Equal(t, []uuid.UUID{channelID}, repo.channelIDs)
	require.Empty(t, state.channelIDs)
	requireNoStatsCalls(t, statsProvider)
}

func TestGetDataMapsActiveLinkedMatch(t *testing.T) {
	channelID := uuid.New()
	settingsAccountID := "111"
	repo := &fakeRepository{settings: model.ChannelDotaSettings{
		Enabled:        true,
		SteamAccountID: &settingsAccountID,
		Mmr:            4200,
		SessionWins:    8,
		SessionLosses:  3,
	}}
	state := &fakeStateProvider{snapshot: match.Snapshot{
		InGame:         true,
		HeroName:       "Pudge",
		MatchID:        42,
		IsRadiant:      true,
		SteamAccountID: "222",
		RadiantScore:   21,
		DireScore:      18,
		GameTime:       1134,
	}}
	statsProvider := &fakeStatsProvider{
		winProbability: 0.73,
		notablePlayers: []string{"Player One", "Player Two"},
	}

	response, err := newTestListener(repo, state, statsProvider).GetData(
		context.Background(),
		busdota.GetDataRequest{ChannelID: channelID.String()},
	)

	require.NoError(t, err)
	require.Equal(t, busdota.GetDataResponse{
		Enabled:        true,
		Linked:         true,
		InGame:         true,
		Mmr:            4200,
		SessionWins:    8,
		SessionLosses:  3,
		HeroName:       "Pudge",
		MatchID:        42,
		TeamIsRadiant:  true,
		RadiantScore:   21,
		DireScore:      18,
		GameTime:       1134,
		WinProbability: 0.73,
		NotablePlayers: []string{"Player One", "Player Two"},
	}, response)
	require.Equal(t, []int64{42}, statsProvider.winProbabilityCalls)
	require.Equal(
		t,
		[]notablePlayersCall{{matchID: 42, streamerAccountID: "222"}},
		statsProvider.notablePlayersCalls,
	)
	require.Empty(t, statsProvider.lastGameCalls)
}

func TestGetDataHandlesConcurrentRequests(t *testing.T) {
	const requests = 16

	channelID := uuid.New()
	accountID := "12345"
	repo := &fakeRepository{settings: model.ChannelDotaSettings{
		Enabled:        true,
		SteamAccountID: &accountID,
	}}
	state := &fakeStateProvider{snapshot: match.Snapshot{
		InGame:         true,
		MatchID:        42,
		SteamAccountID: accountID,
	}}
	statsProvider := &fakeStatsProvider{}
	listener := newTestListener(repo, state, statsProvider)

	errs := make(chan error, requests)
	var requestsGroup sync.WaitGroup
	for range requests {
		requestsGroup.Add(1)
		go func() {
			defer requestsGroup.Done()

			_, err := listener.GetData(
				context.Background(),
				busdota.GetDataRequest{ChannelID: channelID.String()},
			)
			errs <- err
		}()
	}

	requestsGroup.Wait()
	close(errs)

	for err := range errs {
		require.NoError(t, err)
	}
	winProbabilityCalls, notablePlayersCalls, lastGameCalls := statsProvider.callCounts()
	require.Equal(t, requests, winProbabilityCalls)
	require.Equal(t, requests, notablePlayersCalls)
	require.Zero(t, lastGameCalls)
}

func TestGetDataUsesSettingsAccountForNotablePlayersWhenSnapshotAccountIsEmpty(t *testing.T) {
	channelID := uuid.New()
	settingsAccountID := "111"
	repo := &fakeRepository{settings: model.ChannelDotaSettings{
		Enabled:        true,
		SteamAccountID: &settingsAccountID,
	}}
	state := &fakeStateProvider{snapshot: match.Snapshot{InGame: true, MatchID: 42}}
	statsProvider := &fakeStatsProvider{}

	_, err := newTestListener(repo, state, statsProvider).GetData(
		context.Background(),
		busdota.GetDataRequest{ChannelID: channelID.String()},
	)

	require.NoError(t, err)
	require.Equal(
		t,
		[]notablePlayersCall{{matchID: 42, streamerAccountID: settingsAccountID}},
		statsProvider.notablePlayersCalls,
	)
}

func TestGetDataNormalizesSteamID64FallbackForNotablePlayers(t *testing.T) {
	channelID := uuid.New()
	steamID64 := "76561197960278073"
	repo := &fakeRepository{settings: model.ChannelDotaSettings{
		Enabled:        true,
		SteamAccountID: &steamID64,
	}}
	state := &fakeStateProvider{snapshot: match.Snapshot{InGame: true, MatchID: 42}}
	statsProvider := &fakeStatsProvider{}

	_, err := newTestListener(repo, state, statsProvider).GetData(
		context.Background(),
		busdota.GetDataRequest{ChannelID: channelID.String()},
	)

	require.NoError(t, err)
	require.Equal(
		t,
		[]notablePlayersCall{{matchID: 42, streamerAccountID: "12345"}},
		statsProvider.notablePlayersCalls,
	)
}

func TestDotaAccountIDSanitizesInvalidInput(t *testing.T) {
	invalidAccountID := "invalid-steam-id"

	_, err := dotaAccountID(invalidAccountID)

	require.EqualError(t, err, "Steam account ID must be a decimal integer")
	require.NotContains(t, err.Error(), invalidAccountID)
	var numberErr *strconv.NumError
	require.False(t, errors.As(err, &numberErr))
}

func TestGetDataKeepsActiveMatchResponseWhenStatsFail(t *testing.T) {
	channelID := uuid.New()
	accountID := "111"
	logHandler := &testLogHandler{}
	winProbabilityErr := errors.New("win probability unavailable")
	notablePlayersErr := errors.New("notable players unavailable")
	repo := &fakeRepository{settings: model.ChannelDotaSettings{
		Enabled:        true,
		SteamAccountID: &accountID,
	}}
	state := &fakeStateProvider{snapshot: match.Snapshot{
		InGame:   true,
		HeroName: "Drow Ranger",
		MatchID:  88,
	}}
	statsProvider := &fakeStatsProvider{
		winProbability:    0.6,
		winProbabilityErr: winProbabilityErr,
		notablePlayers:    []string{"ignored"},
		notablePlayersErr: notablePlayersErr,
	}

	response, err := newTestListenerWithLogger(repo, state, statsProvider, slog.New(logHandler)).GetData(
		context.Background(),
		busdota.GetDataRequest{ChannelID: channelID.String()},
	)

	require.NoError(t, err)
	require.True(t, response.InGame)
	require.Equal(t, int64(88), response.MatchID)
	require.Zero(t, response.WinProbability)
	require.Nil(t, response.NotablePlayers)
	require.Equal(t, []int64{88}, statsProvider.winProbabilityCalls)
	require.Equal(t, []notablePlayersCall{{matchID: 88, streamerAccountID: accountID}}, statsProvider.notablePlayersCalls)
	require.Equal(t, []string{
		"dota bus listener: failed to fetch win probability",
		"dota bus listener: failed to fetch notable players",
	}, logHandler.Messages())
}

func TestGetDataLogsAndUsesEmptyNotablePlayersFallbackForInvalidSteamID(t *testing.T) {
	channelID := uuid.New()
	invalidAccountID := "invalid-steam-id"
	logHandler := &testLogHandler{}
	repo := &fakeRepository{settings: model.ChannelDotaSettings{
		Enabled:        true,
		SteamAccountID: &invalidAccountID,
	}}
	state := &fakeStateProvider{snapshot: match.Snapshot{InGame: true, MatchID: 42}}
	statsProvider := &fakeStatsProvider{}

	_, err := newTestListenerWithLogger(repo, state, statsProvider, slog.New(logHandler)).GetData(
		context.Background(),
		busdota.GetDataRequest{ChannelID: channelID.String()},
	)

	require.NoError(t, err)
	require.Equal(
		t,
		[]notablePlayersCall{{matchID: 42, streamerAccountID: ""}},
		statsProvider.notablePlayersCalls,
	)

	records := logHandler.Records()
	require.Len(t, records, 1)
	record := records[0]
	require.Equal(t, "dota bus listener: invalid Steam account ID", record.message)
	channelAttr, ok := record.Attr("channel_id")
	require.True(t, ok)
	require.Equal(t, channelID.String(), channelAttr.Value.String())
	for _, attr := range record.attrs {
		require.NotContains(t, logAttrText(attr), invalidAccountID)
	}
}

func TestGetDataMapsLastGameForIdleSteamID64Account(t *testing.T) {
	channelID := uuid.New()
	steamID64 := "76561197960278073"
	repo := &fakeRepository{settings: model.ChannelDotaSettings{
		Enabled:        true,
		SteamAccountID: &steamID64,
		Mmr:            5000,
		SessionWins:    10,
		SessionLosses:  4,
	}}
	state := &fakeStateProvider{snapshot: match.Snapshot{
		InGame:       false,
		RadiantScore: 33,
		DireScore:    29,
		GameTime:     2450,
	}}
	statsProvider := &fakeStatsProvider{lastGame: &stats.LastGame{
		HeroName:  "Crystal Maiden",
		Kills:     3,
		Deaths:    5,
		Assists:   17,
		Win:       true,
		DurationS: 2142,
	}}

	response, err := newTestListener(repo, state, statsProvider).GetData(
		context.Background(),
		busdota.GetDataRequest{ChannelID: channelID.String()},
	)

	require.NoError(t, err)
	require.Equal(t, []int64{12345}, statsProvider.lastGameCalls)
	require.Equal(t, &busdota.LastGameInfo{
		HeroName:  "Crystal Maiden",
		Kills:     3,
		Deaths:    5,
		Assists:   17,
		Win:       true,
		DurationS: 2142,
	}, response.LastGame)
	require.Empty(t, statsProvider.winProbabilityCalls)
	require.Empty(t, statsProvider.notablePlayersCalls)
}

func TestGetDataSkipsStatsForDisabledOrUnlinkedSettings(t *testing.T) {
	linkedAccountID := "12345"
	emptyAccountID := ""

	tests := []struct {
		name     string
		settings model.ChannelDotaSettings
		snapshot match.Snapshot
	}{
		{
			name: "disabled",
			settings: model.ChannelDotaSettings{
				Enabled:        false,
				SteamAccountID: &linkedAccountID,
			},
			snapshot: match.Snapshot{InGame: true, MatchID: 99},
		},
		{
			name: "unlinked",
			settings: model.ChannelDotaSettings{
				Enabled:        true,
				SteamAccountID: &emptyAccountID,
			},
			snapshot: match.Snapshot{InGame: false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			channelID := uuid.New()
			repo := &fakeRepository{settings: tt.settings}
			state := &fakeStateProvider{snapshot: tt.snapshot}
			statsProvider := &fakeStatsProvider{}

			response, err := newTestListener(repo, state, statsProvider).GetData(
				context.Background(),
				busdota.GetDataRequest{ChannelID: channelID.String()},
			)

			require.NoError(t, err)
			require.Equal(t, tt.settings.Enabled, response.Enabled)
			require.Equal(t, tt.settings.SteamAccountID != nil && *tt.settings.SteamAccountID != "", response.Linked)
			require.Equal(t, tt.snapshot.InGame, response.InGame)
			requireNoStatsCalls(t, statsProvider)
		})
	}
}

func TestGetDataKeepsIdleResponseWhenLastGameFails(t *testing.T) {
	channelID := uuid.New()
	accountID := "12345"
	logHandler := &testLogHandler{}
	lastGameErr := errors.New("last game unavailable")
	repo := &fakeRepository{settings: model.ChannelDotaSettings{
		Enabled:        true,
		SteamAccountID: &accountID,
		Mmr:            3300,
	}}
	state := &fakeStateProvider{snapshot: match.Snapshot{InGame: false, HeroName: "Lina"}}
	statsProvider := &fakeStatsProvider{lastGameErr: lastGameErr}

	response, err := newTestListenerWithLogger(repo, state, statsProvider, slog.New(logHandler)).GetData(
		context.Background(),
		busdota.GetDataRequest{ChannelID: channelID.String()},
	)

	require.NoError(t, err)
	require.Equal(t, busdota.GetDataResponse{
		Enabled:  true,
		Linked:   true,
		Mmr:      3300,
		HeroName: "Lina",
	}, response)
	require.Equal(t, []int64{12345}, statsProvider.lastGameCalls)
	require.Empty(t, statsProvider.winProbabilityCalls)
	require.Empty(t, statsProvider.notablePlayersCalls)
	require.Equal(t, []string{"dota bus listener: failed to fetch last game"}, logHandler.Messages())
}

func TestGetDataKeepsPartialResponseForInvalidStoredSteamID(t *testing.T) {
	channelID := uuid.New()
	accountID := "invalid-steam-id"
	logHandler := &testLogHandler{}
	repo := &fakeRepository{settings: model.ChannelDotaSettings{
		Enabled:        true,
		SteamAccountID: &accountID,
		Mmr:            3300,
	}}
	state := &fakeStateProvider{snapshot: match.Snapshot{
		InGame:       false,
		HeroName:     "Lina",
		RadiantScore: 17,
		DireScore:    12,
	}}
	statsProvider := &fakeStatsProvider{}

	response, err := newTestListenerWithLogger(repo, state, statsProvider, slog.New(logHandler)).GetData(
		context.Background(),
		busdota.GetDataRequest{ChannelID: channelID.String()},
	)

	require.NoError(t, err)
	require.Equal(t, busdota.GetDataResponse{
		Enabled:      true,
		Linked:       true,
		Mmr:          3300,
		HeroName:     "Lina",
		RadiantScore: 17,
		DireScore:    12,
	}, response)
	requireNoStatsCalls(t, statsProvider)

	records := logHandler.Records()
	require.Len(t, records, 1)
	record := records[0]
	require.Equal(t, "dota bus listener: invalid Steam account ID", record.message)
	channelAttr, ok := record.Attr("channel_id")
	require.True(t, ok)
	require.Equal(t, channelID.String(), channelAttr.Value.String())
	_, hasSteamAccountID := record.Attr("steam_account_id")
	require.False(t, hasSteamAccountID)
	require.NotContains(t, record.message, accountID)
	for _, attr := range record.attrs {
		require.NotContains(t, logAttrText(attr), accountID)
	}
}

func TestGetDataReturnsSnapshotError(t *testing.T) {
	channelID := uuid.New()
	accountID := "12345"
	snapshotErr := errors.New("snapshot unavailable")
	repo := &fakeRepository{settings: model.ChannelDotaSettings{
		Enabled:        true,
		SteamAccountID: &accountID,
	}}
	state := &fakeStateProvider{err: snapshotErr}
	statsProvider := &fakeStatsProvider{}

	_, err := newTestListener(repo, state, statsProvider).GetData(
		context.Background(),
		busdota.GetDataRequest{ChannelID: channelID.String()},
	)

	require.ErrorIs(t, err, snapshotErr)
	requireNoStatsCalls(t, statsProvider)
}

func TestGetDataReturnsSettingsError(t *testing.T) {
	channelID := uuid.New()
	settingsErr := errors.New("settings unavailable")
	repo := &fakeRepository{err: settingsErr}
	state := &fakeStateProvider{}
	statsProvider := &fakeStatsProvider{}

	_, err := newTestListener(repo, state, statsProvider).GetData(
		context.Background(),
		busdota.GetDataRequest{ChannelID: channelID.String()},
	)

	require.ErrorIs(t, err, settingsErr)
	require.Empty(t, state.channelIDs)
	requireNoStatsCalls(t, statsProvider)
}

func TestBusListenerOnStopWaitsForActiveGetData(t *testing.T) {
	channelID := uuid.New()
	accountID := "12345"
	requestStarted := make(chan struct{})
	releaseRequest := make(chan struct{})
	unsubscribeStarted := make(chan struct{})
	allowUnsubscribe := make(chan struct{})
	lifecycle := &fakeLifecycle{}
	queue := &fakeGetDataQueue{
		unsubscribed:     unsubscribeStarted,
		allowUnsubscribe: allowUnsubscribe,
	}
	statsProvider := &fakeStatsProvider{
		winProbabilityStart: requestStarted,
		winProbabilityDone:  releaseRequest,
	}

	newBusListener(
		&fakeStateProvider{snapshot: match.Snapshot{
			InGame:         true,
			MatchID:        42,
			SteamAccountID: accountID,
		}},
		&fakeRepository{settings: model.ChannelDotaSettings{
			Enabled:        true,
			SteamAccountID: &accountID,
		}},
		statsProvider,
		testLogger(),
		lifecycle,
		queue,
	)

	require.NoError(t, lifecycle.hooks[0].OnStart(context.Background()))
	callback := queue.Callback()
	require.NotNil(t, callback)

	requestResult := make(chan error, 1)
	go func() {
		_, err := callback(context.Background(), busdota.GetDataRequest{ChannelID: channelID.String()})
		requestResult <- err
	}()
	<-requestStarted
	t.Cleanup(func() {
		close(releaseRequest)
		require.NoError(t, <-requestResult)
	})

	stopContext, cancelStop := context.WithCancel(context.Background())
	defer cancelStop()
	stopResult := make(chan error, 1)
	go func() {
		stopResult <- lifecycle.hooks[0].OnStop(stopContext)
	}()
	<-unsubscribeStarted
	cancelStop()
	close(allowUnsubscribe)

	require.ErrorIs(t, <-stopResult, context.Canceled)
	require.Equal(t, 1, queue.UnsubscribeCount())
}

func TestBusListenerOnStopCompletesAfterActiveGetDataFinishes(t *testing.T) {
	channelID := uuid.New()
	accountID := "12345"
	requestStarted := make(chan struct{})
	releaseRequest := make(chan struct{})
	unsubscribeStarted := make(chan struct{})
	lifecycle := &fakeLifecycle{}
	queue := &fakeGetDataQueue{unsubscribed: unsubscribeStarted}
	statsProvider := &fakeStatsProvider{
		winProbabilityStart: requestStarted,
		winProbabilityDone:  releaseRequest,
	}

	newBusListener(
		&fakeStateProvider{snapshot: match.Snapshot{
			InGame:         true,
			MatchID:        42,
			SteamAccountID: accountID,
		}},
		&fakeRepository{settings: model.ChannelDotaSettings{
			Enabled:        true,
			SteamAccountID: &accountID,
		}},
		statsProvider,
		testLogger(),
		lifecycle,
		queue,
	)

	require.NoError(t, lifecycle.hooks[0].OnStart(context.Background()))
	callback := queue.Callback()
	require.NotNil(t, callback)

	requestResult := make(chan error, 1)
	go func() {
		_, err := callback(context.Background(), busdota.GetDataRequest{ChannelID: channelID.String()})
		requestResult <- err
	}()
	<-requestStarted

	stopResult := make(chan error, 1)
	go func() {
		stopResult <- lifecycle.hooks[0].OnStop(context.Background())
	}()
	<-unsubscribeStarted
	close(releaseRequest)

	require.NoError(t, <-requestResult)
	require.NoError(t, <-stopResult)
}

func TestBusListenerRejectsGetDataRequestsAfterStop(t *testing.T) {
	channelID := uuid.New()
	accountID := "12345"
	lifecycle := &fakeLifecycle{}
	queue := &fakeGetDataQueue{}
	statsProvider := &fakeStatsProvider{}

	newBusListener(
		&fakeStateProvider{snapshot: match.Snapshot{
			InGame:         true,
			MatchID:        42,
			SteamAccountID: accountID,
		}},
		&fakeRepository{settings: model.ChannelDotaSettings{
			Enabled:        true,
			SteamAccountID: &accountID,
		}},
		statsProvider,
		testLogger(),
		lifecycle,
		queue,
	)

	require.NoError(t, lifecycle.hooks[0].OnStart(context.Background()))
	require.NoError(t, lifecycle.hooks[0].OnStop(context.Background()))

	callback := queue.Callback()
	require.NotNil(t, callback)
	response, err := callback(context.Background(), busdota.GetDataRequest{ChannelID: channelID.String()})

	require.ErrorIs(t, err, context.Canceled)
	require.Zero(t, response)
	requireNoStatsCalls(t, statsProvider)
}

func TestBusListenerLifecycleSubscribesAndUnsubscribesGetData(t *testing.T) {
	lifecycle := &fakeLifecycle{}
	queue := &fakeGetDataQueue{}

	newBusListener(
		&fakeStateProvider{},
		&fakeRepository{},
		&fakeStatsProvider{},
		testLogger(),
		lifecycle,
		queue,
	)

	require.Len(t, lifecycle.hooks, 1)
	require.NoError(t, lifecycle.hooks[0].OnStart(context.Background()))
	require.Equal(t, "dota", queue.Group())
	require.NotNil(t, queue.Callback())
	require.NoError(t, lifecycle.hooks[0].OnStop(context.Background()))
	require.Equal(t, 1, queue.UnsubscribeCount())
}

func newTestListener(
	repo *fakeRepository,
	state *fakeStateProvider,
	statsProvider *fakeStatsProvider,
) *BusListener {
	return newTestListenerWithLogger(repo, state, statsProvider, testLogger())
}

func newTestListenerWithLogger(
	repo *fakeRepository,
	state *fakeStateProvider,
	statsProvider *fakeStatsProvider,
	logger *slog.Logger,
) *BusListener {
	return &BusListener{
		repository: repo,
		state:      state,
		stats:      statsProvider,
		logger:     logger,
	}
}

func requireNoStatsCalls(t *testing.T, statsProvider *fakeStatsProvider) {
	t.Helper()
	winProbabilityCalls, notablePlayersCalls, lastGameCalls := statsProvider.callCounts()
	require.Zero(t, winProbabilityCalls)
	require.Zero(t, notablePlayersCalls)
	require.Zero(t, lastGameCalls)
}

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}
