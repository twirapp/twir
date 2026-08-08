package spotify_song_requests

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	cfg "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/entities/song_request_mode"
	songrequestssettingsentity "github.com/twirapp/twir/libs/entities/song_requests_settings"
	spotify_song_request "github.com/twirapp/twir/libs/entities/spotify_song_request"
	"github.com/twirapp/twir/libs/integrations/spotify"
	channelsintegrationsspotify "github.com/twirapp/twir/libs/repositories/channels_integrations_spotify"
	channelsintegrationsspotifymodel "github.com/twirapp/twir/libs/repositories/channels_integrations_spotify/model"
	songrequestssettingsrepository "github.com/twirapp/twir/libs/repositories/song_requests_settings"
	spotify_song_requests_repository "github.com/twirapp/twir/libs/repositories/spotify_song_requests"
)

const testChannelID = "test-channel"

type roundTripperFunc func(req *http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func fixtureResponse(req *http.Request, statusCode int, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(body)),
		Request:    req,
	}
}

func swapHTTPClient(t *testing.T, handler roundTripperFunc) {
	t.Helper()

	defaultClient := http.DefaultClient
	http.DefaultClient = &http.Client{Transport: handler}
	t.Cleanup(func() {
		http.DefaultClient = defaultClient
	})
}

type fakeSettingsRepository struct {
	settings songrequestssettingsentity.Settings
	err      error
}

func (r *fakeSettingsRepository) GetByChannelID(
	context.Context,
	string,
) (songrequestssettingsentity.Settings, error) {
	if r.err != nil {
		return songrequestssettingsentity.Nil, r.err
	}
	return r.settings, nil
}

func (r *fakeSettingsRepository) Upsert(
	context.Context,
	songrequestssettingsentity.Settings,
) (songrequestssettingsentity.Settings, error) {
	return songrequestssettingsentity.Nil, nil
}

func (r *fakeSettingsRepository) SetVolume(context.Context, string, int) error {
	return nil
}

type fakeIntegrationsRepository struct {
	integration channelsintegrationsspotifymodel.ChannelIntegrationSpotify
	err         error
}

func (r *fakeIntegrationsRepository) GetByChannelID(
	context.Context,
	string,
) (channelsintegrationsspotifymodel.ChannelIntegrationSpotify, error) {
	if r.err != nil {
		return channelsintegrationsspotifymodel.Nil, r.err
	}
	return r.integration, nil
}

func (r *fakeIntegrationsRepository) Create(
	context.Context,
	channelsintegrationsspotify.CreateInput,
) (channelsintegrationsspotifymodel.ChannelIntegrationSpotify, error) {
	return channelsintegrationsspotifymodel.Nil, nil
}

func (r *fakeIntegrationsRepository) Update(
	context.Context,
	uuid.UUID,
	channelsintegrationsspotify.UpdateInput,
) error {
	return nil
}

func (r *fakeIntegrationsRepository) Delete(context.Context, uuid.UUID) error {
	return nil
}

type fakeRequestsRepository struct {
	mu           sync.Mutex
	requests     []spotify_song_request.SpotifySongRequest
	created      []spotify_song_request.SpotifySongRequest
	statusByID   map[string]spotify_song_request.Status
	cancelledIDs []string
}

func newFakeRequestsRepository() *fakeRequestsRepository {
	return &fakeRequestsRepository{statusByID: map[string]spotify_song_request.Status{}}
}

func (r *fakeRequestsRepository) Create(
	_ context.Context,
	req spotify_song_request.SpotifySongRequest,
) (spotify_song_request.SpotifySongRequest, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.created = append(r.created, req)
	r.requests = append(r.requests, req)
	return req, nil
}

func (r *fakeRequestsRepository) GetByID(
	_ context.Context,
	id string,
) (spotify_song_request.SpotifySongRequest, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, request := range r.requests {
		if request.ID.String() == id {
			return request, nil
		}
	}
	return spotify_song_request.Nil, spotify_song_requests_repository.ErrNotFound
}

func (r *fakeRequestsRepository) GetActiveByChannel(
	_ context.Context,
	channelID string,
) ([]spotify_song_request.SpotifySongRequest, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var result []spotify_song_request.SpotifySongRequest
	for _, request := range r.requests {
		if request.ChannelID != channelID {
			continue
		}
		if status, ok := r.statusByID[request.ID.String()]; ok {
			request.Status = status
		}
		if request.Status == spotify_song_request.StatusQueued ||
			request.Status == spotify_song_request.StatusPlaying ||
			request.Status == spotify_song_request.StatusCancelledPendingSkip {
			result = append(result, request)
		}
	}
	return result, nil
}

func (r *fakeRequestsRepository) GetActiveByRequester(
	ctx context.Context,
	channelID string,
	requesterName string,
) ([]spotify_song_request.SpotifySongRequest, error) {
	active, err := r.GetActiveByChannel(ctx, channelID)
	if err != nil {
		return nil, err
	}
	var result []spotify_song_request.SpotifySongRequest
	for _, request := range active {
		if request.RequesterName == requesterName {
			result = append(result, request)
		}
	}
	return result, nil
}

