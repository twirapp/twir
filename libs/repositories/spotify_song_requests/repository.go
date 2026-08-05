package spotify_song_requests

import (
	"context"

	"github.com/twirapp/twir/libs/entities/spotify_song_request"
)

type Repository interface {
	Create(ctx context.Context, req spotify_song_request.SpotifySongRequest) (spotify_song_request.SpotifySongRequest, error)
	GetByID(ctx context.Context, id string) (spotify_song_request.SpotifySongRequest, error)
	GetActiveByChannel(ctx context.Context, channelID string) ([]spotify_song_request.SpotifySongRequest, error)
	GetActiveByRequester(ctx context.Context, channelID, requesterName string) ([]spotify_song_request.SpotifySongRequest, error)
	CountActiveByChannel(ctx context.Context, channelID string) (int64, error)
	CountActiveByRequester(ctx context.Context, channelID, requesterName string) (int64, error)
	ListByChannel(ctx context.Context, channelID string, limit int) ([]spotify_song_request.SpotifySongRequest, error)
	UpdateStatus(ctx context.Context, id string, status spotify_song_request.Status) error
	CancelPendingSkip(ctx context.Context, id string) error
}
