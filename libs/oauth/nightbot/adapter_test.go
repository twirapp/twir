package nightbot

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/twirapp/twir/libs/oauth"
	channelsintegrations "github.com/twirapp/twir/libs/repositories/channels_integrations"
	channelModel "github.com/twirapp/twir/libs/repositories/channels_integrations/model"
	integrations "github.com/twirapp/twir/libs/repositories/integrations"
	integrationsmodel "github.com/twirapp/twir/libs/repositories/integrations/model"
)

func TestToken_refreshes_nightbot_with_form_credentials_and_preserves_omitted_refresh_token(t *testing.T) {
	// Given
	address := nightbotRedisAddress(t)
	repository := &nightbotRepositoryFake{integration: channelModel.ChannelIntegration{ID: "integration-1", AccessToken: stringPointer("old-access"), RefreshToken: stringPointer("old-refresh")}}
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Header.Get("Authorization") != "" || request.ParseForm() != nil || request.Form.Get("grant_type") != "refresh_token" || request.Form.Get("client_id") != "client-id" || request.Form.Get("client_secret") != "client-secret" || request.Form.Get("refresh_token") != "old-refresh" { t.Fatal("unexpected Nightbot request") }
		_, _ = writer.Write([]byte(`{"access_token":"new-access","expires_in":3600}`))
	}))
	t.Cleanup(server.Close)
	source := newNightbotSource(t, address, server.URL, repository, nightbotIntegrationSettings())

	// When
	credential, err := source.Token(context.Background(), "channel-1")

	// Then
	if err != nil || credential.AccessToken != "new-access" || repository.updated.RefreshToken != nil || repository.updated.Enabled == nil || !*repository.updated.Enabled { t.Fatalf("credential/update = %#v/%#v, %v", credential, repository.updated, err) }
}

func TestToken_returns_redacted_typed_error_when_credentials_missing_or_response_invalid(t *testing.T) {
	for _, test := range []struct { name string; settings integrations.Repository; body string; want error }{
		{name: "missing credentials", settings: integrationSettingsFake{}, body: "", want: ErrMissingClientCredentials},
		{name: "invalid json", settings: nightbotIntegrationSettings(), body: "provider-secret"},
	} {
		t.Run(test.name, func(t *testing.T) {
			// Given
			address := nightbotRedisAddress(t)
			server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) { _, _ = writer.Write([]byte(test.body)) }))
			t.Cleanup(server.Close)
			source := newNightbotSource(t, address, server.URL, &nightbotRepositoryFake{integration: channelModel.ChannelIntegration{ID: "integration-1", AccessToken: stringPointer("old-access"), RefreshToken: stringPointer("old-refresh")}}, test.settings)

			// When
			_, err := source.Token(context.Background(), "channel-1")

			// Then
			if !errors.Is(err, oauth.ErrRefresh) || (test.want != nil && !errors.Is(err, test.want)) || strings.Contains(err.Error(), "provider-secret") { t.Fatalf("error = %v", err) }
		})
	}
}

func newNightbotSource(t *testing.T, address, tokenURL string, repository channelsintegrations.Repository, settings integrations.Repository) TokenSource {
	t.Helper()
	client := redis.NewClient(&redis.Options{Addr: address})
	t.Cleanup(func() { _ = client.Close() })
	source, err := NewTokenSource(SourceOptions{Redis: client, HTTPClient: http.DefaultClient, TokenURL: tokenURL}, repository, settings)
	if err != nil { t.Fatal(err) }
	return source
}

func nightbotRedisAddress(t *testing.T) string { t.Helper(); address := os.Getenv("TWIR_INTEGRATIONS_TEST_REDIS_ADDR"); if address == "" { t.Skip("TWIR_INTEGRATIONS_TEST_REDIS_ADDR is required") }; return address }
func stringPointer(value string) *string { return &value }
func nightbotIntegrationSettings() integrations.Repository { return integrationSettingsFake{integration: integrationsmodel.Integration{ClientID: stringPointer("client-id"), ClientSecret: stringPointer("client-secret")}} }

type integrationSettingsFake struct { integration integrationsmodel.Integration }
func (fake integrationSettingsFake) GetByService(context.Context, integrationsmodel.Service) (integrationsmodel.Integration, error) { return fake.integration, nil }

type nightbotRepositoryFake struct { mu sync.Mutex; integration channelModel.ChannelIntegration; updated channelsintegrations.UpdateInput }
func (fake *nightbotRepositoryFake) GetByChannelAndService(context.Context, string, integrationsmodel.Service) (channelModel.ChannelIntegration, error) { fake.mu.Lock(); defer fake.mu.Unlock(); return fake.integration, nil }
func (fake *nightbotRepositoryFake) Update(_ context.Context, _ string, input channelsintegrations.UpdateInput) error { fake.mu.Lock(); defer fake.mu.Unlock(); fake.updated = input; fake.integration.AccessToken = input.AccessToken; if input.RefreshToken != nil { fake.integration.RefreshToken = input.RefreshToken }; return nil }
func (*nightbotRepositoryFake) Create(context.Context, channelsintegrations.CreateInput) (channelModel.ChannelIntegration, error) { return channelModel.Nil, errors.New("unexpected Create") }
