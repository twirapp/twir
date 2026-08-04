package requested_songs

import (
	"context"
	"time"

	"github.com/twirapp/twir/libs/entities/requested_song"
)

type Repository interface {
	GetByVideoID(
		ctx context.Context,
		channelID string,
		videoID string,
	) (requested_song.RequestedSong, error)
	GetQueue(ctx context.Context, channelID string) ([]requested_song.RequestedSong, error)
	CountByChannelID(ctx context.Context, channelID string, createdAfter time.Time) (int64, error)
	SoftDeleteByVideoID(ctx context.Context, channelID string, videoID string) error
	SoftDeleteAll(ctx context.Context, channelID string) error
	UpdateQueuePositions(ctx context.Context, channelID string, videoIDs []string) error
}
