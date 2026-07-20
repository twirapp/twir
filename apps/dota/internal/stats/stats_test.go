package stats

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/twirapp/kv"
	kvoptions "github.com/twirapp/kv/options"
	"github.com/twirapp/twir/libs/integrations/opendota"
	"github.com/twirapp/twir/libs/integrations/stratz"
)

type fakeValuer struct {
	value []byte
	err   error
}

func (v fakeValuer) Int() (int64, error)     { return 0, kv.ErrInvalidType }
func (v fakeValuer) String() (string, error) { return string(v.value), v.err }
func (v fakeValuer) Bytes() ([]byte, error)  { return v.value, v.err }
func (v fakeValuer) Bool() (bool, error)     { return false, kv.ErrInvalidType }
func (v fakeValuer) Float() (float64, error) { return 0, kv.ErrInvalidType }
func (v fakeValuer) Scan(dest any) error {
	if v.err != nil {
		return v.err
	}
	return json.Unmarshal(v.value, dest)
}
func (v fakeValuer) Err() error { return v.err }

type fakeKV struct {
	mu   sync.Mutex
	data map[string][]byte
}

func newFakeKV() *fakeKV {
	return &fakeKV{data: make(map[string][]byte)}
}

func (f *fakeKV) Get(_ context.Context, key string) kv.Valuer {
	f.mu.Lock()
	defer f.mu.Unlock()

	value, ok := f.data[key]
	if !ok {
		return fakeValuer{err: kv.ErrKeyNil}
	}

	return fakeValuer{value: value}
}

func (f *fakeKV) Set(_ context.Context, key string, value any, _ ...kvoptions.Option) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	bytes, ok := value.([]byte)
	if !ok {
		marshaled, err := json.Marshal(value)
		if err != nil {
			return err
		}
		bytes = marshaled
	}

	f.data[key] = bytes
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
	delete(f.data, key)
	return nil
}

func (f *fakeKV) DeleteMany(_ context.Context, keys []string) error {
	for _, key := range keys {
		_ = f.Delete(context.Background(), key)
	}
	return nil
}

func (f *fakeKV) Exists(_ context.Context, key string) (bool, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	_, ok := f.data[key]
	return ok, nil
}

func (f *fakeKV) ExistsMany(ctx context.Context, keys []string) ([]bool, error) {
	result := make([]bool, len(keys))
	for i, key := range keys {
		exists, err := f.Exists(ctx, key)
		if err != nil {
			return nil, err
		}
		result[i] = exists
	}
	return result, nil
}

func (f *fakeKV) GetKeysByPattern(_ context.Context, _ string) ([]string, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	keys := make([]string, 0, len(f.data))
	for key := range f.data {
		keys = append(keys, key)
	}
	return keys, nil
}

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

type stratzMock struct {
	server          *httptest.Server
	winProbRequests atomic.Int32
	notableRequests atomic.Int32
}

func newStratzMock(t *testing.T) *stratzMock {
	t.Helper()

	mock := &stratzMock{}

	mock.server = httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Header.Get("Authorization") != "Bearer test-token" {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				if r.Header.Get("User-Agent") != "STRATZ_API" {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				var body struct {
					Query string `json:"query"`
				}
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				w.Header().Set("Content-Type", "application/json")

				if body.Query == "" || !strings.HasPrefix(strings.TrimSpace(body.Query), "query") {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				if strings.Contains(body.Query, "liveWinRateValues") {
					mock.winProbRequests.Add(1)
					fmt.Fprint(
						w,
						`{"data":{"live":{"match":{"liveWinRateValues":[{"time":100,"winRate":41.2},{"time":200,"winRate":62.5}]}}}}`,
					)
					return
				}

				if strings.Contains(body.Query, "proSteamAccount") {
					mock.notableRequests.Add(1)
					fmt.Fprint(
						w,
						`{"data":{"live":{"match":{"players":[
							{"steamAccount":{"id":100,"name":"stratzName100","proSteamAccount":{"name":"ProOne","team":{"tag":"T1"}}}},
							{"steamAccount":{"id":200,"name":"regularPlayer","proSteamAccount":null}},
							{"steamAccount":{"id":300,"name":"streamerName","proSteamAccount":{"name":"StreamerPro","team":{"tag":"TT"}}}}
						]}}}}`,
					)
					return
				}

				w.WriteHeader(http.StatusBadRequest)
			},
		),
	)

	t.Cleanup(mock.server.Close)
	return mock
}