func (r *fakeRequestsRepository) GetActiveChannels(context.Context) ([]string, error) {
	return nil, nil
}

func (r *fakeRequestsRepository) CountActiveByChannel(
	ctx context.Context,
	channelID string,
) (int64, error) {
	active, err := r.GetActiveByChannel(ctx, channelID)
	if err != nil {
		return 0, err
	}
	return int64(len(active)), nil
}

func (r *fakeRequestsRepository) CountActiveByRequester(
	ctx context.Context,
	channelID string,
	requesterName string,
) (int64, error) {
	active, err := r.GetActiveByRequester(ctx, channelID, requesterName)
	if err != nil {
		return 0, err
	}
	return int64(len(active)), nil
}

func (r *fakeRequestsRepository) ListByChannel(
	context.Context,
	string,
	int,
) ([]spotify_song_request.SpotifySongRequest, error) {
	return nil, nil
}

func (r *fakeRequestsRepository) UpdateStatus(
	_ context.Context,
	id string,
	status spotify_song_request.Status,
) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.statusByID[id] = status
	return nil
}

func (r *fakeRequestsRepository) CancelPendingSkip(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.cancelledIDs = append(r.cancelledIDs, id)
	return nil
}

func spotifySettings() songrequestssettingsentity.Settings {
	return songrequestssettingsentity.Settings{
		Enabled: true,
		Mode:    song_request_mode.ModeSpotify,
	}
}

func spotifyIntegration(scopes ...string) channelsintegrationsspotifymodel.ChannelIntegrationSpotify {
	return channelsintegrationsspotifymodel.ChannelIntegrationSpotify{
		ID:           uuid.New(),
		ChannelID:    testChannelID,
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		Scopes:       scopes,
	}
}

const searchResponseBody = `{"tracks":{"items":[{"id":"track-1","uri":"spotify:track:track-1","name":"Song","type":"track","artists":[{"name":"Artist"}],"album":{"name":"Album","images":[{"url":"https://img/1.jpg"}]},"duration_ms":180000}]}}`
const devicesResponseBody = `{"devices":[{"id":"device-1","name":"Desktop","type":"Computer","is_active":true,"is_restricted":false,"is_private_session":false}]}`

func newService(
	settingsRepo *fakeSettingsRepository,
	integrationsRepo *fakeIntegrationsRepository,
	requestsRepo *fakeRequestsRepository,
) *Service {
	return New(
		requestsRepo,
		settingsRepo,
		integrationsRepo,
		cfg.Config{SpotifyClientID: "id", SpotifySecret: "secret"},
		slog.Default(),
		nil,
		nil,
	)
}

func TestService_CreateRequest_with_spotify_link_fetches_track_directly(t *testing.T) {
	const trackID = "6IqfoZlee8TYVSUhIiug0P"
	var searchCalled bool
	swapHTTPClient(t, func(req *http.Request) (*http.Response, error) {
		switch req.URL.Path {
		case "/v1/search":
			searchCalled = true
			return fixtureResponse(req, http.StatusOK, searchResponseBody), nil
		case "/v1/tracks/" + trackID:
			body := `{"id":"` + trackID + `","uri":"spotify:track:` + trackID + `","name":"Linked Song","type":"track","artists":[{"name":"Linked Artist"}],"album":{"name":"Album","images":[]},"duration_ms":180000}`
			return fixtureResponse(req, http.StatusOK, body), nil
		case "/v1/me/player/devices":
			return fixtureResponse(req, http.StatusOK, devicesResponseBody), nil
		case "/v1/me/player/queue":
			return fixtureResponse(req, http.StatusNoContent, ""), nil
		default:
			t.Fatalf("unexpected request path %s", req.URL.Path)
			return nil, nil
		}
	})

	requestsRepo := newFakeRequestsRepository()
	service := newService(
		&fakeSettingsRepository{settings: spotifySettings()},
		&fakeIntegrationsRepository{
			integration: spotifyIntegration("user-read-playback-state", "user-modify-playback-state"),
		},
		requestsRepo,
	)

	request, err := service.CreateRequest(
		context.Background(),
		testChannelID,
		"user-1",
		"viewer",
		"Viewer",
		"chat",
		"https://open.spotify.com/track/"+trackID+"?si=1a13b96e0ad14bf0",
	)
	if err != nil {
		t.Fatalf("CreateRequest() error = %v", err)
	}
	if searchCalled {
		t.Fatal("search endpoint was called for a spotify link query")
	}
	if request.TrackID != trackID || request.Title != "Linked Song" || request.Artist != "Linked Artist" {
		t.Fatalf("request = %#v, want linked track fields", request)
	}
}

