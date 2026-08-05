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

	currentlyPlaying, err := client.GetCurrentlyPlaying(ctx)
	if err != nil {
		if errors.Is(err, spotify.ErrInsufficientScope) {
			r.markRequestsUnknown(ctx, channelID, activeRequests)
			return false
		}
		if errors.Is(err, spotify.ErrNoActiveDevice) {
			return false
		}
		if errors.Is(err, spotify.ErrRateLimited) {
			r.service.logger.ErrorContext(ctx, "spotify currently-playing rate limited", logger.Error(err), slog.String("channel_id", channelID))
			return true
		}
		r.service.logger.ErrorContext(ctx, "failed to get spotify currently playing track", logger.Error(err), slog.String("channel_id", channelID))
		return true
	}

	queue, err := client.GetQueue(ctx)
	if err != nil {
		if errors.Is(err, spotify.ErrRateLimited) {
			r.service.logger.ErrorContext(ctx, "spotify queue rate limited", logger.Error(err), slog.String("channel_id", channelID))
			return true
		}
		r.service.logger.ErrorContext(ctx, "failed to get spotify queue", logger.Error(err), slog.String("channel_id", channelID))
		return true
	}

	queueURIs := make(map[string]struct{}, len(queue)+1)
	currentURI := ""
	currentDeviceID := ""
	if currentlyPlaying != nil {
		currentDeviceID = currentlyPlaying.Device.ID
		if currentlyPlaying.Item != nil {
			currentURI = currentlyPlaying.Item.URI
			queueURIs[currentURI] = struct{}{}
		}
	}
	for _, track := range queue {
		queueURIs[track.URI] = struct{}{}
	}

	if r.anyRequestPresent(activeRequests, queueURIs, currentURI) {
		r.clearMissingSince(channelID)
	} else if r.shouldRemoveMissing(channelID) {
		for _, request := range activeRequests {
			if request.Status == spotify_song_request.StatusQueued && r.transitionStatus(ctx, channelID, request, spotify_song_request.StatusRemovedOrReconciled) {
				return true
			}
		}
	}

	playingMatched := false
	for _, request := range activeRequests {
		if request.Status == spotify_song_request.StatusCancelledPendingSkip && currentURI != "" && request.TrackURI == currentURI {
			deviceID, deviceErr := r.resolveSkipDevice(ctx, channelID, currentDeviceID)
			if deviceErr != nil {
				if errors.Is(deviceErr, spotify.ErrNoActiveDevice) {
					if r.transitionStatus(ctx, channelID, request, spotify_song_request.StatusUnknown) {
						return true
					}
					continue
				}
				r.service.logger.ErrorContext(ctx, "failed to resolve spotify device for skip", logger.Error(deviceErr), slog.String("channel_id", channelID), slog.String("request_id", request.ID.String()))
				return true
			}

			if err := client.SkipNext(ctx, deviceID); err != nil {
				if errors.Is(err, spotify.ErrNoActiveDevice) {
					if r.transitionStatus(ctx, channelID, request, spotify_song_request.StatusUnknown) {
						return true
					}
					continue
				}
				if errors.Is(err, spotify.ErrRateLimited) {
					r.service.logger.ErrorContext(ctx, "spotify skip rate limited", logger.Error(err), slog.String("channel_id", channelID), slog.String("request_id", request.ID.String()))
					return true
				}
				r.service.logger.ErrorContext(ctx, "failed to skip spotify track", logger.Error(err), slog.String("channel_id", channelID), slog.String("request_id", request.ID.String()))
				return true
			}

			if r.transitionStatus(ctx, channelID, request, spotify_song_request.StatusSkippedByTwir) {
				return true
			}
			continue
		}

		if request.TrackURI == currentURI && !playingMatched {
			playingMatched = true
			if request.Status == spotify_song_request.StatusQueued && r.transitionStatus(ctx, channelID, request, spotify_song_request.StatusPlaying) {
				return true
			}
			continue
		}

		if _, present := queueURIs[request.TrackURI]; present {
			if request.Status != spotify_song_request.StatusQueued && request.Status != spotify_song_request.StatusCancelledPendingSkip && r.transitionStatus(ctx, channelID, request, spotify_song_request.StatusQueued) {
				return true
			}
			continue
		}

		if request.Status == spotify_song_request.StatusPlaying && currentURI == "" && r.transitionStatus(ctx, channelID, request, spotify_song_request.StatusPlayed) {
			return true
		}
	}

	if !r.anyQueuedRequestPresent(activeRequests, queueURIs, currentURI) {
		r.markMissingSince(channelID)
	}

	return false
}