type opendotaMock struct {
	server         *httptest.Server
	recentRequests atomic.Int32
	heroesRequests atomic.Int32
	heroesStatus   atomic.Int32
	proRequests    atomic.Int32
}

func newOpendotaMock(t *testing.T) *opendotaMock {
	t.Helper()

	mock := &opendotaMock{}

	mux := http.NewServeMux()
	mux.HandleFunc(
		"GET /players/{accountID}/recentMatches",
		func(w http.ResponseWriter, _ *http.Request) {
			mock.recentRequests.Add(1)
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(
				w,
				`[
					{"match_id":1001,"hero_id":1,"kills":10,"deaths":2,"assists":8,"duration":1900,"player_slot":3,"radiant_win":false,"game_mode":22,"lobby_type":7,"start_time":1700000000},
					{"match_id":1000,"hero_id":2,"kills":1,"deaths":9,"assists":3,"duration":2400,"player_slot":130,"radiant_win":false,"game_mode":22,"lobby_type":7,"start_time":1699990000}
				]`,
			)
		},
	)
	mux.HandleFunc(
		"GET /constants/heroes",
		func(w http.ResponseWriter, _ *http.Request) {
			mock.heroesRequests.Add(1)
			if status := mock.heroesStatus.Load(); status != 0 {
				http.Error(w, "heroes unavailable", int(status))
				return
			}

			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(
				w,
				`{
					"1":{"id":1,"name":"npc_dota_hero_antimage","localized_name":"Anti-Mage"},
					"2":{"id":2,"name":"npc_dota_hero_axe","localized_name":"Axe"}
				}`,
			)
		},
	)
	mux.HandleFunc(
		"GET /proPlayers",
		func(w http.ResponseWriter, _ *http.Request) {
			mock.proRequests.Add(1)
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(
				w,
				`[
					{"account_id":100,"name":"Miracle-","team_name":"Nigma","team_tag":"Nigma"},
					{"account_id":300,"name":"StreamerProName","team_name":"TT","team_tag":"TT"}
				]`,
			)
		},
	)

	mock.server = httptest.NewServer(mux)
	t.Cleanup(mock.server.Close)
	return mock
}

func newTestStats(t *testing.T, stratzToken string) (*Stats, *stratzMock, *opendotaMock, *fakeKV) {
	t.Helper()

	stratzMock := newStratzMock(t)
	opendotaMock := newOpendotaMock(t)
	kvStore := newFakeKV()

	stratzClient := stratz.New(stratzToken, stratz.WithBaseURL(stratzMock.server.URL))
	opendotaClient := opendota.New(opendota.WithBaseURL(opendotaMock.server.URL))

	return New(stratzClient, opendotaClient, kvStore, testLogger()), stratzMock, opendotaMock, kvStore
}

func TestLastGame_ParsesMatchAndResolvesHero(t *testing.T) {
	stats, _, opendotaMock, _ := newTestStats(t, "test-token")

	game, err := stats.LastGame(context.Background(), 42)
	if err != nil {
		t.Fatalf("LastGame returned error: %v", err)
	}
	if game == nil {
		t.Fatal("LastGame returned nil")
	}

	if game.HeroName != "Anti-Mage" {
		t.Errorf("expected hero name Anti-Mage, got %q", game.HeroName)
	}
	if game.Kills != 10 || game.Deaths != 2 || game.Assists != 8 {
		t.Errorf("unexpected KDA: %d/%d/%d", game.Kills, game.Deaths, game.Assists)
	}
	if game.DurationS != 1900 {
		t.Errorf("expected duration 1900, got %d", game.DurationS)
	}
	// player_slot 3 < 128 means radiant; radiant_win false means loss.
	if game.Win {
		t.Error("expected loss for radiant player when radiant_win is false")
	}

	if opendotaMock.recentRequests.Load() != 1 {
		t.Errorf("expected 1 recentMatches request, got %d", opendotaMock.recentRequests.Load())
	}
	if opendotaMock.heroesRequests.Load() != 1 {
		t.Errorf("expected 1 heroes request, got %d", opendotaMock.heroesRequests.Load())
	}
}