func TestService_CreateRequest_success(t *testing.T) {
	var queueAddedURI string
	swapHTTPClient(t, func(req *http.Request) (*http.Response, error) {
		switch req.URL.Path {
		case "/v1/search":
			return fixtureResponse(req, http.StatusOK, searchResponseBody), nil
		case "/v1/me/player/devices":
			return fixtureResponse(req, http.StatusOK, devicesResponseBody), nil
		case "/v1/me/player/queue":
			queueAddedURI = req.URL.Query().Get("uri")
			return fixtureResponse(req, http.StatusNoContent, ""), nil
		default:
			t.Fatalf("unexpected request path %s", req.URL.Path)
			return nil, nil
		}
	})

	requestsRepo := newFakeRequestsRepository()
	service := newService(
		&fakeSettingsRepository{settings: spotifySettings()},
		&fakeIntegrationsRepository{
			integration: spotifyIntegration("user-read-playback-state", "user-modify-playback-state"),
		},
		requestsRepo,
	)

	request, err := service.CreateRequest(
		context.Background(),
		testChannelID,
		"user-1",
		"viewer",
		"Viewer",
		"chat",
		"song",
	)
	if err != nil {
		t.Fatalf("CreateRequest() error = %v", err)
	}
	if queueAddedURI != "spotify:track:track-1" {
		t.Fatalf("queued uri = %q, want spotify:track:track-1", queueAddedURI)
	}
	if request.Status != spotify_song_request.StatusQueued {
		t.Fatalf("status = %q, want queued", request.Status)
	}
	if request.QueuePosition != 1 {
		t.Fatalf("queue position = %d, want 1", request.QueuePosition)
	}
	if len(requestsRepo.created) != 1 {
		t.Fatalf("created requests = %d, want 1", len(requestsRepo.created))
	}
}

func TestService_CreateRequest_rejects_duplicate_track(t *testing.T) {
	requestsRepo := newFakeRequestsRepository()
	requestsRepo.requests = append(
		requestsRepo.requests,
		spotify_song_request.SpotifySongRequest{
			ID:        uuid.New(),
			ChannelID: testChannelID,
			TrackURI:  "spotify:track:track-1",
			Status:    spotify_song_request.StatusQueued,
		},
	)

	service := newService(
		&fakeSettingsRepository{settings: spotifySettings()},
		&fakeIntegrationsRepository{
			integration: spotifyIntegration("user-modify-playback-state"),
		},
		requestsRepo,
	)

	swapHTTPClient(t, func(req *http.Request) (*http.Response, error) {
		return fixtureResponse(req, http.StatusOK, searchResponseBody), nil
	})

	_, err := service.CreateRequest(context.Background(), testChannelID, "", "viewer", "Viewer", "chat", "song")
	if !errors.Is(err, ErrTrackAlreadyInQueue) {
		t.Fatalf("CreateRequest() error = %v, want %v", err, ErrTrackAlreadyInQueue)
	}
}

func TestService_CreateRequest_rejects_youtube_mode(t *testing.T) {
	settings := spotifySettings()
	settings.Mode = song_request_mode.ModeYouTube

	service := newService(
		&fakeSettingsRepository{settings: settings},
		&fakeIntegrationsRepository{integration: spotifyIntegration()},
		newFakeRequestsRepository(),
	)

	_, err := service.CreateRequest(context.Background(), testChannelID, "", "viewer", "Viewer", "chat", "song")
	if !errors.Is(err, ErrNotSpotifyMode) {
		t.Fatalf("CreateRequest() error = %v, want %v", err, ErrNotSpotifyMode)
	}
}

func TestService_CreateRequest_rejects_missing_settings(t *testing.T) {
	service := newService(
		&fakeSettingsRepository{err: songrequestssettingsrepository.ErrNotFound},
		&fakeIntegrationsRepository{integration: spotifyIntegration()},
		newFakeRequestsRepository(),
	)

	_, err := service.CreateRequest(context.Background(), testChannelID, "", "viewer", "Viewer", "chat", "song")
	if !errors.Is(err, ErrNotSpotifyMode) {
		t.Fatalf("CreateRequest() error = %v, want %v", err, ErrNotSpotifyMode)
	}
}

func TestService_CreateRequest_rejects_not_connected(t *testing.T) {
	service := newService(
		&fakeSettingsRepository{settings: spotifySettings()},
		&fakeIntegrationsRepository{integration: channelsintegrationsspotifymodel.Nil},
		newFakeRequestsRepository(),
	)

	_, err := service.CreateRequest(context.Background(), testChannelID, "", "viewer", "Viewer", "chat", "song")
	if !errors.Is(err, spotify.ErrNotConnected) {
		t.Fatalf("CreateRequest() error = %v, want %v", err, spotify.ErrNotConnected)
	}
}

func TestService_CreateRequest_rejects_missing_scope(t *testing.T) {
	service := newService(
		&fakeSettingsRepository{settings: spotifySettings()},
		&fakeIntegrationsRepository{integration: spotifyIntegration("user-read-playback-state")},
		newFakeRequestsRepository(),
	)

	_, err := service.CreateRequest(context.Background(), testChannelID, "", "viewer", "Viewer", "chat", "song")
	if !errors.Is(err, spotify.ErrInsufficientScope) {
		t.Fatalf("CreateRequest() error = %v, want %v", err, spotify.ErrInsufficientScope)
	}
}

