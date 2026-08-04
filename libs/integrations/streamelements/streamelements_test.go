package streamelements

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"
)

type legacyAuthLinker interface {
	GetAuthLink(redirectURL string) string
}

var _ legacyAuthLinker = (*StreamElements)(nil)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return f(request)
}

type observedBody struct {
	read bool
}

func (b *observedBody) Read([]byte) (int, error) {
	b.read = true
	return 0, errors.New("error response body must not be read")
}

func (b *observedBody) Close() error {
	return nil
}

type fakeTokenStore struct {
	mu sync.Mutex

	tokens      Tokens
	getErr      error
	updateErr   error
	getCalls    int
	updateCalls []Tokens
	getChannels []string
	setChannels []string
}

func (s *fakeTokenStore) GetTokens(_ context.Context, channelID string) (Tokens, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.getCalls++
	s.getChannels = append(s.getChannels, channelID)
	return s.tokens, s.getErr
}

func (s *fakeTokenStore) UpdateTokens(_ context.Context, channelID string, tokens Tokens) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.updateCalls = append(s.updateCalls, tokens)
	s.setChannels = append(s.setChannels, channelID)
	if s.updateErr == nil {
		s.tokens = tokens
	}
	return s.updateErr
}

func (s *fakeTokenStore) replace(tokens Tokens) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tokens = tokens
}

func (s *fakeTokenStore) snapshot() (Tokens, int, []Tokens) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.tokens, s.getCalls, append([]Tokens(nil), s.updateCalls...)
}

func (s *fakeTokenStore) channels() ([]string, []string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]string(nil), s.getChannels...), append([]string(nil), s.setChannels...)
}

type fakeLocker struct {
	before func()
	err    error
	keys   []string
}

func (l *fakeLocker) WithLock(
	ctx context.Context,
	key string,
	fn func(context.Context) error,
) error {
	l.keys = append(l.keys, key)
	if l.before != nil {
		l.before()
	}
	if l.err != nil {
		return l.err
	}
	return fn(ctx)
}

func TestStaticOAuthAuthorizationAndCodeExchange(t *testing.T) {
	var gotForm url.Values
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/oauth2/token" {
			t.Errorf("path = %q, want /oauth2/token", r.URL.Path)
		}
		if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
			t.Errorf("Content-Type = %q", r.Header.Get("Content-Type"))
		}
		if err := r.ParseForm(); err != nil {
			t.Errorf("ParseForm() error = %v", err)
		}
		gotForm = r.PostForm
		_, _ = w.Write([]byte(`{"access_token":"access","refresh_token":"refresh","token_type":"Bearer","expires_in":3600}`))
	}))
	defer server.Close()

	client := NewStatic(
		"client-id",
		"client-secret",
		WithBaseURL(server.URL),
		WithHTTPClient(server.Client()),
	)
	authLink, err := client.GetAuthLinkWithState("https://twir.example/callback", "csrf-state")
	if err != nil {
		t.Fatalf("GetAuthLinkWithState() error = %v", err)
	}
	parsed, err := url.Parse(authLink)
	if err != nil {
		t.Fatalf("Parse(auth link) error = %v", err)
	}
	if parsed.Path != "/oauth2/authorize" {
		t.Fatalf("auth path = %q, want /oauth2/authorize", parsed.Path)
	}
	wantQuery := url.Values{
		"client_id":     {"client-id"},
		"redirect_uri":  {"https://twir.example/callback"},
		"response_type": {"code"},
		"scope":         {"channel:read bot:read tips:read"},
		"state":         {"csrf-state"},
	}
	if !reflect.DeepEqual(parsed.Query(), wantQuery) {
		t.Fatalf("auth query = %#v, want %#v", parsed.Query(), wantQuery)
	}

	tokens, err := client.ExchangeCode(context.Background(), "auth-code", "https://twir.example/callback")
	if err != nil {
		t.Fatalf("ExchangeCode() error = %v", err)
	}
	if tokens.AccessToken != "access" || tokens.RefreshToken != "refresh" {
		t.Fatalf("ExchangeCode() tokens = %#v", tokens)
	}
	wantForm := url.Values{
		"grant_type":    {"authorization_code"},
		"client_id":     {"client-id"},
		"client_secret": {"client-secret"},
		"code":          {"auth-code"},
		"redirect_uri":  {"https://twir.example/callback"},
	}
	if !reflect.DeepEqual(gotForm, wantForm) {
		t.Fatalf("exchange form = %#v, want %#v", gotForm, wantForm)
	}
}

