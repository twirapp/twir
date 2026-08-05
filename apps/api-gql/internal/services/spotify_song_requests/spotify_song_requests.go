package spotify_song_requests

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/twirapp/kv"
	cfg "github.com/twirapp/twir/libs/config"
	spotify_song_request "github.com/twirapp/twir/libs/entities/spotify_song_request"
	"github.com/twirapp/twir/libs/integrations/spotify"
	"github.com/twirapp/twir/libs/logger"
	channelsintegrationsspotify "github.com/twirapp/twir/libs/repositories/channels_integrations_spotify"
	songrequestssettingsrepository "github.com/twirapp/twir/libs/repositories/song_requests_settings"
	spotify_song_requests_repository "github.com/twirapp/twir/libs/repositories/spotify_song_requests"
	"github.com/twirapp/twir/libs/wsrouter"
)

var (
	ErrNotSpotifyMode          = errors.New("spotify song requests: not spotify mode")
	ErrTrackNotFound           = errors.New("spotify song requests: track not found")
	ErrTrackAlreadyInQueue     = errors.New("spotify song requests: track already in queue")
	ErrMaxRequestsExceeded     = errors.New("spotify song requests: max requests exceeded")
	ErrUserMaxRequestsExceeded = errors.New("spotify song requests: user max requests exceeded")
	ErrDurationNotAllowed      = errors.New("spotify song requests: duration not allowed")
)

const spotifySongRequestDeviceCachePrefix = "spotify:songrequests:device:"

const spotifySongRequestDeviceCacheTTL = 5 * time.Minute

type Service struct {
	spotifySongRequestsRepository  spotify_song_requests_repository.Repository
	songRequestsSettingsRepository songrequestssettingsrepository.Repository
	spotifyIntegrationsRepository  channelsintegrationsspotify.Repository
	config                         cfg.Config
	logger                         *slog.Logger
	kv                             kv.KV
	wsRouter                       wsrouter.WsRouter
}

func New(
	spotifySongRequestsRepository spotify_song_requests_repository.Repository,
	songRequestsSettingsRepository songrequestssettingsrepository.Repository,
	spotifyIntegrationsRepository channelsintegrationsspotify.Repository,
	config cfg.Config,
	logger *slog.Logger,
	kv kv.KV,
	wsRouter wsrouter.WsRouter,
) *Service {
	return &Service{
		spotifySongRequestsRepository:  spotifySongRequestsRepository,
		songRequestsSettingsRepository: songRequestsSettingsRepository,
		spotifyIntegrationsRepository:  spotifyIntegrationsRepository,
		config:                         config,
		logger:                         logger,
		kv:                             kv,
		wsRouter:                       wsRouter,
	}
}

func SpotifyQueueWsKey(channelID string) string {
	return "api.spotifySongRequestQueue." + channelID
}

func (s *Service) publishQueueChanged(ctx context.Context, channelID string) {
	if s.wsRouter == nil {
		return
	}
	if err := s.wsRouter.Publish(SpotifyQueueWsKey(channelID), struct{}{}); err != nil {
		s.logger.ErrorContext(
			ctx,
			"failed to publish spotify queue change",
			logger.Error(err),
			slog.String("channel_id", channelID),
		)
	}
}