func TestService_CreateRequest_rejects_max_requests(t *testing.T) {
	settings := spotifySettings()
	settings.MaxRequests = 1

	requestsRepo := newFakeRequestsRepository()
	requestsRepo.requests = append(
		requestsRepo.requests,
		spotify_song_request.SpotifySongRequest{
			ID:        uuid.New(),
			ChannelID: testChannelID,
			Status:    spotify_song_request.StatusQueued,
		},
	)

	service := newService(
		&fakeSettingsRepository{settings: settings},
		&fakeIntegrationsRepository{
			integration: spotifyIntegration("user-modify-playback-state"),
		},
		requestsRepo,
	)

	swapHTTPClient(t, func(req *http.Request) (*http.Response, error) {
		return fixtureResponse(req, http.StatusOK, searchResponseBody), nil
	})

	_, err := service.CreateRequest(context.Background(), testChannelID, "", "viewer", "Viewer", "chat", "song")
	if !errors.Is(err, ErrMaxRequestsExceeded) {
		t.Fatalf("CreateRequest() error = %v, want %v", err, ErrMaxRequestsExceeded)
	}
}

func TestService_CreateRequest_rejects_user_max_requests(t *testing.T) {
	settings := spotifySettings()
	settings.UserMaxRequests = 1

	requestsRepo := newFakeRequestsRepository()
	requestsRepo.requests = append(
		requestsRepo.requests,
		spotify_song_request.SpotifySongRequest{
			ID:            uuid.New(),
			ChannelID:     testChannelID,
			RequesterName: "viewer",
			Status:        spotify_song_request.StatusQueued,
		},
	)

	service := newService(
		&fakeSettingsRepository{settings: settings},
		&fakeIntegrationsRepository{
			integration: spotifyIntegration("user-modify-playback-state"),
		},
		requestsRepo,
	)

	swapHTTPClient(t, func(req *http.Request) (*http.Response, error) {
		return fixtureResponse(req, http.StatusOK, searchResponseBody), nil
	})

	_, err := service.CreateRequest(context.Background(), testChannelID, "", "viewer", "Viewer", "chat", "song")
	if !errors.Is(err, ErrUserMaxRequestsExceeded) {
		t.Fatalf("CreateRequest() error = %v, want %v", err, ErrUserMaxRequestsExceeded)
	}
}

func TestService_CreateRequest_rejects_duration(t *testing.T) {
	settings := spotifySettings()
	settings.SongMaxLength = 2 // minutes; track fixture is 3 minutes

	service := newService(
		&fakeSettingsRepository{settings: settings},
		&fakeIntegrationsRepository{
			integration: spotifyIntegration("user-modify-playback-state"),
		},
		newFakeRequestsRepository(),
	)

	swapHTTPClient(t, func(req *http.Request) (*http.Response, error) {
		return fixtureResponse(req, http.StatusOK, searchResponseBody), nil
	})

	_, err := service.CreateRequest(context.Background(), testChannelID, "", "viewer", "Viewer", "chat", "song")
	if !errors.Is(err, ErrDurationNotAllowed) {
		t.Fatalf("CreateRequest() error = %v, want %v", err, ErrDurationNotAllowed)
	}
}

func TestService_CreateRequest_allows_duration_within_minute_limit(t *testing.T) {
	settings := spotifySettings()
	settings.SongMaxLength = 10 // minutes; track fixture is 3 minutes

	service := newService(
		&fakeSettingsRepository{settings: settings},
		&fakeIntegrationsRepository{
			integration: spotifyIntegration("user-modify-playback-state"),
		},
		newFakeRequestsRepository(),
	)

	swapHTTPClient(t, func(req *http.Request) (*http.Response, error) {
		switch req.URL.Path {
		case "/v1/search":
			return fixtureResponse(req, http.StatusOK, searchResponseBody), nil
		case "/v1/me/player/devices":
			return fixtureResponse(req, http.StatusOK, devicesResponseBody), nil
		case "/v1/me/player/queue":
			return fixtureResponse(req, http.StatusNoContent, ""), nil
		default:
			t.Fatalf("unexpected request path %s", req.URL.Path)
			return nil, nil
		}
	})

	_, err := service.CreateRequest(context.Background(), testChannelID, "", "viewer", "Viewer", "chat", "song")
	if err != nil {
		t.Fatalf("CreateRequest() error = %v, want nil for 3-minute track with 10-minute limit", err)
	}
}

func TestService_CreateRequest_rejects_unknown_track(t *testing.T) {
	service := newService(
		&fakeSettingsRepository{settings: spotifySettings()},
		&fakeIntegrationsRepository{
			integration: spotifyIntegration("user-modify-playback-state"),
		},
		newFakeRequestsRepository(),
	)

	swapHTTPClient(t, func(req *http.Request) (*http.Response, error) {
		return fixtureResponse(req, http.StatusOK, `{"tracks":{"items":[]}}`), nil
	})

	_, err := service.CreateRequest(context.Background(), testChannelID, "", "viewer", "Viewer", "chat", "song")
	if !errors.Is(err, ErrTrackNotFound) {
		t.Fatalf("CreateRequest() error = %v, want %v", err, ErrTrackNotFound)
	}
}