func TestLegacyGetAuthLinkRetainsExactMethodSignature(t *testing.T) {
	client := New("client-id", "client-secret")
	var methodValue func(string) string = client.GetAuthLink
	link := methodValue("https://twir.example/callback")
	parsed, err := url.Parse(link)
	if err != nil {
		t.Fatalf("Parse(auth link) error = %v", err)
	}
	if parsed.Query().Get("state") != "" {
		t.Fatalf("legacy auth link state = %q, want empty", parsed.Query().Get("state"))
	}
}

func TestGetAuthLinkWithStateRejectsBlankState(t *testing.T) {
	client := New("client-id", "client-secret")
	if _, err := client.GetAuthLinkWithState("https://twir.example/callback", "   "); err == nil {
		t.Fatal("GetAuthLinkWithState() error = nil, want blank-state rejection")
	}
}

func TestAuthenticatedProfileCommandsAndTimersUseBearerToken(t *testing.T) {
	var paths []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		paths = append(paths, r.URL.Path)
		if r.Header.Get("Authorization") != "Bearer access" {
			t.Errorf("Authorization = %q, want bearer access", r.Header.Get("Authorization"))
		}
		switch r.URL.Path {
		case "/kappa/v2/channels/me":
			_, _ = w.Write([]byte(`{"_id":"profile-id","username":"streamer"}`))
		case "/kappa/v2/bot/commands/provider-channel":
			_, _ = w.Write([]byte(`[{"_id":"command-id","command":"hello"}]`))
		case "/kappa/v2/bot/timers/provider-channel":
			_, _ = w.Write([]byte(`[{"_id":"timer-id","name":"timer"}]`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client := NewAuthorized(
		"client-id", "client-secret", "twir-channel", "https://twir.example/callback",
		Tokens{AccessToken: "access", RefreshToken: "refresh"},
		&fakeTokenStore{}, &fakeLocker{},
		WithBaseURL(server.URL), WithHTTPClient(server.Client()),
	)
	profile, err := client.GetProfile(context.Background())
	if err != nil {
		t.Fatalf("GetProfile() error = %v", err)
	}
	commands, err := client.GetCommands(context.Background(), "provider-channel")
	if err != nil {
		t.Fatalf("GetCommands() error = %v", err)
	}
	timers, err := client.GetTimers(context.Background(), "provider-channel")
	if err != nil {
		t.Fatalf("GetTimers() error = %v", err)
	}
	if profile.ID != "profile-id" || len(commands) != 1 || commands[0].ID != "command-id" ||
		len(timers) != 1 || timers[0].ID != "timer-id" {
		t.Fatalf("decoded responses = profile %#v, commands %#v, timers %#v", profile, commands, timers)
	}
	wantPaths := []string{
		"/kappa/v2/channels/me",
		"/kappa/v2/bot/commands/provider-channel",
		"/kappa/v2/bot/timers/provider-channel",
	}
	if !reflect.DeepEqual(paths, wantPaths) {
		t.Fatalf("paths = %#v, want %#v", paths, wantPaths)
	}
}

func TestProviderErrorsDoNotExposeResponseBodies(t *testing.T) {
	const secretBody = "provider-secret-that-must-not-leak"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.FormValue("code") == "malformed" {
			_, _ = w.Write([]byte(`{"access_token":`))
			return
		}
		w.WriteHeader(http.StatusBadGateway)
		_, _ = w.Write([]byte(secretBody))
	}))
	defer server.Close()

	client := NewStatic("id", "secret", WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	for _, code := range []string{"rejected", "malformed"} {
		_, err := client.ExchangeCode(context.Background(), code, "https://twir.example/callback")
		if err == nil {
			t.Fatalf("ExchangeCode(%q) error = nil", code)
		}
		if strings.Contains(err.Error(), secretBody) {
			t.Fatalf("ExchangeCode(%q) error leaked response body: %v", code, err)
		}
	}
}

func TestNon2xxResponseIsRejectedWithoutReadingProviderBody(t *testing.T) {
	body := &observedBody{}
	client := NewStatic(
		"id",
		"secret",
		WithHTTPClient(&http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusBadGateway,
				Header:     make(http.Header),
				Body:       body,
			}, nil
		})}),
	)

	_, err := client.ExchangeCode(context.Background(), "code", "callback")
	if err == nil || !strings.Contains(err.Error(), "status 502") {
		t.Fatalf("ExchangeCode() error = %v, want sanitized status error", err)
	}
	if body.read {
		t.Fatal("ExchangeCode() read non-2xx provider body")
	}
}