func (s *Service) CreateRequest(
	ctx context.Context,
	channelID string,
	requesterUserID string,
	requesterName string,
	requesterDisplayName string,
	source string,
	query string,
) (spotify_song_request.SpotifySongRequest, error) {
	settings, err := s.loadSpotifySettings(ctx, channelID)
	if err != nil {
		return spotify_song_request.Nil, err
	}

	client, integration, err := s.loadSpotifyClient(ctx, channelID)
	if err != nil {
		return spotify_song_request.Nil, err
	}
	if !slices.Contains(integration.Scopes, "user-modify-playback-state") {
		return spotify_song_request.Nil, spotify.ErrInsufficientScope
	}

	track, err := resolveTrack(ctx, client, query)
	if err != nil {
		return spotify_song_request.Nil, err
	}

	activeRequests, err := s.spotifySongRequestsRepository.GetActiveByChannel(ctx, channelID)
	if err != nil {
		return spotify_song_request.Nil, fmt.Errorf("get active spotify song requests: %w", err)
	}
	for _, active := range activeRequests {
		if active.TrackURI == track.URI {
			return spotify_song_request.Nil, ErrTrackAlreadyInQueue
		}
	}
	if settings.MaxRequests > 0 {
		count, err := s.spotifySongRequestsRepository.CountActiveByChannel(ctx, channelID)
		if err != nil {
			return spotify_song_request.Nil, fmt.Errorf("count active spotify song requests by channel: %w", err)
		}
		if count >= int64(settings.MaxRequests) {
			return spotify_song_request.Nil, ErrMaxRequestsExceeded
		}
	}
	if settings.UserMaxRequests > 0 {
		count, err := s.spotifySongRequestsRepository.CountActiveByRequester(ctx, channelID, requesterName)
		if err != nil {
			return spotify_song_request.Nil, fmt.Errorf("count active spotify song requests by requester: %w", err)
		}
		if count >= int64(settings.UserMaxRequests) {
			return spotify_song_request.Nil, ErrUserMaxRequestsExceeded
		}
	}
	trackDurationMinutes := int(
		math.Round((time.Duration(track.DurationMs) * time.Millisecond).Minutes()),
	)
	if settings.SongMaxLength > 0 && trackDurationMinutes > settings.SongMaxLength {
		return spotify_song_request.Nil, ErrDurationNotAllowed
	}
	if settings.SongMinLength > 0 && trackDurationMinutes < settings.SongMinLength {
		return spotify_song_request.Nil, ErrDurationNotAllowed
	}

	deviceID, err := s.selectDevice(ctx, channelID, client)
	if err != nil {
		return spotify_song_request.Nil, err
	}
	if err := client.AddToQueue(ctx, track.URI, deviceID); err != nil {
		if errors.Is(err, spotify.ErrNoActiveDevice) {
			s.invalidateDeviceCache(ctx, channelID)
			deviceID, selectErr := s.selectDevice(ctx, channelID, client)
			if selectErr != nil {
				return spotify_song_request.Nil, selectErr
			}
			if err = client.AddToQueue(ctx, track.URI, deviceID); err != nil {
				s.logger.ErrorContext(
					ctx,
					"failed to add track to spotify queue",
					logger.Error(err),
					slog.String("channel_id", channelID),
					slog.String("track_uri", track.URI),
					slog.String("device_id", deviceID),
				)
				return spotify_song_request.Nil, err
			}
		} else {
			s.logger.ErrorContext(
				ctx,
				"failed to add track to spotify queue",
				logger.Error(err),
				slog.String("channel_id", channelID),
				slog.String("track_uri", track.URI),
				slog.String("device_id", deviceID),
			)
			return spotify_song_request.Nil, err
		}
	}

	now := time.Now()
	activeCount, err := s.spotifySongRequestsRepository.CountActiveByChannel(ctx, channelID)
	if err != nil {
		return spotify_song_request.Nil, fmt.Errorf("count active spotify song requests by channel: %w", err)
	}

	request := spotify_song_request.SpotifySongRequest{
		ID:            uuid.New(),
		ChannelID:     channelID,
		TrackID:       track.ID,
		TrackURI:      track.URI,
		Title:         track.Name,
		Artist:        track.ArtistName,
		Album:         track.AlbumName,
		DurationMs:    track.DurationMs,
		RequesterName: requesterName,
		Source:        source,
		QueuePosition: int(activeCount) + 1,
		Status:        spotify_song_request.StatusQueued,
		QueuedAt:      now,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	if requesterUserID != "" {
		request.RequesterUserID = &requesterUserID
	}
	if requesterDisplayName != "" {
		request.RequesterDisplayName = &requesterDisplayName
	}

	created, err := s.spotifySongRequestsRepository.Create(ctx, request)
	if err != nil {
		return spotify_song_request.Nil, fmt.Errorf("create spotify song request: %w", err)
	}

	s.publishQueueChanged(ctx, channelID)

	return created, nil
}

func (s *Service) CancelRequest(
	ctx context.Context,
	channelID string,
	requesterName string,
) (spotify_song_request.SpotifySongRequest, error) {
	if _, err := s.loadSpotifySettings(ctx, channelID); err != nil {
		return spotify_song_request.Nil, err
	}

	var (
		requests []spotify_song_request.SpotifySongRequest
		err      error
	)
	if requesterName != "" {
		requests, err = s.spotifySongRequestsRepository.GetActiveByRequester(ctx, channelID, requesterName)
	} else {
		requests, err = s.spotifySongRequestsRepository.GetActiveByChannel(ctx, channelID)
	}
	if err != nil {
		return spotify_song_request.Nil, err
	}
	if len(requests) == 0 {
		return spotify_song_request.Nil, spotify.ErrTrackNotFound
	}

	request := requests[len(requests)-1]
	if err := s.spotifySongRequestsRepository.CancelPendingSkip(ctx, request.ID.String()); err != nil {
		return spotify_song_request.Nil, err
	}
	request.Status = spotify_song_request.StatusCancelledPendingSkip

	s.publishQueueChanged(ctx, channelID)

	return request, nil
}

func (s *Service) SkipRequest(ctx context.Context, channelID, requestID string) error {
	if err := s.ensureRequestChannel(ctx, channelID, requestID); err != nil {
		return err
	}

	if err := s.spotifySongRequestsRepository.UpdateStatus(
		ctx,
		requestID,
		spotify_song_request.StatusSkippedByTwir,
	); err != nil {
		return err
	}

	s.publishQueueChanged(ctx, channelID)

	return nil
}

func (s *Service) CancelRequestByID(ctx context.Context, channelID, requestID string) error {
	if err := s.ensureRequestChannel(ctx, channelID, requestID); err != nil {
		return err
	}

	if err := s.spotifySongRequestsRepository.CancelPendingSkip(ctx, requestID); err != nil {
		return err
	}

	s.publishQueueChanged(ctx, channelID)

	return nil
}

func (s *Service) GetActiveQueue(
	ctx context.Context,
	channelID string,
) ([]spotify_song_request.SpotifySongRequest, error) {
	return s.spotifySongRequestsRepository.GetActiveByChannel(ctx, channelID)
}

func (s *Service) SelectAndCacheDevice(ctx context.Context, channelID string) (string, error) {
	client, _, err := s.loadSpotifyClient(ctx, channelID)
	if err != nil {
		return "", err
	}

	return s.selectDevice(ctx, channelID, client)
}

func (s *Service) SetSelectedDevice(ctx context.Context, channelID, deviceID string) error {
	client, _, err := s.loadSpotifyClient(ctx, channelID)
	if err != nil {
		return err
	}

	devices, err := client.GetDevices(ctx)
	if err != nil {
		return err
	}
	for _, device := range devices {
		if device.ID == deviceID {
			return s.cacheDevice(ctx, channelID, device)
		}
	}

	return spotify.ErrNoActiveDevice
}

func (s *Service) GetDevices(ctx context.Context, channelID string) ([]spotify.Device, error) {
	client, _, err := s.loadSpotifyClient(ctx, channelID)
	if err != nil {
		return nil, err
	}

	return client.GetDevices(ctx)
}

func (s *Service) SearchTracks(
	ctx context.Context,
	channelID string,
	query string,
	limit int,
) ([]spotify.SpotifyTrack, error) {
	client, _, err := s.loadSpotifyClient(ctx, channelID)
	if err != nil {
		return nil, err
	}

	return client.SearchTracks(ctx, query, limit)
}

func (s *Service) ensureRequestChannel(ctx context.Context, channelID, requestID string) error {
	request, err := s.spotifySongRequestsRepository.GetByID(ctx, requestID)
	if err != nil {
		return err
	}
	if request.ChannelID != channelID {
		return spotify_song_requests_repository.ErrNotFound
	}

	return nil
}