func TestService_CancelRequest_cancels_latest_requester_request(t *testing.T) {
	settings := spotifySettings()

	older := spotify_song_request.SpotifySongRequest{
		ID:            uuid.New(),
		ChannelID:     testChannelID,
		RequesterName: "viewer",
		Title:         "older",
		Status:        spotify_song_request.StatusQueued,
	}
	latest := spotify_song_request.SpotifySongRequest{
		ID:            uuid.New(),
		ChannelID:     testChannelID,
		RequesterName: "viewer",
		Title:         "latest",
		Status:        spotify_song_request.StatusQueued,
	}

	requestsRepo := newFakeRequestsRepository()
	requestsRepo.requests = append(requestsRepo.requests, older, latest)

	service := newService(
		&fakeSettingsRepository{settings: settings},
		&fakeIntegrationsRepository{integration: spotifyIntegration()},
		requestsRepo,
	)

	cancelled, err := service.CancelRequest(context.Background(), testChannelID, "viewer")
	if err != nil {
		t.Fatalf("CancelRequest() error = %v", err)
	}
	if cancelled.ID != latest.ID {
		t.Fatalf("cancelled request = %v, want latest %v", cancelled.ID, latest.ID)
	}
	if cancelled.Status != spotify_song_request.StatusCancelledPendingSkip {
		t.Fatalf("status = %q, want cancelled_pending_skip", cancelled.Status)
	}
	if len(requestsRepo.cancelledIDs) != 1 || requestsRepo.cancelledIDs[0] != latest.ID.String() {
		t.Fatalf("cancelled ids = %v, want [%s]", requestsRepo.cancelledIDs, latest.ID)
	}
}

func TestService_CancelRequest_no_active_requests(t *testing.T) {
	service := newService(
		&fakeSettingsRepository{settings: spotifySettings()},
		&fakeIntegrationsRepository{integration: spotifyIntegration()},
		newFakeRequestsRepository(),
	)

	_, err := service.CancelRequest(context.Background(), testChannelID, "viewer")
	if !errors.Is(err, spotify.ErrTrackNotFound) {
		t.Fatalf("CancelRequest() error = %v, want %v", err, spotify.ErrTrackNotFound)
	}
}

const nothingPlayingBody = `{"currently_playing_type":"","is_playing":false,"device":{"id":"device-1","name":"Desktop","type":"Computer","is_active":true},"item":null}`

func TestService_SkipRequest_defers_when_track_not_playing(t *testing.T) {
	requestID := uuid.New()
	requestsRepo := newFakeRequestsRepository()
	requestsRepo.requests = append(
		requestsRepo.requests,
		spotify_song_request.SpotifySongRequest{
			ID:        requestID,
			ChannelID: testChannelID,
			TrackURI:  "spotify:track:track-1",
			Status:    spotify_song_request.StatusQueued,
		},
	)

	swapHTTPClient(t, func(req *http.Request) (*http.Response, error) {
		if req.URL.Path != "/v1/me/player/currently-playing" {
			t.Fatalf("unexpected request path %s", req.URL.Path)
		}
		return fixtureResponse(req, http.StatusOK, nothingPlayingBody), nil
	})

	service := newService(
		&fakeSettingsRepository{settings: spotifySettings()},
		&fakeIntegrationsRepository{integration: spotifyIntegration()},
		requestsRepo,
	)

	if err := service.SkipRequest(context.Background(), testChannelID, requestID.String()); err != nil {
		t.Fatalf("SkipRequest() error = %v", err)
	}
	if len(requestsRepo.cancelledIDs) != 1 || requestsRepo.cancelledIDs[0] != requestID.String() {
		t.Fatalf("cancelled ids = %v, want [%s] (deferred skip)", requestsRepo.cancelledIDs, requestID)
	}
}

func TestService_SkipRequest_skips_immediately_when_track_is_playing(t *testing.T) {
	requestID := uuid.New()
	requestsRepo := newFakeRequestsRepository()
	requestsRepo.requests = append(
		requestsRepo.requests,
		spotify_song_request.SpotifySongRequest{
			ID:        requestID,
			ChannelID: testChannelID,
			TrackURI:  "spotify:track:track-1",
			Status:    spotify_song_request.StatusPlaying,
		},
	)

	var skippedDeviceID string
	swapHTTPClient(t, func(req *http.Request) (*http.Response, error) {
		switch req.URL.Path {
		case "/v1/me/player/currently-playing":
			body := `{"currently_playing_type":"track","is_playing":true,"progress_ms":1,"device":{"id":"device-1","name":"Desktop","type":"Computer","is_active":true},"item":{"id":"track-1","uri":"spotify:track:track-1","name":"Song","type":"track","artists":[{"name":"Artist"}],"album":{"name":"Album","images":[]},"duration_ms":180000}}`
			return fixtureResponse(req, http.StatusOK, body), nil
		case "/v1/me/player/next":
			skippedDeviceID = req.URL.Query().Get("device_id")
			return fixtureResponse(req, http.StatusNoContent, ""), nil
		default:
			t.Fatalf("unexpected request path %s", req.URL.Path)
			return nil, nil
		}
	})

	service := newService(
		&fakeSettingsRepository{settings: spotifySettings()},
		&fakeIntegrationsRepository{integration: spotifyIntegration()},
		requestsRepo,
	)

	if err := service.SkipRequest(context.Background(), testChannelID, requestID.String()); err != nil {
		t.Fatalf("SkipRequest() error = %v", err)
	}
	if skippedDeviceID != "device-1" {
		t.Fatalf("skipped device = %q, want device-1", skippedDeviceID)
	}
	if got := requestsRepo.statusByID[requestID.String()]; got != spotify_song_request.StatusSkippedByTwir {
		t.Fatalf("status = %q, want skipped_by_twir", got)
	}
}