func TestSuccessfulResponseBodyIsCappedAtOneMiB(t *testing.T) {
	const maxBody = 1 << 20
	for _, test := range []struct {
		name    string
		size    int
		wantErr bool
	}{
		{name: "exact limit", size: maxBody},
		{name: "overflow", size: maxBody + 1, wantErr: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			prefix := `{"access_token":"access","refresh_token":"refresh"}`
			body := prefix + strings.Repeat(" ", test.size-len(prefix))
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				_, _ = io.WriteString(w, body)
			}))
			defer server.Close()

			client := NewStatic("id", "secret", WithBaseURL(server.URL), WithHTTPClient(server.Client()))
			_, err := client.ExchangeCode(context.Background(), "code", "callback")
			if test.wantErr {
				if err == nil || !strings.Contains(err.Error(), "exceeds 1048576 bytes") {
					t.Fatalf("ExchangeCode() error = %v, want 1 MiB overflow error", err)
				}
			} else if err != nil {
				t.Fatalf("ExchangeCode() error at exact limit = %v", err)
			}
		})
	}
}

func TestUnauthorizedRefreshPersistsResolvedTokensBeforeRetry(t *testing.T) {
	for _, test := range []struct {
		name            string
		refreshResponse string
		wantRefresh     string
	}{
		{name: "rotated refresh token", refreshResponse: "new-refresh", wantRefresh: "new-refresh"},
		{name: "omitted refresh token falls back", refreshResponse: "", wantRefresh: "old-refresh"},
	} {
		t.Run(test.name, func(t *testing.T) {
			var apiCalls, refreshCalls int
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.URL.Path {
				case "/kappa/v2/channels/me":
					apiCalls++
					if r.Header.Get("Authorization") == "Bearer old-access" {
						w.WriteHeader(http.StatusUnauthorized)
						return
					}
					if r.Header.Get("Authorization") != "Bearer new-access" {
						t.Errorf("retry Authorization = %q", r.Header.Get("Authorization"))
					}
					_, _ = w.Write([]byte(`{"_id":"profile-id"}`))
				case "/oauth2/token":
					refreshCalls++
					if err := r.ParseForm(); err != nil {
						t.Errorf("ParseForm() error = %v", err)
					}
					want := url.Values{
						"grant_type":    {"refresh_token"},
						"client_id":     {"client-id"},
						"client_secret": {"client-secret"},
						"refresh_token": {"old-refresh"},
					}
					if !reflect.DeepEqual(r.PostForm, want) {
						t.Errorf("refresh form = %#v, want %#v", r.PostForm, want)
					}
					_, _ = fmt.Fprintf(w, `{"access_token":"new-access","refresh_token":%q}`, test.refreshResponse)
				default:
					http.NotFound(w, r)
				}
			}))
			defer server.Close()

			store := &fakeTokenStore{tokens: Tokens{AccessToken: "old-access", RefreshToken: "old-refresh"}}
			locker := &fakeLocker{}
			client := NewAuthorized(
				"client-id", "client-secret", "channel-id", "https://twir.example/callback",
				Tokens{AccessToken: "old-access", RefreshToken: "old-refresh"}, store, locker,
				WithBaseURL(server.URL), WithHTTPClient(server.Client()),
			)

			profile, err := client.GetProfile(context.Background())
			if err != nil {
				t.Fatalf("GetProfile() error = %v", err)
			}
			if profile.ID != "profile-id" {
				t.Fatalf("profile ID = %q", profile.ID)
			}
			_, getCalls, updates := store.snapshot()
			wantTokens := Tokens{AccessToken: "new-access", RefreshToken: test.wantRefresh}
			if getCalls != 1 || !reflect.DeepEqual(updates, []Tokens{wantTokens}) {
				t.Fatalf("store calls = get %d updates %#v, want one atomic update %#v", getCalls, updates, wantTokens)
			}
			getChannels, setChannels := store.channels()
			if !reflect.DeepEqual(getChannels, []string{"channel-id"}) ||
				!reflect.DeepEqual(setChannels, []string{"channel-id"}) {
				t.Fatalf("store channel IDs = get %#v update %#v", getChannels, setChannels)
			}
			if !reflect.DeepEqual(client.tokens, wantTokens) {
				t.Fatalf("in-memory tokens = %#v, want %#v", client.tokens, wantTokens)
			}
			if apiCalls != 2 || refreshCalls != 1 {
				t.Fatalf("calls = api %d refresh %d, want api 2 refresh 1", apiCalls, refreshCalls)
			}
			wantKey := "twir:integration-token-refresh:streamelements:channel-id"
			if !reflect.DeepEqual(locker.keys, []string{wantKey}) {
				t.Fatalf("lock keys = %#v, want [%q]", locker.keys, wantKey)
			}
		})
	}
}

