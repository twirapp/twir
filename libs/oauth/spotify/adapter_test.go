package spotify

import (
	"context"
	"encoding/base64"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"
	"strings"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/twirapp/twir/libs/oauth"
	channelsintegrationsspotify "github.com/twirapp/twir/libs/repositories/channels_integrations_spotify"
	spotifyModel "github.com/twirapp/twir/libs/repositories/channels_integrations_spotify/model"
	integrations "github.com/twirapp/twir/libs/repositories/integrations"
	integrationsmodel "github.com/twirapp/twir/libs/repositories/integrations/model"
)

func TestToken_refreshes_spotify_with_basic_auth_and_preserves_omitted_refresh_token(t *testing.T) {
	// Given
	address := integrationRedisAddress(t)
	repository := &spotifyRepositoryFake{integration: spotifyModel.ChannelIntegrationSpotify{ID: uuid.New(), ChannelID: "channel-1", AccessToken: "old-access", RefreshToken: "old-refresh", Scopes: []string{"scope"}, UpdatedAt: time.Now().Add(-2 * time.Hour)}}
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost || request.Header.Get("Authorization") != "Basic "+base64.StdEncoding.EncodeToString([]byte("client-id:client-secret")) {
			t.Fatal("unexpected Spotify request authentication")
		}
		if err := request.ParseForm(); err != nil || request.Form.Get("grant_type") != "refresh_token" || request.Form.Get("refresh_token") != "old-refresh" {
			t.Fatal("unexpected Spotify request form")
		}
		_, _ = writer.Write([]byte(`{"access_token":"new-access","expires_in":3600}`))
	}))
	t.Cleanup(server.Close)
	source := newSpotifySource(t, address, server.URL, repository, spotifyIntegrationSettings())

	// When
	credential, err := source.Token(context.Background(), "channel-1")

	// Then
	if err != nil || credential.AccessToken != "new-access" || repository.updated.RefreshToken != nil {
		t.Fatalf("credential/update = %#v/%#v, %v", credential, repository.updated, err)
	}
}

func TestToken_returns_redacted_typed_error_for_provider_failure(t *testing.T) {
	for _, test := range []struct {
		name string
		body string
		status int
	}{
		{name: "non_2xx", body: "provider-secret", status: http.StatusUnauthorized},
		{name: "invalid_json", body: "{", status: http.StatusOK},
	} {
		t.Run(test.name, func(t *testing.T) {
			// Given
			address := integrationRedisAddress(t)
			server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
				writer.WriteHeader(test.status)
				_, _ = writer.Write([]byte(test.body))
			}))
			t.Cleanup(server.Close)
			source := newSpotifySource(t, address, server.URL, &spotifyRepositoryFake{integration: expiredSpotifyIntegration()}, spotifyIntegrationSettings())

			// When
			_, err := source.Token(context.Background(), "channel-1")

			// Then
			var providerError ProviderError
			if !errors.Is(err, oauth.ErrRefresh) || !errors.As(err, &providerError) || strings.Contains(err.Error(), "provider-secret") {
				t.Fatalf("error = %v", err)
			}
		})
	}
}

func TestToken_refreshes_once_across_two_runtimes(t *testing.T) {
	// Given
	address := integrationRedisAddress(t)
	repository := &spotifyRepositoryFake{integration: expiredSpotifyIntegration()}
	started, proceed := make(chan struct{}), make(chan struct{})
	var calls atomic.Int64
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		calls.Add(1); close(started); <-proceed
		_, _ = writer.Write([]byte(`{"access_token":"new-access","expires_in":3600}`))
	}))
	t.Cleanup(server.Close)
	first := newSpotifySource(t, address, server.URL, repository, spotifyIntegrationSettings())
	second := newSpotifySource(t, address, server.URL, repository, spotifyIntegrationSettings())
	results := make(chan error, 2)
	go func() { _, err := first.Token(context.Background(), "channel-1"); results <- err }()
	<-started
	go func() { _, err := second.Token(context.Background(), "channel-1"); results <- err }()

	// When
	close(proceed)
	err1, err2 := <-results, <-results

	// Then
	if err1 != nil || err2 != nil || calls.Load() != 1 { t.Fatalf("errors/calls = %v/%v/%d", err1, err2, calls.Load()) }
}

func newSpotifySource(t *testing.T, address, tokenURL string, repository channelsintegrationsspotify.Repository, settings integrations.Repository) TokenSource {
	t.Helper()
	client := redis.NewClient(&redis.Options{Addr: address})
	t.Cleanup(func() { _ = client.Close() })
	source, err := NewTokenSource(SourceOptions{Redis: client, HTTPClient: http.DefaultClient, TokenURL: tokenURL}, repository, settings)
	if err != nil { t.Fatal(err) }
	return source
}

func integrationRedisAddress(t *testing.T) string {
	t.Helper()
	address := os.Getenv("TWIR_INTEGRATIONS_TEST_REDIS_ADDR")
	if address == "" { t.Skip("TWIR_INTEGRATIONS_TEST_REDIS_ADDR is required") }
	return address
}

func spotifyIntegrationSettings() integrations.Repository { return integrationSettingsFake{integration: integrationsmodel.Integration{ClientID: pointer("client-id"), ClientSecret: pointer("client-secret")}} }
func expiredSpotifyIntegration() spotifyModel.ChannelIntegrationSpotify { return spotifyModel.ChannelIntegrationSpotify{ID: uuid.New(), ChannelID: "channel-1", AccessToken: "old-access", RefreshToken: "old-refresh", UpdatedAt: time.Now().Add(-2 * time.Hour)} }
func pointer(value string) *string { return &value }
type integrationSettingsFake struct { integration integrationsmodel.Integration }
func (fake integrationSettingsFake) GetByService(context.Context, integrationsmodel.Service) (integrationsmodel.Integration, error) { return fake.integration, nil }

type spotifyRepositoryFake struct { mu sync.Mutex; integration spotifyModel.ChannelIntegrationSpotify; updated channelsintegrationsspotify.UpdateInput }
func (fake *spotifyRepositoryFake) GetByChannelID(context.Context, string) (spotifyModel.ChannelIntegrationSpotify, error) { fake.mu.Lock(); defer fake.mu.Unlock(); return fake.integration, nil }
func (fake *spotifyRepositoryFake) Update(_ context.Context, _ uuid.UUID, input channelsintegrationsspotify.UpdateInput) error { fake.mu.Lock(); defer fake.mu.Unlock(); fake.updated = input; if input.AccessToken != nil { fake.integration.AccessToken = *input.AccessToken }; if input.RefreshToken != nil { fake.integration.RefreshToken = *input.RefreshToken }; fake.integration.UpdatedAt = time.Now(); return nil }
func (*spotifyRepositoryFake) Create(context.Context, channelsintegrationsspotify.CreateInput) (spotifyModel.ChannelIntegrationSpotify, error) { return spotifyModel.Nil, errors.New("unexpected Create") }
func (*spotifyRepositoryFake) Delete(context.Context, uuid.UUID) error { return errors.New("unexpected Delete") }