func TestService_SkipRequest_rejects_foreign_channel(t *testing.T) {
	requestID := uuid.New()
	requestsRepo := newFakeRequestsRepository()
	requestsRepo.requests = append(
		requestsRepo.requests,
		spotify_song_request.SpotifySongRequest{
			ID:        requestID,
			ChannelID: "other-channel",
			Status:    spotify_song_request.StatusQueued,
		},
	)

	service := newService(
		&fakeSettingsRepository{settings: spotifySettings()},
		&fakeIntegrationsRepository{integration: spotifyIntegration()},
		requestsRepo,
	)

	err := service.SkipRequest(context.Background(), testChannelID, requestID.String())
	if !errors.Is(err, spotify_song_requests_repository.ErrNotFound) {
		t.Fatalf("SkipRequest() error = %v, want %v", err, spotify_song_requests_repository.ErrNotFound)
	}
}

func TestService_CancelRequestByID_marks_pending_skip(t *testing.T) {
	requestID := uuid.New()
	requestsRepo := newFakeRequestsRepository()
	requestsRepo.requests = append(
		requestsRepo.requests,
		spotify_song_request.SpotifySongRequest{
			ID:        requestID,
			ChannelID: testChannelID,
			TrackURI:  "spotify:track:track-1",
			Status:    spotify_song_request.StatusQueued,
		},
	)

	swapHTTPClient(t, func(req *http.Request) (*http.Response, error) {
		if req.URL.Path != "/v1/me/player/currently-playing" {
			t.Fatalf("unexpected request path %s", req.URL.Path)
		}
		return fixtureResponse(req, http.StatusOK, nothingPlayingBody), nil
	})

	service := newService(
		&fakeSettingsRepository{settings: spotifySettings()},
		&fakeIntegrationsRepository{integration: spotifyIntegration()},
		requestsRepo,
	)

	if err := service.CancelRequestByID(context.Background(), testChannelID, requestID.String()); err != nil {
		t.Fatalf("CancelRequestByID() error = %v", err)
	}
	if len(requestsRepo.cancelledIDs) != 1 || requestsRepo.cancelledIDs[0] != requestID.String() {
		t.Fatalf("cancelled ids = %v, want [%s]", requestsRepo.cancelledIDs, requestID)
	}
}

func TestReconciler_removes_request_missing_from_spotify_queue(t *testing.T) {
	newReconcilerFixture := func(requestID uuid.UUID, missingSince time.Time) (*Reconciler, *fakeRequestsRepository) {
		requestsRepo := newFakeRequestsRepository()
		requestsRepo.requests = append(
			requestsRepo.requests,
			spotify_song_request.SpotifySongRequest{
				ID:        requestID,
				ChannelID: testChannelID,
				TrackURI:  "spotify:track:track-1",
				Status:    spotify_song_request.StatusQueued,
			},
		)

		swapHTTPClient(t, func(req *http.Request) (*http.Response, error) {
			switch req.URL.Path {
			case "/v1/me/player/currently-playing":
				body := `{"currently_playing_type":"track","is_playing":true,"progress_ms":1,"device":{"id":"device-1","name":"Desktop","type":"Computer","is_active":true},"item":{"id":"other-track","uri":"spotify:track:other-track","name":"Other","type":"track","artists":[{"name":"Artist"}],"album":{"name":"Album","images":[]},"duration_ms":180000}}`
				return fixtureResponse(req, http.StatusOK, body), nil
			case "/v1/me/player/queue":
				return fixtureResponse(req, http.StatusOK, `{"queue":[]}`), nil
			default:
				t.Fatalf("unexpected request path %s", req.URL.Path)
				return nil, nil
			}
		})

		service := newService(
			&fakeSettingsRepository{settings: spotifySettings()},
			&fakeIntegrationsRepository{
				integration: spotifyIntegration("user-read-playback-state", "user-modify-playback-state"),
			},
			requestsRepo,
		)
		reconciler := &Reconciler{
			service: service,
			missing: map[string]time.Time{requestID.String(): missingSince},
		}

		return reconciler, requestsRepo
	}

	t.Run("marks missing on first observation without removing", func(t *testing.T) {
		requestID := uuid.New()
		reconciler, requestsRepo := newReconcilerFixture(requestID, time.Now())

		if hadError := reconciler.reconcileChannel(context.Background(), testChannelID); hadError {
			t.Fatal("reconcileChannel() reported error")
		}
		if _, ok := requestsRepo.statusByID[requestID.String()]; ok {
			t.Fatalf("status changed to %q on first missing observation", requestsRepo.statusByID[requestID.String()])
		}
		if _, ok := reconciler.missing[requestID.String()]; !ok {
			t.Fatal("missing timestamp was not recorded")
		}
	})

	t.Run("removes after missing threshold", func(t *testing.T) {
		requestID := uuid.New()
		reconciler, requestsRepo := newReconcilerFixture(requestID, time.Now().Add(-time.Minute))

		if hadError := reconciler.reconcileChannel(context.Background(), testChannelID); hadError {
			t.Fatal("reconcileChannel() reported error")
		}
		if got := requestsRepo.statusByID[requestID.String()]; got != spotify_song_request.StatusRemovedOrReconciled {
			t.Fatalf("status = %q, want removed_or_reconciled", got)
		}
	})
}