func TestUnauthorizedRereadsTokensAndSkipsRedundantRefresh(t *testing.T) {
	var apiCalls, refreshCalls int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/kappa/v2/channels/me":
			apiCalls++
			if r.Header.Get("Authorization") == "Bearer old-access" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			if r.Header.Get("Authorization") != "Bearer other-access" {
				t.Errorf("retry Authorization = %q", r.Header.Get("Authorization"))
			}
			_, _ = w.Write([]byte(`{"_id":"profile-id"}`))
		case "/oauth2/token":
			refreshCalls++
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	store := &fakeTokenStore{tokens: Tokens{AccessToken: "old-access", RefreshToken: "old-refresh"}}
	locker := &fakeLocker{before: func() {
		store.replace(Tokens{AccessToken: "other-access", RefreshToken: "other-refresh"})
	}}
	client := NewAuthorized(
		"id", "secret", "channel-id", "callback",
		Tokens{AccessToken: "old-access", RefreshToken: "old-refresh"}, store, locker,
		WithBaseURL(server.URL), WithHTTPClient(server.Client()),
	)

	if _, err := client.GetProfile(context.Background()); err != nil {
		t.Fatalf("GetProfile() error = %v", err)
	}
	_, getCalls, updates := store.snapshot()
	if apiCalls != 2 || refreshCalls != 0 || getCalls != 1 || len(updates) != 0 {
		t.Fatalf("calls = api %d refresh %d get %d updates %d", apiCalls, refreshCalls, getCalls, len(updates))
	}
	want := Tokens{AccessToken: "other-access", RefreshToken: "other-refresh"}
	if !reflect.DeepEqual(client.tokens, want) {
		t.Fatalf("in-memory tokens = %#v, want reread %#v", client.tokens, want)
	}
}