func TestLastGame_ReturnsHeroLookupError(t *testing.T) {
	stats, _, opendotaMock, _ := newTestStats(t, "test-token")
	opendotaMock.heroesStatus.Store(http.StatusInternalServerError)

	game, err := stats.LastGame(context.Background(), 42)
	if err == nil {
		t.Fatal("expected LastGame to return the heroes error")
	}
	if game != nil {
		t.Errorf("expected nil game, got %+v", game)
	}
	if opendotaMock.recentRequests.Load() != 1 {
		t.Errorf("expected 1 successful recentMatches request, got %d", opendotaMock.recentRequests.Load())
	}
	if opendotaMock.heroesRequests.Load() != 1 {
		t.Errorf("expected 1 heroes request, got %d", opendotaMock.heroesRequests.Load())
	}
}

func TestRecentMatchWinLogic(t *testing.T) {
	tests := []struct {
		name       string
		playerSlot int
		radiantWin bool
		wantWin    bool
	}{
		{"radiant player radiant win", 3, true, true},
		{"radiant player dire win", 3, false, false},
		{"dire player dire win", 130, false, true},
		{"dire player radiant win", 130, true, false},
		{"boundary slot 127 is radiant", 127, true, true},
		{"boundary slot 128 is dire", 128, false, true},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				match := opendota.RecentMatch{PlayerSlot: tt.playerSlot, RadiantWin: tt.radiantWin}
				if got := match.Won(); got != tt.wantWin {
					t.Errorf("Won() = %v, want %v", got, tt.wantWin)
				}
			},
		)
	}
}

func TestHeroesCached(t *testing.T) {
	stats, _, opendotaMock, _ := newTestStats(t, "test-token")

	for i := 0; i < 3; i++ {
		if _, err := stats.LastGame(context.Background(), 42); err != nil {
			t.Fatalf("LastGame call %d returned error: %v", i, err)
		}
	}

	if opendotaMock.heroesRequests.Load() != 1 {
		t.Errorf(
			"expected heroes endpoint to be hit once, got %d",
			opendotaMock.heroesRequests.Load(),
		)
	}
}

func TestWinProbability_NormalizesPercentAndCaches(t *testing.T) {
	stats, stratzMock, _, _ := newTestStats(t, "test-token")

	prob, err := stats.WinProbability(context.Background(), 777)
	if err != nil {
		t.Fatalf("WinProbability returned error: %v", err)
	}
	if prob != 0.625 {
		t.Errorf("expected 0.625, got %v", prob)
	}

	// second call must be served from cache
	prob2, err := stats.WinProbability(context.Background(), 777)
	if err != nil {
		t.Fatalf("WinProbability second call returned error: %v", err)
	}
	if prob2 != prob {
		t.Errorf("expected cached value %v, got %v", prob, prob2)
	}
	if stratzMock.winProbRequests.Load() != 1 {
		t.Errorf(
			"expected 1 stratz request, got %d",
			stratzMock.winProbRequests.Load(),
		)
	}
}

func TestWinProbability_StratzDisabled(t *testing.T) {
	stats, stratzMock, _, _ := newTestStats(t, "")

	prob, err := stats.WinProbability(context.Background(), 777)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if prob != 0 {
		t.Errorf("expected 0 probability, got %v", prob)
	}
	if stratzMock.winProbRequests.Load() != 0 {
		t.Error("disabled stratz client must not hit the server")
	}
}

func TestWinProbability_StratzDisabledBypassesCache(t *testing.T) {
	stats, stratzMock, _, kvStore := newTestStats(t, "")
	ctx := context.Background()

	if err := setCached(
		ctx,
		kvStore,
		winProbabilityCachePrefix+"777",
		winProbabilityCacheTTL,
		0.625,
	); err != nil {
		t.Fatalf("failed to populate win probability cache: %v", err)
	}

	probability, err := stats.WinProbability(ctx, 777)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if probability != 0 {
		t.Errorf("expected 0 probability, got %v", probability)
	}
	if stratzMock.winProbRequests.Load() != 0 {
		t.Error("disabled stratz client must not hit the server")
	}
}

