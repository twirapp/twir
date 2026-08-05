package resolvers

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/google/uuid"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlerrors"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/libs/entities/song_request_mode"
	spotify_song_request "github.com/twirapp/twir/libs/entities/spotify_song_request"
	apperrors "github.com/twirapp/twir/libs/errors"
	"github.com/twirapp/twir/libs/integrations/spotify"
	songrequestssettingsrepository "github.com/twirapp/twir/libs/repositories/song_requests_settings"
)

const (
	defaultSpotifySearchLimit = 5
	maxSpotifySearchLimit     = 20
)

func (r *mutationResolver) spotifySongRequestSelectDevice(ctx context.Context, deviceID string) (bool, error) {
	dashboardID, err := r.deps.Sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return false, gqlerrors.HandleError(err)
	}

	if err := r.deps.SpotifySongRequestsService.SetSelectedDevice(ctx, dashboardID, deviceID); err != nil {
		return false, spotifyGraphQLError(err)
	}

	return true, nil
}

func (r *mutationResolver) spotifySongRequestRefreshDevice(ctx context.Context) (*gqlmodel.SpotifySongRequestDevice, error) {
	dashboardID, err := r.deps.Sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return nil, gqlerrors.HandleError(err)
	}

	device, err := r.spotifyDevice(ctx, dashboardID)
	if err != nil {
		return nil, spotifyGraphQLError(err)
	}

	return device, nil
}

func (r *mutationResolver) spotifySongRequestSkip(ctx context.Context, requestID uuid.UUID) (bool, error) {
	dashboardID, err := r.deps.Sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return false, gqlerrors.HandleError(err)
	}

	if err := r.deps.SpotifySongRequestsService.SkipRequest(ctx, dashboardID, requestID.String()); err != nil {
		return false, spotifyGraphQLError(err)
	}

	return true, nil
}

func (r *mutationResolver) spotifySongRequestCancel(ctx context.Context, requestID uuid.UUID) (bool, error) {
	dashboardID, err := r.deps.Sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return false, gqlerrors.HandleError(err)
	}

	if err := r.deps.SpotifySongRequestsService.CancelRequestByID(ctx, dashboardID, requestID.String()); err != nil {
		return false, spotifyGraphQLError(err)
	}

	return true, nil
}

func (r *queryResolver) spotifySongRequestsQueue(ctx context.Context) (*gqlmodel.SpotifySongRequestQueue, error) {
	dashboardID, err := r.deps.Sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return nil, gqlerrors.HandleError(err)
	}

	settings, err := r.deps.SongRequestsService.GetSettings(ctx, dashboardID)
	if err != nil {
		if errors.Is(err, songrequestssettingsrepository.ErrNotFound) {
			return emptySpotifyQueue(), nil
		}
		return nil, gqlerrors.HandleError(fmt.Errorf("failed to get song requests settings: %w", err))
	}
	if settings.Mode != song_request_mode.ModeSpotify {
		return emptySpotifyQueue(), nil
	}

	requests, err := r.deps.SpotifySongRequestsService.GetActiveQueue(ctx, dashboardID)
	if err != nil {
		return nil, spotifyGraphQLError(err)
	}

	result := &gqlmodel.SpotifySongRequestQueue{
		Requests: make([]gqlmodel.SpotifySongRequest, 0, len(requests)),
	}
	for _, request := range requests {
		result.Requests = append(result.Requests, mapSpotifySongRequest(request))
	}

	result.CurrentDevice, err = r.spotifyDevice(ctx, dashboardID)
	if err != nil {
		return nil, spotifyGraphQLError(err)
	}

	return result, nil
}

func (r *queryResolver) spotifySongRequestsSearch(ctx context.Context, query string, limit *int) ([]gqlmodel.SpotifySongRequestSearchResult, error) {
	dashboardID, err := r.deps.Sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return nil, gqlerrors.HandleError(err)
	}

	searchLimit := defaultSpotifySearchLimit
	if limit != nil && *limit > 0 {
		searchLimit = min(*limit, maxSpotifySearchLimit)
	}

	tracks, err := r.deps.SpotifySongRequestsService.SearchTracks(ctx, dashboardID, query, searchLimit)
	if err != nil {
		return nil, spotifyGraphQLError(err)
	}

	result := make([]gqlmodel.SpotifySongRequestSearchResult, 0, len(tracks))
	for _, track := range tracks {
		imageURL := track.ImageURL
		var imageURLPtr *string
		if imageURL != "" {
			imageURLPtr = &imageURL
		}
		result = append(result, gqlmodel.SpotifySongRequestSearchResult{
			ID:         track.ID,
			Title:      track.Name,
			Artist:     track.ArtistName,
			Album:      track.AlbumName,
			DurationMs: track.DurationMs,
			ImageURL:   imageURLPtr,
		})
	}

	return result, nil
}

