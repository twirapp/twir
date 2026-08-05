package spotify_song_requests

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/twirapp/twir/libs/baseapp/lifecycle"
	spotify_song_request "github.com/twirapp/twir/libs/entities/spotify_song_request"
	"github.com/twirapp/twir/libs/integrations/spotify"
	"github.com/twirapp/twir/libs/logger"
)

const (
	reconcilerTickInterval     = 5 * time.Second
	reconcilerIdleInterval     = 60 * time.Second
	reconcilerMissingThreshold = 15 * time.Second
)

type Reconciler struct {
	service *Service
	ctx     context.Context
	cancel  context.CancelFunc
	wg      sync.WaitGroup
	mu      sync.Mutex
	missing map[string]time.Time
}

func NewReconciler(lc *lifecycle.Lifecycle, service *Service) *Reconciler {
	workerCtx, cancel := context.WithCancel(context.Background())
	r := &Reconciler{
		service: service,
		ctx:     workerCtx,
		cancel:  cancel,
		missing: map[string]time.Time{},
	}

	lc.Append(lifecycle.Hook{
		OnStart: func(context.Context) error {
			r.wg.Add(1)
			go func() {
				defer r.wg.Done()
				r.run()
			}()
			return nil
		},
		OnStop: func(context.Context) error {
			r.cancel()
			r.wg.Wait()
			return nil
		},
	})

	return r
}

func (r *Reconciler) run() {
	interval := reconcilerTickInterval
	for {
		timer := time.NewTimer(interval)
		select {
		case <-r.ctx.Done():
			timer.Stop()
			return
		case <-timer.C:
		}

		hadError, activeChannels := r.tick(r.ctx)
		interval = reconcilerTickInterval
		if !hadError && activeChannels == 0 {
			interval = reconcilerIdleInterval
		}
	}
}

func (r *Reconciler) tick(ctx context.Context) (bool, int) {
	type activeChannelsRepository interface {
		GetActiveChannels(ctx context.Context) ([]string, error)
	}

	repo := r.service.spotifySongRequestsRepository.(activeChannelsRepository)
	channels, err := repo.GetActiveChannels(ctx)
	if err != nil {
		r.service.logger.ErrorContext(ctx, "failed to load active spotify song request channels", logger.Error(err))
		return true, 0
	}
	if len(channels) == 0 {
		return false, 0
	}

	hadError := false
	for _, channelID := range channels {
		if r.reconcileChannel(ctx, channelID) {
			hadError = true
		}
	}

	return hadError, len(channels)
}

type playbackSnapshot struct {
	currentURI      string
	currentDeviceID string
	queueURIs       map[string]struct{}
}

func (r *Reconciler) reconcileChannel(ctx context.Context, channelID string) bool {
	activeRequests, err := r.service.spotifySongRequestsRepository.GetActiveByChannel(ctx, channelID)
	if err != nil {
		r.service.logger.ErrorContext(ctx, "failed to load active spotify song requests", logger.Error(err), slog.String("channel_id", channelID))
		return true
	}

	client, integration, err := r.service.loadSpotifyClient(ctx, channelID)
	if err != nil {
		if errors.Is(err, spotify.ErrNotConnected) {
			r.markRequestsUnknown(ctx, channelID, activeRequests)
			return false
		}
		r.service.logger.ErrorContext(ctx, "failed to load spotify client", logger.Error(err), slog.String("channel_id", channelID))
		return true
	}
	if !r.hasPlaybackScopes(integration.Scopes) {
		r.markRequestsUnknown(ctx, channelID, activeRequests)
		return false
	}

	snapshot, hadError, skip := r.loadPlaybackSnapshot(ctx, channelID, client, activeRequests)
	if skip {
		return hadError
	}

	if r.reconcileMissingRequests(ctx, channelID, activeRequests, snapshot) {
		return true
	}

	playingMatched := false
	for _, request := range activeRequests {
		if request.Status == spotify_song_request.StatusCancelledPendingSkip && snapshot.currentURI != "" && request.TrackURI == snapshot.currentURI {
			if r.executeDeferredSkip(ctx, channelID, client, request, snapshot.currentDeviceID) {
				return true
			}
			continue
		}

		if request.TrackURI == snapshot.currentURI && !playingMatched {
			playingMatched = true
			if request.Status == spotify_song_request.StatusQueued && r.transitionStatus(ctx, channelID, request, spotify_song_request.StatusPlaying) {
				return true
			}
			continue
		}

		if _, present := snapshot.queueURIs[request.TrackURI]; present {
			if request.Status != spotify_song_request.StatusQueued && request.Status != spotify_song_request.StatusCancelledPendingSkip && r.transitionStatus(ctx, channelID, request, spotify_song_request.StatusQueued) {
				return true
			}
			continue
		}

		if request.Status == spotify_song_request.StatusPlaying && request.TrackURI != snapshot.currentURI && r.transitionStatus(ctx, channelID, request, spotify_song_request.StatusPlayed) {
			return true
		}
	}

	if !r.anyQueuedRequestPresent(activeRequests, snapshot.queueURIs, snapshot.currentURI) {
		r.markMissingSince(channelID)
	}

	return false
}

