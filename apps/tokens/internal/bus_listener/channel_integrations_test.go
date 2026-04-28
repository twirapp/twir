package bus_listener

import (
	"context"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"

	buscoretokens "github.com/twirapp/twir/libs/bus-core/tokens"
	channelsintegrations "github.com/twirapp/twir/libs/repositories/channels_integrations"
	channelsintegrationsmodel "github.com/twirapp/twir/libs/repositories/channels_integrations/model"
	channelsintegrationsspotify "github.com/twirapp/twir/libs/repositories/channels_integrations_spotify"
	channelsintegrationsspotifymodel "github.com/twirapp/twir/libs/repositories/channels_integrations_spotify/model"
	integrationsrepo "github.com/twirapp/twir/libs/repositories/integrations"
	integrationsmodel "github.com/twirapp/twir/libs/repositories/integrations/model"
	"github.com/google/uuid"
)

func TestRequestChannelIntegrationToken_SpotifyRefreshesViaTokensService(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if got := r.Header.Get("Authorization"); got != "Basic "+base64.StdEncoding.EncodeToString([]byte("client-id:client-secret")) {
			t.Fatalf("unexpected authorization header: %q", got)
		}
		if err := r.ParseForm(); err != nil {
			t.Fatalf("parse form: %v", err)
		}
		if r.Form.Get("grant_type") != "refresh_token" {
			t.Fatalf("unexpected grant_type: %q", r.Form.Get("grant_type"))
		}
		if r.Form.Get("refresh_token") != "spotify-refresh" {
			t.Fatalf("unexpected refresh token: %q", r.Form.Get("refresh_token"))
		}

		_, _ = w.Write([]byte(`{"access_token":"spotify-access-new","refresh_token":"spotify-refresh-new"}`))
	}))
	t.Cleanup(server.Close)

	spotifyRepo := &fakeSpotifyChannelIntegrationsRepository{
		integration: channelsintegrationsspotifymodel.ChannelIntegrationSpotify{
			ID:           uuid.New(),
			ChannelID:    "channel-1",
			AccessToken:  "spotify-access-old",
			RefreshToken: "spotify-refresh",
			Scopes:       []string{"user-read-playback-state"},
			Enabled:      true,
		},
	}

	impl := &tokensImpl{
		httpClient:               server.Client(),
		spotifyIntegrationsRepo: spotifyRepo,
		integrationsRepo: &fakeIntegrationsRepository{
			integration: integrationsmodel.Integration{
				Service:      integrationsmodel.ServiceSpotify,
				ClientID:     ptr("client-id"),
				ClientSecret: ptr("client-secret"),
			},
		},
		newMutex: func(name string) lockableMutex { return fakeMutex{} },
		spotifyTokenURL: server.URL,
	}

	resp, err := impl.RequestChannelIntegrationToken(context.Background(), buscoretokens.GetChannelIntegrationTokenRequest{
		ChannelID: "channel-1",
		Service:   integrationsmodel.ServiceSpotify,
	})
	if err != nil {
		t.Fatalf("RequestChannelIntegrationToken returned error: %v", err)
	}

	if resp.AccessToken != "spotify-access-new" {
		t.Fatalf("unexpected access token: %q", resp.AccessToken)
	}
	if spotifyRepo.updateCalls != 1 {
		t.Fatalf("expected one spotify update, got %d", spotifyRepo.updateCalls)
	}
	if spotifyRepo.updated.RefreshToken == nil || *spotifyRepo.updated.RefreshToken != "spotify-refresh-new" {
		t.Fatalf("unexpected persisted refresh token")
	}
}

