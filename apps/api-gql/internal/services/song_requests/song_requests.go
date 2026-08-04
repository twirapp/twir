package song_requests

import (
	"context"
	"fmt"
	"time"

	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/api"
	model "github.com/twirapp/twir/libs/gomodels"
	"gorm.io/gorm"
)

type Opts struct {
	Gorm          *gorm.DB
	Bus           *buscore.Bus
	PlaybackState *PlaybackStateService
}

func New(opts Opts) *Service {
	return &Service{
		gorm:          opts.Gorm,
		bus:           opts.Bus,
		playbackState: opts.PlaybackState,
	}
}

type Service struct {
	gorm          *gorm.DB
	bus           *buscore.Bus
	playbackState *PlaybackStateService
}

type QueueItem struct {
	ID                   string    `json:"id"`
	Title                string    `json:"title"`
	SongLink             string    `json:"songLink"`
	DurationSeconds      int       `json:"durationSeconds"`
	OrderedByName        string    `json:"orderedByName"`
	OrderedByDisplayName string    `json:"orderedByDisplayName"`
	QueuePosition        int       `json:"queuePosition"`
	CreatedAt            time.Time `json:"createdAt"`
}

func (s *Service) GetQueue(ctx context.Context, channelID string) ([]QueueItem, error) {
	var queue []model.RequestedSong
	if err := s.gorm.WithContext(ctx).
		Where(`"channelId" = ? AND "deletedAt" IS NULL`, channelID).
		Order(`"queuePosition" asc`).Find(&queue).Error; err != nil {
		return nil, fmt.Errorf("failed to get queue: %w", err)
	}

	result := make([]QueueItem, 0, len(queue))
	for _, song := range queue {
		link := fmt.Sprintf("https://youtu.be/%s", song.VideoID)
		if song.SongLink.Valid {
			link = song.SongLink.String
		}
		result = append(result, QueueItem{
			ID: song.VideoID, Title: song.Title, SongLink: link,
			DurationSeconds: int(song.Duration), OrderedByName: song.OrderedByName,
			OrderedByDisplayName: song.OrderedByDisplayName.String,
			QueuePosition:        song.QueuePosition, CreatedAt: song.CreatedAt,
		})
	}
	return result, nil
}

func (s *Service) GetCurrentSong(ctx context.Context, channelID string) (*PlaybackState, error) {
	return s.playbackState.GetState(ctx, channelID)
}

func (s *Service) Skip(ctx context.Context, channelID string) error {
	state, err := s.playbackState.GetState(ctx, channelID)
	if err != nil {
		return fmt.Errorf("failed to get playback state: %w", err)
	}
	if state != nil {
		now := time.Now()
		if err := s.gorm.WithContext(ctx).Model(&model.RequestedSong{}).
			Where(`"videoId" = ? AND "channelId" = ? AND "deletedAt" IS NULL`, state.VideoID, channelID).
			Update("deletedAt", now).Error; err != nil {
			return fmt.Errorf("failed to remove current song: %w", err)
		}
		_ = s.bus.Api.SongRequestRemoveFromQueue.Publish(ctx, api.SongRequestRemoveFromQueue{ChannelID: channelID, VideoID: state.VideoID})
	}
	if err := s.playbackState.ClearState(ctx, channelID); err != nil {
		return err
	}
	s.playbackState.PublishClearedState(channelID)
	return nil
}

func (s *Service) DeleteFromQueue(ctx context.Context, channelID, videoID string) error {
	now := time.Now()
	result := s.gorm.WithContext(ctx).Model(&model.RequestedSong{}).
		Where(`"videoId" = ? AND "channelId" = ? AND "deletedAt" IS NULL`, videoID, channelID).
		Update("deletedAt", now)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("song not found")
	}
	return s.bus.Api.SongRequestRemoveFromQueue.Publish(ctx, api.SongRequestRemoveFromQueue{ChannelID: channelID, VideoID: videoID})
}

func (s *Service) ClearQueue(ctx context.Context, channelID string) error {
	now := time.Now()
	if err := s.gorm.WithContext(ctx).Model(&model.RequestedSong{}).
		Where(`"channelId" = ? AND "deletedAt" IS NULL`, channelID).
		Update("deletedAt", now).Error; err != nil {
		return err
	}
	if err := s.playbackState.ClearState(ctx, channelID); err != nil {
		return err
	}
	s.playbackState.PublishClearedState(channelID)
	return nil
}

func (s *Service) ReorderQueue(ctx context.Context, channelID string, videoIDs []string) error {
	if err := s.gorm.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for position, videoID := range videoIDs {
			result := tx.Model(&model.RequestedSong{}).
				Where(`"videoId" = ? AND "channelId" = ? AND "deletedAt" IS NULL`, videoID, channelID).
				Update("queuePosition", position)
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 0 {
				return fmt.Errorf("song %s not found", videoID)
			}
		}
		return nil
	}); err != nil {
		return err
	}
	return s.bus.Api.SongRequestAddToQueue.Publish(ctx, api.SongRequestAddToQueue{ChannelID: channelID})
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