func TestSecondUnauthorizedReturnsErrUnauthorizedWithoutLooping(t *testing.T) {
	var apiCalls, refreshCalls int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/kappa/v2/channels/me":
			apiCalls++
			w.WriteHeader(http.StatusUnauthorized)
		case "/oauth2/token":
			refreshCalls++
			_, _ = w.Write([]byte(`{"access_token":"new-access","refresh_token":"new-refresh"}`))
		}
	}))
	defer server.Close()

	store := &fakeTokenStore{tokens: Tokens{AccessToken: "old-access", RefreshToken: "old-refresh"}}
	client := NewAuthorized(
		"id", "secret", "channel-id", "callback",
		Tokens{AccessToken: "old-access", RefreshToken: "old-refresh"}, store, &fakeLocker{},
		WithBaseURL(server.URL), WithHTTPClient(server.Client()),
	)

	_, err := client.GetProfile(context.Background())
	if !errors.Is(err, ErrUnauthorized) {
		t.Fatalf("GetProfile() error = %v, want ErrUnauthorized", err)
	}
	if apiCalls != 2 || refreshCalls != 1 {
		t.Fatalf("calls = api %d refresh %d, want api 2 refresh 1", apiCalls, refreshCalls)
	}
}

func TestFailedPersistenceDoesNotReplaceInMemoryTokens(t *testing.T) {
	var apiCalls int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/kappa/v2/channels/me":
			apiCalls++
			w.WriteHeader(http.StatusUnauthorized)
		case "/oauth2/token":
			_, _ = w.Write([]byte(`{"access_token":"new-access","refresh_token":"new-refresh"}`))
		}
	}))
	defer server.Close()

	wantErr := errors.New("database unavailable")
	store := &fakeTokenStore{
		tokens:    Tokens{AccessToken: "old-access", RefreshToken: "old-refresh"},
		updateErr: wantErr,
	}
	client := NewAuthorized(
		"id", "secret", "channel-id", "callback",
		Tokens{AccessToken: "old-access", RefreshToken: "old-refresh"}, store, &fakeLocker{},
		WithBaseURL(server.URL), WithHTTPClient(server.Client()),
	)

	_, err := client.GetProfile(context.Background())
	if !errors.Is(err, wantErr) {
		t.Fatalf("GetProfile() error = %v, want persistence error", err)
	}
	wantTokens := Tokens{AccessToken: "old-access", RefreshToken: "old-refresh"}
	if !reflect.DeepEqual(client.tokens, wantTokens) {
		t.Fatalf("in-memory tokens = %#v, want unchanged %#v", client.tokens, wantTokens)
	}
	if apiCalls != 1 {
		t.Fatalf("API calls = %d, want no retry after failed persistence", apiCalls)
	}
}

func TestRefreshRejectionIsSanitized(t *testing.T) {
	const secretBody = "refresh-token-is-invalid-and-secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/kappa/v2/channels/me" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(secretBody))
	}))
	defer server.Close()

	store := &fakeTokenStore{tokens: Tokens{AccessToken: "old", RefreshToken: "refresh"}}
	client := NewAuthorized(
		"id", "secret", "channel", "callback", store.tokens, store, &fakeLocker{},
		WithBaseURL(server.URL), WithHTTPClient(server.Client()),
	)
	_, err := client.GetProfile(context.Background())
	if err == nil || strings.Contains(err.Error(), secretBody) {
		t.Fatalf("GetProfile() error = %v, want sanitized refresh rejection", err)
	}
}

func TestDefaultHTTPClientHasFifteenSecondTimeout(t *testing.T) {
	client := NewAuthorized(
		"id", "secret", "channel", "callback", Tokens{}, &fakeTokenStore{}, &fakeLocker{},
	)
	if client.httpClient.Timeout != 15*time.Second {
		t.Fatalf("HTTP timeout = %v, want 15s", client.httpClient.Timeout)
	}
}