func TestRequestChannelIntegrationToken_NightbotRefreshesViaTokensService(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Fatalf("parse form: %v", err)
		}
		if r.Form.Get("client_id") != "nightbot-client" {
			t.Fatalf("unexpected client_id: %q", r.Form.Get("client_id"))
		}
		if r.Form.Get("client_secret") != "nightbot-secret" {
			t.Fatalf("unexpected client_secret: %q", r.Form.Get("client_secret"))
		}
		if r.Form.Get("refresh_token") != "nightbot-refresh" {
			t.Fatalf("unexpected refresh token: %q", r.Form.Get("refresh_token"))
		}

		_, _ = w.Write([]byte(`{"access_token":"nightbot-access-new","refresh_token":"nightbot-refresh-new","expires_in":3600}`))
	}))
	t.Cleanup(server.Close)

	channelRepo := &fakeChannelIntegrationsRepository{
		integration: channelsintegrationsmodel.ChannelIntegration{
			ID:           "integration-1",
			ChannelID:    "channel-1",
			Enabled:      true,
			AccessToken:  ptr("nightbot-access-old"),
			RefreshToken: ptr("nightbot-refresh"),
		},
	}

	impl := &tokensImpl{
		httpClient:               server.Client(),
		channelIntegrationsRepo: channelRepo,
		integrationsRepo: &fakeIntegrationsRepository{
			integration: integrationsmodel.Integration{
				Service:      integrationsmodel.ServiceNightbot,
				ClientID:     ptr("nightbot-client"),
				ClientSecret: ptr("nightbot-secret"),
			},
		},
		newMutex: func(name string) lockableMutex { return fakeMutex{} },
		nightbotTokenURL: server.URL,
	}

	resp, err := impl.RequestChannelIntegrationToken(context.Background(), buscoretokens.GetChannelIntegrationTokenRequest{
		ChannelID: "channel-1",
		Service:   integrationsmodel.ServiceNightbot,
	})
	if err != nil {
		t.Fatalf("RequestChannelIntegrationToken returned error: %v", err)
	}

	if resp.AccessToken != "nightbot-access-new" {
		t.Fatalf("unexpected access token: %q", resp.AccessToken)
	}
	if channelRepo.updateCalls != 1 {
		t.Fatalf("expected one channel integration update, got %d", channelRepo.updateCalls)
	}
}

type fakeIntegrationsRepository struct {
	integration integrationsmodel.Integration
	service     integrationsmodel.Service
	getCalls    int
}

func (f *fakeIntegrationsRepository) GetByService(ctx context.Context, service integrationsmodel.Service) (integrationsmodel.Integration, error) {
	f.service = service
	f.getCalls++
	return f.integration, nil
}

type fakeSpotifyChannelIntegrationsRepository struct {
	integration channelsintegrationsspotifymodel.ChannelIntegrationSpotify
	updated     channelsintegrationsspotify.UpdateInput
	updateCalls int
}

func (f *fakeSpotifyChannelIntegrationsRepository) GetByChannelID(ctx context.Context, channelID string) (channelsintegrationsspotifymodel.ChannelIntegrationSpotify, error) {
	return f.integration, nil
}

func (f *fakeSpotifyChannelIntegrationsRepository) Create(ctx context.Context, input channelsintegrationsspotify.CreateInput) (channelsintegrationsspotifymodel.ChannelIntegrationSpotify, error) {
	panic("unexpected call")
}

func (f *fakeSpotifyChannelIntegrationsRepository) Update(ctx context.Context, id uuid.UUID, input channelsintegrationsspotify.UpdateInput) error {
	f.updateCalls++
	f.updated = input
	if input.AccessToken != nil {
		f.integration.AccessToken = *input.AccessToken
	}
	if input.RefreshToken != nil {
		f.integration.RefreshToken = *input.RefreshToken
	}
	return nil
}

func (f *fakeSpotifyChannelIntegrationsRepository) Delete(ctx context.Context, id uuid.UUID) error {
	panic("unexpected call")
}

type fakeChannelIntegrationsRepository struct {
	integration channelsintegrationsmodel.ChannelIntegration
	updated     channelsintegrations.UpdateInput
	updateCalls int
}

func (f *fakeChannelIntegrationsRepository) GetByChannelAndService(ctx context.Context, channelID string, service integrationsmodel.Service) (channelsintegrationsmodel.ChannelIntegration, error) {
	return f.integration, nil
}

func (f *fakeChannelIntegrationsRepository) Create(ctx context.Context, input channelsintegrations.CreateInput) (channelsintegrationsmodel.ChannelIntegration, error) {
	panic("unexpected call")
}

func (f *fakeChannelIntegrationsRepository) Update(ctx context.Context, id string, input channelsintegrations.UpdateInput) error {
	f.updateCalls++
	f.updated = input
	if input.AccessToken != nil {
		f.integration.AccessToken = input.AccessToken
	}
	if input.RefreshToken != nil {
		f.integration.RefreshToken = input.RefreshToken
	}
	return nil
}

func ptr[T any](v T) *T { return &v }

var _ integrationsrepo.Repository = (*fakeIntegrationsRepository)(nil)
var _ channelsintegrationsspotify.Repository = (*fakeSpotifyChannelIntegrationsRepository)(nil)
var _ channelsintegrations.Repository = (*fakeChannelIntegrationsRepository)(nil)
