package gsi

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	cfg "github.com/twirapp/twir/libs/config"
	dotarepository "github.com/twirapp/twir/libs/repositories/dota"
	"github.com/twirapp/twir/libs/repositories/dota/model"
	"golang.org/x/time/rate"
)

type fakeRepo struct {
	settings map[string]model.ChannelDotaSettings
}

var _ dotarepository.Repository = (*fakeRepo)(nil)

func (f *fakeRepo) GetByGsiToken(
	_ context.Context,
	token string,
) (model.ChannelDotaSettings, error) {
	settings, ok := f.settings[token]
	if !ok {
		return model.Nil, dotarepository.ErrNotFound
	}
	return settings, nil
}

func (f *fakeRepo) GetByChannelID(
	context.Context,
	uuid.UUID,
) (model.ChannelDotaSettings, error) {
	return model.Nil, dotarepository.ErrNotFound
}

func (f *fakeRepo) Create(
	context.Context,
	dotarepository.CreateInput,
) (model.ChannelDotaSettings, error) {
	return model.Nil, nil
}

func (f *fakeRepo) Update(
	context.Context,
	uuid.UUID,
	dotarepository.UpdateInput,
) (model.ChannelDotaSettings, error) {
	return model.Nil, nil
}

func (f *fakeRepo) UpdateMatchResult(
	context.Context,
	uuid.UUID,
	bool,
	int,
) (model.ChannelDotaSettings, error) {
	return model.Nil, nil
}

func (f *fakeRepo) ApplyMatchResultOnce(
	context.Context,
	dotarepository.ApplyMatchResultInput,
) (model.ChannelDotaSettings, error) {
	return model.Nil, errors.New("not implemented")
}

func (f *fakeRepo) GetMatchState(context.Context, uuid.UUID) (model.MatchState, error) {
	return model.MatchState{}, errors.New("not implemented")
}

func (f *fakeRepo) ApplyMatchStateTransition(
	context.Context,
	dotarepository.ApplyMatchStateTransitionInput,
) (bool, error) {
	return false, errors.New("not implemented")
}

func (f *fakeRepo) ClaimPredictionActions(
	context.Context,
	dotarepository.ClaimPredictionActionsInput,
) ([]model.ClaimedOutboxAction, error) {
	return nil, errors.New("not implemented")
}

func (f *fakeRepo) CompletePredictionAction(context.Context, uuid.UUID, uuid.UUID) error {
	return errors.New("not implemented")
}

func (f *fakeRepo) RetryPredictionAction(
	context.Context,
	uuid.UUID,
	uuid.UUID,
	time.Time,
) error {
	return errors.New("not implemented")
}

func (f *fakeRepo) ResetSession(
	context.Context,
	uuid.UUID,
) (model.ChannelDotaSettings, error) {
	return model.Nil, nil
}

func (f *fakeRepo) RegenerateGsiToken(
	context.Context,
	uuid.UUID,
) (model.ChannelDotaSettings, error) {
	return model.Nil, nil
}

type processCall struct {
	channelID uuid.UUID
	payload   Payload
}

type fakeProcessor struct {
	mu    sync.Mutex
	calls []processCall
	err   error
}

func (f *fakeProcessor) Process(
	_ context.Context,
	channelID uuid.UUID,
	payload Payload,
) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.calls = append(f.calls, processCall{channelID: channelID, payload: payload})
	return f.err
}

func (f *fakeProcessor) lastCall() (processCall, bool) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if len(f.calls) == 0 {
		return processCall{}, false
	}
	return f.calls[len(f.calls)-1], true
}

const testToken = "test-token"

var testChannelID = uuid.MustParse("11111111-1111-1111-1111-111111111111")

func newTestServer(t *testing.T, processor *fakeProcessor, opts ...Option) *httptest.Server {
	t.Helper()

	repo := &fakeRepo{
		settings: map[string]model.ChannelDotaSettings{
			testToken: {ChannelID: testChannelID, GsiToken: testToken, Enabled: true},
		},
	}

	srv := New(
		cfg.Config{DotaHttpPort: 0},
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		repo,
		processor,
		opts...,
	)

	ts := httptest.NewServer(srv.httpServer.Handler)
	t.Cleanup(ts.Close)

	return ts
}