func TestReconciler_skips_cancelled_request_playing_now(t *testing.T) {
	requestID := uuid.New()
	requestsRepo := newFakeRequestsRepository()
	requestsRepo.requests = append(
		requestsRepo.requests,
		spotify_song_request.SpotifySongRequest{
			ID:        requestID,
			ChannelID: testChannelID,
			TrackURI:  "spotify:track:track-1",
			Status:    spotify_song_request.StatusCancelledPendingSkip,
		},
	)

	var skippedDeviceID string
	swapHTTPClient(t, func(req *http.Request) (*http.Response, error) {
		switch req.URL.Path {
		case "/v1/me/player/currently-playing":
			body := `{"currently_playing_type":"track","is_playing":true,"progress_ms":1,"device":{"id":"device-1","name":"Desktop","type":"Computer","is_active":true},"item":{"id":"track-1","uri":"spotify:track:track-1","name":"Song","type":"track","artists":[{"name":"Artist"}],"album":{"name":"Album","images":[]},"duration_ms":180000}}`
			return fixtureResponse(req, http.StatusOK, body), nil
		case "/v1/me/player/queue":
			return fixtureResponse(req, http.StatusOK, `{"queue":[]}`), nil
		case "/v1/me/player/next":
			skippedDeviceID = req.URL.Query().Get("device_id")
			return fixtureResponse(req, http.StatusNoContent, ""), nil
		default:
			t.Fatalf("unexpected request path %s", req.URL.Path)
			return nil, nil
		}
	})

	service := newService(
		&fakeSettingsRepository{settings: spotifySettings()},
		&fakeIntegrationsRepository{
			integration: spotifyIntegration("user-read-playback-state", "user-modify-playback-state"),
		},
		requestsRepo,
	)
	reconciler := &Reconciler{service: service, missing: map[string]time.Time{}}

	if hadError := reconciler.reconcileChannel(context.Background(), testChannelID); hadError {
		t.Fatal("reconcileChannel() reported error")
	}
	if skippedDeviceID != "device-1" {
		t.Fatalf("skipped device = %q, want device-1", skippedDeviceID)
	}
	if got := requestsRepo.statusByID[requestID.String()]; got != spotify_song_request.StatusSkippedByTwir {
		t.Fatalf("status = %q, want skipped_by_twir", got)
	}
}

func TestReconciler_marks_queued_request_as_playing(t *testing.T) {
	requestID := uuid.New()
	requestsRepo := newFakeRequestsRepository()
	requestsRepo.requests = append(
		requestsRepo.requests,
		spotify_song_request.SpotifySongRequest{
			ID:        requestID,
			ChannelID: testChannelID,
			TrackURI:  "spotify:track:track-1",
			Status:    spotify_song_request.StatusQueued,
		},
	)

	swapHTTPClient(t, func(req *http.Request) (*http.Response, error) {
		switch req.URL.Path {
		case "/v1/me/player/currently-playing":
			body := `{"currently_playing_type":"track","is_playing":true,"progress_ms":1,"device":{"id":"device-1","name":"Desktop","type":"Computer","is_active":true},"item":{"id":"track-1","uri":"spotify:track:track-1","name":"Song","type":"track","artists":[{"name":"Artist"}],"album":{"name":"Album","images":[]},"duration_ms":180000}}`
			return fixtureResponse(req, http.StatusOK, body), nil
		case "/v1/me/player/queue":
			return fixtureResponse(req, http.StatusOK, `{"queue":[]}`), nil
		default:
			t.Fatalf("unexpected request path %s", req.URL.Path)
			return nil, nil
		}
	})

	service := newService(
		&fakeSettingsRepository{settings: spotifySettings()},
		&fakeIntegrationsRepository{
			integration: spotifyIntegration("user-read-playback-state", "user-modify-playback-state"),
		},
		requestsRepo,
	)
	reconciler := &Reconciler{service: service, missing: map[string]time.Time{}}

	if hadError := reconciler.reconcileChannel(context.Background(), testChannelID); hadError {
		t.Fatal("reconcileChannel() reported error")
	}
	if got := requestsRepo.statusByID[requestID.String()]; got != spotify_song_request.StatusPlaying {
		t.Fatalf("status = %q, want playing", got)
	}
}

