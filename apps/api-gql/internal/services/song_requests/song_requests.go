package song_requests

import (
	"context"
	"fmt"

	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	model "github.com/twirapp/twir/libs/gomodels"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In

	Gorm *gorm.DB
}

func New(opts Opts) *Service {
	return &Service{
		gorm: opts.Gorm,
	}
}

type Service struct {
	gorm *gorm.DB
}

func (s *Service) GetPublicQueue(ctx context.Context, channelID string) (
	[]entity.SongRequestPublic,
	error,
) {
	var queue []model.RequestedSong
	if err := s.gorm.
		WithContext(ctx).
		Where(`"channelId" = ? AND "deletedAt" IS NULL`, channelID).
		Order(`"queuePosition" asc`).
		Find(&queue).Error; err != nil {
		return nil, fmt.Errorf("failed to get queue: %w", err)
	}

	songs := make([]entity.SongRequestPublic, 0, len(queue))

	for _, song := range queue {
		songLink := fmt.Sprintf(
			"https://youtu.be/%s",
			song.ID,
		)
		if song.SongLink.Valid {
			songLink = song.SongLink.String
		}

		songs = append(
			songs, entity.SongRequestPublic{
				Title:           song.Title,
				UserID:          song.OrderedById,
				CreatedAt:       song.CreatedAt,
				SongLink:        songLink,
				DurationSeconds: int(song.Duration),
			},
		)
	}

	return songs, nil
}