func postFixture(t *testing.T, ts *httptest.Server, token, fixture string) *http.Response {
	t.Helper()

	body, err := os.ReadFile("testdata/" + fixture)
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}

	resp, err := http.Post(ts.URL+"/gsi/"+token, "application/json", strings.NewReader(string(body)))
	if err != nil {
		t.Fatalf("post: %v", err)
	}

	return resp
}

func TestValidTokenAndPayload(t *testing.T) {
	processor := &fakeProcessor{}
	ts := newTestServer(t, processor)

	resp := postFixture(t, ts, testToken, "in_game.json")
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	call, ok := processor.lastCall()
	if !ok {
		t.Fatal("expected processor to be called")
	}
	if call.channelID != testChannelID {
		t.Fatalf("expected channelID %s, got %s", testChannelID, call.channelID)
	}
	if call.payload.Map == nil {
		t.Fatal("expected map to be present")
	}
	if call.payload.Map.MatchID != 7812345678 {
		t.Fatalf("expected matchid 7812345678, got %d", call.payload.Map.MatchID)
	}
	if call.payload.Map.GameState != GameStateInProgress {
		t.Fatalf("expected game state %s, got %s", GameStateInProgress, call.payload.Map.GameState)
	}
	if len(call.payload.Events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(call.payload.Events))
	}

	if call.payload.Events[1].PlayerID == nil {
		t.Fatal("expected aegis player_id to be present")
	}
	if *call.payload.Events[1].PlayerID != 2 {
		t.Fatalf("expected aegis player_id 2, got %d", *call.payload.Events[1].PlayerID)
	}
}

func TestMenuPayloadWithoutMap(t *testing.T) {
	processor := &fakeProcessor{}
	ts := newTestServer(t, processor)

	resp := postFixture(t, ts, testToken, "menu.json")
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	call, ok := processor.lastCall()
	if !ok {
		t.Fatal("expected processor to be called")
	}
	if call.payload.Map != nil {
		t.Fatal("expected map to be nil in menu payload")
	}
}

func TestUnknownToken(t *testing.T) {
	processor := &fakeProcessor{}
	ts := newTestServer(t, processor)

	resp := postFixture(t, ts, "unknown-token", "in_game.json")
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}
	if _, ok := processor.lastCall(); ok {
		t.Fatal("expected processor not to be called")
	}
}

func TestMismatchedAuthTokenInPayload(t *testing.T) {
	processor := &fakeProcessor{}
	ts := newTestServer(t, processor)

	resp, err := http.Post(
		ts.URL+"/gsi/"+testToken,
		"application/json",
		strings.NewReader(`{"auth": {"token": "other-token"}}`),
	)
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}
}

func TestMissingAuthTokenInPayloadStillAccepted(t *testing.T) {
	processor := &fakeProcessor{}
	ts := newTestServer(t, processor)

	resp, err := http.Post(
		ts.URL+"/gsi/"+testToken,
		"application/json",
		strings.NewReader(`{"provider": {"name": "Dota 2"}}`),
	)
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestMalformedJSON(t *testing.T) {
	processor := &fakeProcessor{}
	ts := newTestServer(t, processor)

	resp, err := http.Post(
		ts.URL+"/gsi/"+testToken,
		"application/json",
		strings.NewReader(`{not json`),
	)
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
	if _, ok := processor.lastCall(); ok {
		t.Fatal("expected processor not to be called")
	}
}

func TestRateLimit(t *testing.T) {
	processor := &fakeProcessor{}
	ts := newTestServer(t, processor, WithRateLimit(rate.Every(time.Hour), 1))

	resp := postFixture(t, ts, testToken, "in_game.json")
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 for first request, got %d", resp.StatusCode)
	}

	resp = postFixture(t, ts, testToken, "in_game.json")
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusTooManyRequests {
		t.Fatalf("expected 429 for second request, got %d", resp.StatusCode)
	}
}

func TestHealth(t *testing.T) {
	processor := &fakeProcessor{}
	ts := newTestServer(t, processor)

	resp, err := http.Get(ts.URL + "/health")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	if string(body) != "OK" {
		t.Fatalf("expected body OK, got %q", string(body))
	}
}
