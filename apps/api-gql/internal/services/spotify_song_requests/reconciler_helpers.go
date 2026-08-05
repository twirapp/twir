package spotify_song_requests

import (
	"context"
	"log/slog"
	"slices"
	"time"

	spotify_song_request "github.com/twirapp/twir/libs/entities/spotify_song_request"
	"github.com/twirapp/twir/libs/logger"
)

func (r *Reconciler) hasPlaybackScopes(scopes []string) bool {
	return slices.Contains(scopes, "user-read-playback-state") && slices.Contains(scopes, "user-modify-playback-state")
}

func (r *Reconciler) anyRequestPresent(requests []spotify_song_request.SpotifySongRequest, queueURIs map[string]struct{}, currentURI string) bool {
	for _, request := range requests {
		if request.TrackURI == currentURI {
			return true
		}
		if _, ok := queueURIs[request.TrackURI]; ok {
			return true
		}
	}
	return false
}

func (r *Reconciler) anyQueuedRequestPresent(requests []spotify_song_request.SpotifySongRequest, queueURIs map[string]struct{}, currentURI string) bool {
	for _, request := range requests {
		if request.Status == spotify_song_request.StatusQueued && (request.TrackURI == currentURI || hasTrackURI(queueURIs, request.TrackURI)) {
			return true
		}
	}
	return false
}

func hasTrackURI(queueURIs map[string]struct{}, trackURI string) bool {
	_, ok := queueURIs[trackURI]
	return ok
}

func (r *Reconciler) resolveSkipDevice(
	ctx context.Context,
	channelID string,
	currentDeviceID string,
) (string, error) {
	if currentDeviceID != "" {
		return currentDeviceID, nil
	}

	return r.service.SelectAndCacheDevice(ctx, channelID)
}

func (r *Reconciler) transitionStatus(
	ctx context.Context,
	channelID string,
	request spotify_song_request.SpotifySongRequest,
	status spotify_song_request.Status,
) bool {
	if request.Status == status {
		return false
	}

	if err := r.service.spotifySongRequestsRepository.UpdateStatus(ctx, request.ID.String(), status); err != nil {
		r.service.logger.ErrorContext(ctx, "failed to update spotify song request status", logger.Error(err), slog.String("channel_id", channelID), slog.String("request_id", request.ID.String()), slog.String("from", request.Status.String()), slog.String("to", status.String()))
		return true
	}

	r.service.publishQueueChanged(ctx, channelID)

	r.service.logger.InfoContext(ctx, "updated spotify song request status", slog.String("channel_id", channelID), slog.String("request_id", request.ID.String()), slog.String("from", request.Status.String()), slog.String("to", status.String()))
	return false
}

func (r *Reconciler) markRequestsUnknown(
	ctx context.Context,
	channelID string,
	requests []spotify_song_request.SpotifySongRequest,
) {
	for _, request := range requests {
		if request.Status == spotify_song_request.StatusUnknown {
			continue
		}
		if r.transitionStatus(ctx, channelID, request, spotify_song_request.StatusUnknown) {
			return
		}
	}
}

func (r *Reconciler) clearMissingSince(channelID string) {
	r.mu.Lock()
	delete(r.missing, channelID)
	r.mu.Unlock()
}

func (r *Reconciler) markMissingSince(channelID string) {
	r.mu.Lock()
	if _, ok := r.missing[channelID]; !ok {
		r.missing[channelID] = time.Now()
	}
	r.mu.Unlock()
}

func (r *Reconciler) shouldRemoveMissing(channelID string) bool {
	r.mu.Lock()
	startedAt, ok := r.missing[channelID]
	r.mu.Unlock()
	if !ok {
		return false
	}

	return time.Since(startedAt) >= reconcilerMissingThreshold
}
