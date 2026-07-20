package buslistener

import (
	"context"
	"errors"
	"io"
	"log/slog"
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

type testLogHandler struct {
	mu       sync.Mutex
	messages []string
}

func (h *testLogHandler) Enabled(context.Context, slog.Level) bool {
	return true
}

func (h *testLogHandler) Handle(_ context.Context, record slog.Record) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.messages = append(h.messages, record.Message)
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

type fakeRepository struct {
	settings   model.ChannelDotaSettings
	err        error
	channelIDs []uuid.UUID
}

var _ dotarepository.Repository = (*fakeRepository)(nil)

func (f *fakeRepository) GetByChannelID(
	_ context.Context,
	channelID uuid.UUID,
) (model.ChannelDotaSettings, error) {
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
	snapshot   match.Snapshot
	err        error
	channelIDs []uuid.UUID
}

func (f *fakeStateProvider) GetSnapshot(
	_ context.Context,
	channelID uuid.UUID,
) (match.Snapshot, error) {
	f.channelIDs = append(f.channelIDs, channelID)
	return f.snapshot, f.err
}

type notablePlayersCall struct {
	matchID           int64
	streamerAccountID string
}

type fakeStatsProvider struct {
	winProbability      float64
	winProbabilityErr   error
	notablePlayers      []string
	notablePlayersErr   error
	lastGame            *stats.LastGame
	lastGameErr         error
	winProbabilityCalls []int64
	notablePlayersCalls []notablePlayersCall
	lastGameCalls       []int64
}

var _ StatsProvider = (*fakeStatsProvider)(nil)

func (f *fakeStatsProvider) WinProbability(_ context.Context, matchID int64) (float64, error) {
	f.winProbabilityCalls = append(f.winProbabilityCalls, matchID)
	return f.winProbability, f.winProbabilityErr
}

func (f *fakeStatsProvider) NotablePlayers(
	_ context.Context,
	matchID int64,
	streamerAccountID string,
) ([]string, error) {
	f.notablePlayersCalls = append(
		f.notablePlayersCalls,
		notablePlayersCall{matchID: matchID, streamerAccountID: streamerAccountID},
	)
	return f.notablePlayers, f.notablePlayersErr
}

func (f *fakeStatsProvider) LastGame(_ context.Context, accountID int64) (*stats.LastGame, error) {
	f.lastGameCalls = append(f.lastGameCalls, accountID)
	return f.lastGame, f.lastGameErr
}

type fakeLifecycle struct {
	hooks []fx.Hook
}

func (f *fakeLifecycle) Append(hook fx.Hook) {
	f.hooks = append(f.hooks, hook)
}

type fakeGetDataQueue struct {
	group        string
	callback     buscore.QueueSubscribeCallback[busdota.GetDataRequest, busdota.GetDataResponse]
	subscribeErr error
	unsubscribes int
}

func (f *fakeGetDataQueue) SubscribeGroup(
	group string,
	callback buscore.QueueSubscribeCallback[busdota.GetDataRequest, busdota.GetDataResponse],
) error {
	f.group = group
	f.callback = callback
	return f.subscribeErr
}

func (f *fakeGetDataQueue) Unsubscribe() {
	f.unsubscribes++
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

	response, err := newTestListener(repo, state, statsProvider).GetData(
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
	require.Equal(t, "dota", queue.group)
	require.NotNil(t, queue.callback)
	require.NoError(t, lifecycle.hooks[0].OnStop(context.Background()))
	require.Equal(t, 1, queue.unsubscribes)
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
	require.Empty(t, statsProvider.winProbabilityCalls)
	require.Empty(t, statsProvider.notablePlayersCalls)
	require.Empty(t, statsProvider.lastGameCalls)
}

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}