func TestReconciler_marks_played_when_another_track_is_playing(t *testing.T) {
	requestID := uuid.New()
	requestsRepo := newFakeRequestsRepository()
	requestsRepo.requests = append(
		requestsRepo.requests,
		spotify_song_request.SpotifySongRequest{
			ID:        requestID,
			ChannelID: testChannelID,
			TrackURI:  "spotify:track:track-1",
			Status:    spotify_song_request.StatusPlaying,
		},
	)

	swapHTTPClient(t, func(req *http.Request) (*http.Response, error) {
		switch req.URL.Path {
		case "/v1/me/player/currently-playing":
			body := `{"currently_playing_type":"track","is_playing":true,"progress_ms":1,"device":{"id":"device-1","name":"Desktop","type":"Computer","is_active":true},"item":{"id":"other-track","uri":"spotify:track:other-track","name":"Other","type":"track","artists":[{"name":"Artist"}],"album":{"name":"Album","images":[]},"duration_ms":180000}}`
			return fixtureResponse(req, http.StatusOK, body), nil
		case "/v1/me/player/queue":
			return fixtureResponse(req, http.StatusOK, `{"queue":[]}`), nil
		default:
			t.Fatalf("unexpected request path %s", req.URL.Path)
			return nil, nil
		}
	})

	service := newService(
		&fakeSettingsRepository{settings: spotifySettings()},
		&fakeIntegrationsRepository{
			integration: spotifyIntegration("user-read-playback-state", "user-modify-playback-state"),
		},
		requestsRepo,
	)
	reconciler := &Reconciler{service: service, missing: map[string]time.Time{}}

	if hadError := reconciler.reconcileChannel(context.Background(), testChannelID); hadError {
		t.Fatal("reconcileChannel() reported error")
	}
	if got := requestsRepo.statusByID[requestID.String()]; got != spotify_song_request.StatusPlayed {
		t.Fatalf("status = %q, want played", got)
	}
}

func TestReconciler_marks_played_when_playback_stopped(t *testing.T) {
	requestID := uuid.New()
	requestsRepo := newFakeRequestsRepository()
	requestsRepo.requests = append(
		requestsRepo.requests,
		spotify_song_request.SpotifySongRequest{
			ID:        requestID,
			ChannelID: testChannelID,
			TrackURI:  "spotify:track:track-1",
			Status:    spotify_song_request.StatusPlaying,
		},
	)

	swapHTTPClient(t, func(req *http.Request) (*http.Response, error) {
		switch req.URL.Path {
		case "/v1/me/player/currently-playing":
			return fixtureResponse(req, http.StatusOK, `{"currently_playing_type":"","is_playing":false,"device":{"id":"device-1","name":"Desktop","type":"Computer","is_active":true},"item":null}`), nil
		case "/v1/me/player/queue":
			return fixtureResponse(req, http.StatusOK, `{"queue":[]}`), nil
		default:
			t.Fatalf("unexpected request path %s", req.URL.Path)
			return nil, nil
		}
	})

	service := newService(
		&fakeSettingsRepository{settings: spotifySettings()},
		&fakeIntegrationsRepository{
			integration: spotifyIntegration("user-read-playback-state", "user-modify-playback-state"),
		},
		requestsRepo,
	)
	reconciler := &Reconciler{service: service, missing: map[string]time.Time{}}

	if hadError := reconciler.reconcileChannel(context.Background(), testChannelID); hadError {
		t.Fatal("reconcileChannel() reported error")
	}
	if got := requestsRepo.statusByID[requestID.String()]; got != spotify_song_request.StatusPlayed {
		t.Fatalf("status = %q, want played", got)
	}
}

func TestReconciler_marks_unknown_without_playback_scopes(t *testing.T) {
	requestID := uuid.New()
	requestsRepo := newFakeRequestsRepository()
	requestsRepo.requests = append(
		requestsRepo.requests,
		spotify_song_request.SpotifySongRequest{
			ID:        requestID,
			ChannelID: testChannelID,
			TrackURI:  "spotify:track:track-1",
			Status:    spotify_song_request.StatusQueued,
		},
	)

	service := newService(
		&fakeSettingsRepository{settings: spotifySettings()},
		&fakeIntegrationsRepository{
			integration: spotifyIntegration("user-read-playback-state"),
		},
		requestsRepo,
	)
	reconciler := &Reconciler{service: service, missing: map[string]time.Time{}}

	if hadError := reconciler.reconcileChannel(context.Background(), testChannelID); hadError {
		t.Fatal("reconcileChannel() reported error")
	}
	if got := requestsRepo.statusByID[requestID.String()]; got != spotify_song_request.StatusUnknown {
		t.Fatalf("status = %q, want unknown", got)
	}
}