func TestStratzClient_DisabledReturnsErrDisabled(t *testing.T) {
	client := stratz.New("")

	if client.Enabled() {
		t.Error("expected Enabled() to be false with empty token")
	}

	if _, err := client.WinProbability(context.Background(), 1); !errors.Is(err, stratz.ErrDisabled) {
		t.Errorf("expected ErrDisabled, got %v", err)
	}
	if _, err := client.NotablePlayers(context.Background(), 1); !errors.Is(err, stratz.ErrDisabled) {
		t.Errorf("expected ErrDisabled, got %v", err)
	}
}

func TestNotablePlayers_CrossReferencesAndExcludesStreamer(t *testing.T) {
	stats, stratzMock, opendotaMock, _ := newTestStats(t, "test-token")

	names, err := stats.NotablePlayers(context.Background(), 555, "300")
	if err != nil {
		t.Fatalf("NotablePlayers returned error: %v", err)
	}

	if len(names) != 1 {
		t.Fatalf("expected 1 notable player, got %v", names)
	}
	// OpenDota name takes precedence over Stratz name for account 100.
	if names[0] != "Miracle-" {
		t.Errorf("expected Miracle-, got %q", names[0])
	}

	// second call must be served from cache
	names2, err := stats.NotablePlayers(context.Background(), 555, "300")
	if err != nil {
		t.Fatalf("NotablePlayers second call returned error: %v", err)
	}
	if len(names2) != 1 || names2[0] != "Miracle-" {
		t.Errorf("unexpected cached result: %v", names2)
	}
	if stratzMock.notableRequests.Load() != 1 {
		t.Errorf(
			"expected 1 stratz notable request, got %d",
			stratzMock.notableRequests.Load(),
		)
	}
	if opendotaMock.proRequests.Load() != 1 {
		t.Errorf(
			"expected 1 proPlayers request, got %d",
			opendotaMock.proRequests.Load(),
		)
	}
}

func TestNotablePlayers_CacheKeyIncludesStreamerAccountID(t *testing.T) {
	stats, _, _, _ := newTestStats(t, "test-token")
	ctx := context.Background()

	first, err := stats.NotablePlayers(ctx, 555, "300")
	if err != nil {
		t.Fatalf("first NotablePlayers call returned error: %v", err)
	}
	if len(first) != 1 || first[0] != "Miracle-" {
		t.Fatalf("unexpected first notable players result: %v", first)
	}

	second, err := stats.NotablePlayers(ctx, 555, "100")
	if err != nil {
		t.Fatalf("second NotablePlayers call returned error: %v", err)
	}
	if len(second) != 1 || second[0] != "StreamerProName" {
		t.Fatalf("unexpected second notable players result: %v", second)
	}
	for _, name := range second {
		if name == "Miracle-" {
			t.Error("second result included account 100")
		}
	}
}

func TestNotablePlayers_StratzDisabled(t *testing.T) {
	stats, stratzMock, _, _ := newTestStats(t, "")

	names, err := stats.NotablePlayers(context.Background(), 555, "300")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if names != nil {
		t.Errorf("expected nil names, got %v", names)
	}
	if stratzMock.notableRequests.Load() != 0 {
		t.Error("disabled stratz client must not hit the server")
	}
}

func TestNotablePlayers_StratzDisabledBypassesCache(t *testing.T) {
	stats, stratzMock, _, kvStore := newTestStats(t, "")
	ctx := context.Background()

	if err := setCached(
		ctx,
		kvStore,
		notablePlayersCachePrefix+"555:300",
		notablePlayersCacheTTL,
		[]string{"cached pro"},
	); err != nil {
		t.Fatalf("failed to populate notable players cache: %v", err)
	}

	names, err := stats.NotablePlayers(ctx, 555, "300")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if names != nil {
		t.Errorf("expected nil names, got %v", names)
	}
	if stratzMock.notableRequests.Load() != 0 {
		t.Error("disabled stratz client must not hit the server")
	}
}