func (r *songRequestsSettingsResolver) spotifyCapabilities(ctx context.Context, _ *gqlmodel.SongRequestsSettings) (*gqlmodel.SpotifySongRequestCapabilities, error) {
	dashboardID, err := r.deps.Sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return nil, gqlerrors.HandleError(err)
	}

	integration, err := r.deps.SpotifyIntegrationService.GetSpotifyData(ctx, dashboardID)
	if err != nil {
		return nil, gqlerrors.HandleError(fmt.Errorf("failed to get Spotify integration: %w", err))
	}

	capabilities := &gqlmodel.SpotifySongRequestCapabilities{}
	if integration == nil {
		return capabilities, nil
	}

	capabilities.Connected = true
	capabilities.HasPlaybackScope = slices.Contains(integration.Scopes, "user-read-playback-state") && slices.Contains(integration.Scopes, "user-modify-playback-state")
	capabilities.CanUseSpotify = capabilities.Connected && capabilities.HasPlaybackScope
	if !capabilities.CanUseSpotify {
		return capabilities, nil
	}

	device, err := r.spotifyDevice(ctx, dashboardID)
	if err != nil {
		return nil, spotifyGraphQLError(err)
	}
	capabilities.ActiveDevice = device
	capabilities.SelectedDevice = device

	return capabilities, nil
}

func (r *Resolver) spotifyDevice(ctx context.Context, dashboardID string) (*gqlmodel.SpotifySongRequestDevice, error) {
	deviceID, err := r.deps.SpotifySongRequestsService.SelectAndCacheDevice(ctx, dashboardID)
	if err != nil {
		if errors.Is(err, spotify.ErrNoActiveDevice) {
			return nil, nil
		}
		return nil, err
	}

	devices, err := r.deps.SpotifySongRequestsService.GetDevices(ctx, dashboardID)
	if err != nil {
		return nil, err
	}
	for _, device := range devices {
		if device.ID == deviceID {
			return &gqlmodel.SpotifySongRequestDevice{
				ID:       device.ID,
				Name:     device.Name,
				Type:     device.Type,
				IsActive: device.IsActive,
			}, nil
		}
	}

	return nil, nil
}

func mapSpotifySongRequest(request spotify_song_request.SpotifySongRequest) gqlmodel.SpotifySongRequest {
	status := gqlmodel.SpotifySongRequestStatusUnknown
	switch request.Status {
	case spotify_song_request.StatusQueued:
		status = gqlmodel.SpotifySongRequestStatusQueued
	case spotify_song_request.StatusPlaying:
		status = gqlmodel.SpotifySongRequestStatusPlaying
	case spotify_song_request.StatusPlayed:
		status = gqlmodel.SpotifySongRequestStatusPlayed
	case spotify_song_request.StatusSkippedByTwir:
		status = gqlmodel.SpotifySongRequestStatus("SKIPPED")
	case spotify_song_request.StatusCancelledPendingSkip:
		status = gqlmodel.SpotifySongRequestStatusCancelledPendingSkip
	}

	return gqlmodel.SpotifySongRequest{
		ID:                   request.ID,
		Title:                request.Title,
		Artist:               request.Artist,
		Album:                request.Album,
		DurationMs:           request.DurationMs,
		RequesterName:        request.RequesterName,
		RequesterDisplayName: request.RequesterDisplayName,
		Source:               request.Source,
		QueuePosition:        request.QueuePosition,
		Status:               status,
		CreatedAt:            request.CreatedAt,
	}
}

func emptySpotifyQueue() *gqlmodel.SpotifySongRequestQueue {
	return &gqlmodel.SpotifySongRequestQueue{Requests: []gqlmodel.SpotifySongRequest{}}
}

func spotifyGraphQLError(err error) error {
	switch {
	case errors.Is(err, spotify.ErrNotConnected):
		return gqlerrors.HandleError(apperrors.NewBadRequestError("Spotify is not connected. Connect Spotify first."))
	case errors.Is(err, spotify.ErrInsufficientScope):
		return gqlerrors.HandleError(apperrors.NewBadRequestError("Spotify playback permissions are missing. Reconnect Spotify."))
	case errors.Is(err, spotify.ErrNoActiveDevice):
		return gqlerrors.HandleError(apperrors.NewBadRequestError("No active Spotify device was found."))
	default:
		return gqlerrors.HandleError(fmt.Errorf("Spotify song request operation failed: %w", err))
	}
}