func (r *Reconciler) loadPlaybackSnapshot(
	ctx context.Context,
	channelID string,
	client spotifyClient,
	activeRequests []spotify_song_request.SpotifySongRequest,
) (snapshot *playbackSnapshot, hadError bool, skip bool) {
	currentlyPlaying, err := client.GetCurrentlyPlaying(ctx)
	if err != nil {
		return nil, r.handleCurrentlyPlayingError(ctx, channelID, activeRequests, err), true
	}

	queue, err := client.GetQueue(ctx)
	if err != nil {
		if errors.Is(err, spotify.ErrRateLimited) {
			r.service.logger.ErrorContext(ctx, "spotify queue rate limited", logger.Error(err), slog.String("channel_id", channelID))
		} else {
			r.service.logger.ErrorContext(ctx, "failed to get spotify queue", logger.Error(err), slog.String("channel_id", channelID))
		}
		return nil, true, true
	}

	snapshot = &playbackSnapshot{
		queueURIs: make(map[string]struct{}, len(queue)+1),
	}
	if currentlyPlaying != nil {
		snapshot.currentDeviceID = currentlyPlaying.Device.ID
		r.cachePlayerDevice(ctx, channelID, currentlyPlaying.Device)
		if currentlyPlaying.Item != nil {
			snapshot.currentURI = currentlyPlaying.Item.URI
			snapshot.queueURIs[snapshot.currentURI] = struct{}{}
		}
	}
	for _, track := range queue {
		snapshot.queueURIs[track.URI] = struct{}{}
	}

	return snapshot, false, false
}

func (r *Reconciler) handleCurrentlyPlayingError(
	ctx context.Context,
	channelID string,
	activeRequests []spotify_song_request.SpotifySongRequest,
	err error,
) bool {
	switch {
	case errors.Is(err, spotify.ErrInsufficientScope):
		r.markRequestsUnknown(ctx, channelID, activeRequests)
		return false
	case errors.Is(err, spotify.ErrNoActiveDevice):
		return false
	case errors.Is(err, spotify.ErrRateLimited):
		r.service.logger.ErrorContext(ctx, "spotify currently-playing rate limited", logger.Error(err), slog.String("channel_id", channelID))
		return true
	default:
		r.service.logger.ErrorContext(ctx, "failed to get spotify currently playing track", logger.Error(err), slog.String("channel_id", channelID))
		return true
	}
}

func (r *Reconciler) cachePlayerDevice(ctx context.Context, channelID string, device spotify.Device) {
	if device.ID == "" {
		return
	}
	if err := r.service.cacheDevice(ctx, channelID, device); err != nil {
		r.service.logger.ErrorContext(
			ctx,
			"failed to cache spotify device",
			logger.Error(err),
			slog.String("channel_id", channelID),
		)
	}
}

func (r *Reconciler) reconcileMissingRequests(
	ctx context.Context,
	channelID string,
	activeRequests []spotify_song_request.SpotifySongRequest,
	snapshot *playbackSnapshot,
) bool {
	if r.anyRequestPresent(activeRequests, snapshot.queueURIs, snapshot.currentURI) {
		r.clearMissingSince(channelID)
		return false
	}

	if !r.shouldRemoveMissing(channelID) {
		return false
	}

	for _, request := range activeRequests {
		if request.Status == spotify_song_request.StatusQueued && r.transitionStatus(ctx, channelID, request, spotify_song_request.StatusRemovedOrReconciled) {
			return true
		}
	}

	return false
}

func (r *Reconciler) executeDeferredSkip(
	ctx context.Context,
	channelID string,
	client spotifyClient,
	request spotify_song_request.SpotifySongRequest,
	currentDeviceID string,
) bool {
	deviceID, err := r.resolveSkipDevice(ctx, channelID, currentDeviceID)
	if err != nil {
		if errors.Is(err, spotify.ErrNoActiveDevice) {
			return r.transitionStatus(ctx, channelID, request, spotify_song_request.StatusUnknown)
		}
		r.service.logger.ErrorContext(ctx, "failed to resolve spotify device for skip", logger.Error(err), slog.String("channel_id", channelID), slog.String("request_id", request.ID.String()))
		return true
	}

	if err := client.SkipNext(ctx, deviceID); err != nil {
		return r.handleSkipError(ctx, channelID, request, err)
	}

	return r.transitionStatus(ctx, channelID, request, spotify_song_request.StatusSkippedByTwir)
}

func (r *Reconciler) handleSkipError(
	ctx context.Context,
	channelID string,
	request spotify_song_request.SpotifySongRequest,
	err error,
) bool {
	if errors.Is(err, spotify.ErrNoActiveDevice) {
		return r.transitionStatus(ctx, channelID, request, spotify_song_request.StatusUnknown)
	}
	if errors.Is(err, spotify.ErrRateLimited) {
		r.service.logger.ErrorContext(ctx, "spotify skip rate limited", logger.Error(err), slog.String("channel_id", channelID), slog.String("request_id", request.ID.String()))
		return true
	}
	r.service.logger.ErrorContext(ctx, "failed to skip spotify track", logger.Error(err), slog.String("channel_id", channelID), slog.String("request_id", request.ID.String()))
	return true
}
