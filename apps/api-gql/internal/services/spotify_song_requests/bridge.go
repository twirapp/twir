package spotify_song_requests

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/twirapp/twir/libs/baseapp/lifecycle"
	buscore "github.com/twirapp/twir/libs/bus-core"
	buscorespotify "github.com/twirapp/twir/libs/bus-core/spotify"
	"github.com/twirapp/twir/libs/entities/spotify_song_request"
	"github.com/twirapp/twir/libs/integrations/spotify"
)

type RequestBridge struct {
	service *Service
	bus     *buscore.Bus
	logger  *slog.Logger
}

func NewRequestBridge(
	lc *lifecycle.Lifecycle,
	service *Service,
	bus *buscore.Bus,
	logger *slog.Logger,
) *RequestBridge {
	b := &RequestBridge{
		service: service,
		bus:     bus,
		logger:  logger,
	}

	lc.Append(lifecycle.Hook{
		OnStart: func(ctx context.Context) error {
			if err := bus.Spotify.Search.SubscribeGroup(
				"api",
				b.handleSearch,
			); err != nil {
				return fmt.Errorf("subscribe to spotify search: %w", err)
			}
			if err := bus.Spotify.CreateSongRequest.SubscribeGroup(
				"api",
				b.handleCreateSongRequest,
			); err != nil {
				return fmt.Errorf("subscribe to spotify create song request: %w", err)
			}
			if err := bus.Spotify.CancelSongRequest.SubscribeGroup(
				"api",
				b.handleCancelSongRequest,
			); err != nil {
				return fmt.Errorf("subscribe to spotify cancel song request: %w", err)
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			bus.Spotify.Search.Unsubscribe()
			bus.Spotify.CreateSongRequest.Unsubscribe()
			bus.Spotify.CancelSongRequest.Unsubscribe()
			return nil
		},
	})

	return b
}

func (b *RequestBridge) handleSearch(
	ctx context.Context,
	req buscorespotify.SearchRequest,
) (buscorespotify.SearchResponse, error) {
	tracks, err := b.service.SearchTracks(ctx, req.ChannelID, req.Query, req.Limit)
	if err != nil {
		return buscorespotify.SearchResponse{}, err
	}

	result := make([]buscorespotify.TrackData, 0, len(tracks))
	for _, track := range tracks {
		result = append(result, buscorespotify.TrackData{
			ID:         track.ID,
			URI:        track.URI,
			Name:       track.Name,
			ArtistName: track.ArtistName,
			AlbumName:  track.AlbumName,
			DurationMs: track.DurationMs,
			ImageURL:   track.ImageURL,
		})
	}

	return buscorespotify.SearchResponse{Tracks: result}, nil
}

func (b *RequestBridge) handleCreateSongRequest(
	ctx context.Context,
	req buscorespotify.CreateSongRequestRequest,
) (buscorespotify.CreateSongRequestResponse, error) {
	request, err := b.service.CreateRequest(
		ctx,
		req.ChannelID,
		req.RequesterUserID,
		req.RequesterName,
		req.RequesterDisplayName,
		req.Source,
		req.Query,
	)
	if err != nil {
		if errors.Is(err, spotify.ErrInsufficientScope) {
			return buscorespotify.CreateSongRequestResponse{}, err
		}
		if errors.Is(err, spotify.ErrNotConnected) {
			return buscorespotify.CreateSongRequestResponse{}, err
		}
		if errors.Is(err, spotify.ErrNoActiveDevice) {
			return buscorespotify.CreateSongRequestResponse{}, err
		}
		return buscorespotify.CreateSongRequestResponse{}, err
	}

	return buscorespotify.CreateSongRequestResponse{Request: mapSongRequest(request)}, nil
}

func (b *RequestBridge) handleCancelSongRequest(
	ctx context.Context,
	req buscorespotify.CancelSongRequestRequest,
) (buscorespotify.CancelSongRequestResponse, error) {
	request, err := b.service.CancelRequest(ctx, req.ChannelID, req.RequesterName)
	if err != nil {
		return buscorespotify.CancelSongRequestResponse{}, err
	}

	return buscorespotify.CancelSongRequestResponse{Request: mapSongRequest(request)}, nil
}

func mapSongRequest(request spotify_song_request.SpotifySongRequest) buscorespotify.SongRequestData {
	data := buscorespotify.SongRequestData{
		ID:            request.ID.String(),
		ChannelID:     request.ChannelID,
		TrackID:       request.TrackID,
		TrackURI:      request.TrackURI,
		Title:         request.Title,
		Artist:        request.Artist,
		Album:         request.Album,
		DurationMs:    request.DurationMs,
		RequesterName: request.RequesterName,
		Source:        request.Source,
		QueuePosition: request.QueuePosition,
		Status:        request.Status.String(),
		CreatedAt:     request.CreatedAt.Format(time.RFC3339),
	}
	if request.RequesterUserID != nil {
		data.RequesterUserID = *request.RequesterUserID
	}
	if request.RequesterDisplayName != nil {
		data.RequesterDisplayName = *request.RequesterDisplayName
	}
	return data
}
