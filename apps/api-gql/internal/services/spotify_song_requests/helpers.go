package spotify_song_requests

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	kvoptions "github.com/twirapp/kv/options"
	"github.com/twirapp/twir/libs/entities/song_request_mode"
	songrequestssettings "github.com/twirapp/twir/libs/entities/song_requests_settings"
	deprecatedgormmodel "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/integrations/spotify"
	"github.com/twirapp/twir/libs/logger"
	channelsintegrationsspotifymodel "github.com/twirapp/twir/libs/repositories/channels_integrations_spotify/model"
	songrequestssettingsrepository "github.com/twirapp/twir/libs/repositories/song_requests_settings"
)

type spotifyClient interface {
	SearchTracks(ctx context.Context, query string, limit int) ([]spotify.SpotifyTrack, error)
	GetTrackByID(ctx context.Context, trackID string) (*spotify.SpotifyTrack, error)
	GetDevices(ctx context.Context) ([]spotify.Device, error)
	AddToQueue(ctx context.Context, trackURI string, deviceID string) error
	GetCurrentlyPlaying(ctx context.Context) (*spotify.CurrentlyPlaying, error)
	GetQueue(ctx context.Context) ([]spotify.SpotifyTrack, error)
	SkipNext(ctx context.Context, deviceID string) error
}

func (s *Service) loadSpotifySettings(
	ctx context.Context,
	channelID string,
) (songrequestssettings.Settings, error) {
	settings, err := s.songRequestsSettingsRepository.GetByChannelID(ctx, channelID)
	if err != nil {
		if errors.Is(err, songrequestssettingsrepository.ErrNotFound) {
			return songrequestssettings.Nil, ErrNotSpotifyMode
		}
		return songrequestssettings.Nil, fmt.Errorf("load song requests settings: %w", err)
	}
	if settings.IsNil() || !settings.Enabled || settings.Mode != song_request_mode.ModeSpotify {
		return songrequestssettings.Nil, ErrNotSpotifyMode
	}

	return settings, nil
}

func (s *Service) loadSpotifyClient(
	ctx context.Context,
	channelID string,
) (spotifyClient, channelsintegrationsspotifymodel.ChannelIntegrationSpotify, error) {
	integration, err := s.spotifyIntegrationsRepository.GetByChannelID(ctx, channelID)
	if err != nil {
		return nil, channelsintegrationsspotifymodel.Nil, fmt.Errorf("load spotify integration: %w", err)
	}
	if integration.ID == uuid.Nil {
		return nil, channelsintegrationsspotifymodel.Nil, spotify.ErrNotConnected
	}

	client := spotify.New(
		deprecatedgormmodel.Integrations{
			ClientID: sql.NullString{String: s.config.SpotifyClientID, Valid: s.config.SpotifyClientID != ""},
			ClientSecret: sql.NullString{
				String: s.config.SpotifySecret,
				Valid:  s.config.SpotifySecret != "",
			},
		},
		channelsintegrationsspotifymodel.ChannelIntegrationSpotify{
			AccessToken:  integration.AccessToken,
			RefreshToken: integration.RefreshToken,
			Scopes:       integration.Scopes,
		},
		s.spotifyIntegrationsRepository,
	)
	if client == nil {
		return nil, channelsintegrationsspotifymodel.Nil, spotify.ErrNotConnected
	}

	return client, integration, nil
}

func (s *Service) selectDevice(ctx context.Context, channelID string, client spotifyClient) (string, error) {
	devices, err := client.GetDevices(ctx)
	if err != nil {
		return "", err
	}

	for _, device := range devices {
		if !device.IsActive || device.IsRestricted {
			continue
		}

		if s.kv != nil {
			if err := s.kv.Set(
				ctx,
				spotifySongRequestDeviceCachePrefix+channelID,
				device.ID,
				kvoptions.WithExpire(spotifySongRequestDeviceCacheTTL),
			); err != nil {
				return "", fmt.Errorf("cache spotify device: %w", err)
			}
		}

		return device.ID, nil
	}

	return "", spotify.ErrNoActiveDevice
}

func (s *Service) invalidateDeviceCache(ctx context.Context, channelID string) {
	if s.kv == nil {
		return
	}
	if err := s.kv.Delete(ctx, spotifySongRequestDeviceCachePrefix+channelID); err != nil {
		s.logger.ErrorContext(
			ctx,
			"failed to invalidate spotify device cache",
			logger.Error(err),
			slog.String("channel_id", channelID),
		)
	}
}

func resolveTrack(
	ctx context.Context,
	client spotifyClient,
	query string,
) (spotify.SpotifyTrack, error) {
	if trackID, ok := spotify.ParseTrackID(query); ok {
		track, err := client.GetTrackByID(ctx, trackID)
		if err != nil {
			return spotify.SpotifyTrack{}, err
		}
		return *track, nil
	}

	tracks, err := client.SearchTracks(ctx, query, 5)
	if err != nil {
		return spotify.SpotifyTrack{}, err
	}
	if len(tracks) == 0 {
		return spotify.SpotifyTrack{}, ErrTrackNotFound
	}

	return tracks[0], nil
}
