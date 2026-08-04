package song_requests

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/api"
	"github.com/twirapp/twir/libs/entities/requested_song"
	songrequestssettingsentity "github.com/twirapp/twir/libs/entities/song_requests_settings"
	requestedsongsrepository "github.com/twirapp/twir/libs/repositories/requested_songs"
	songrequestssettingsrepository "github.com/twirapp/twir/libs/repositories/song_requests_settings"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	Bus           *buscore.Bus
	PlaybackState *PlaybackStateService

	SettingsRepository songrequestssettingsrepository.Repository
	SongsRepository    requestedsongsrepository.Repository
}

func New(opts Opts) *Service {
	return &Service{
		bus:                opts.Bus,
		playbackState:      opts.PlaybackState,
		settingsRepository: opts.SettingsRepository,
		songsRepository:    opts.SongsRepository,
	}
}

type Service struct {
	bus                *buscore.Bus
	playbackState      *PlaybackStateService
	settingsRepository songrequestssettingsrepository.Repository
	songsRepository    requestedsongsrepository.Repository
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

func mapQueueItem(song requested_song.RequestedSong) QueueItem {
	return QueueItem{
		ID:                   song.VideoID,
		Title:                song.Title,
		SongLink:             song.Link(),
		DurationSeconds:      int(song.Duration),
		OrderedByName:        song.OrderedByName,
		OrderedByDisplayName: song.RequesterDisplayName(),
		QueuePosition:        song.QueuePosition,
		CreatedAt:            song.CreatedAt,
	}
}

func (s *Service) GetQueue(ctx context.Context, channelID string) ([]QueueItem, error) {
	queue, err := s.songsRepository.GetQueue(ctx, channelID)
	if err != nil {
		return nil, fmt.Errorf("failed to get queue: %w", err)
	}

	result := make([]QueueItem, 0, len(queue))
	for _, song := range queue {
		result = append(result, mapQueueItem(song))
	}

	return result, nil
}

func (s *Service) GetSongByVideoID(
	ctx context.Context,
	channelID string,
	videoID string,
) (requested_song.RequestedSong, error) {
	return s.songsRepository.GetByVideoID(ctx, channelID, videoID)
}

func (s *Service) GetCurrentSong(ctx context.Context, channelID string) (*PlaybackState, error) {
	return s.playbackState.GetState(ctx, channelID)
}

func (s *Service) GetSettings(
	ctx context.Context,
	channelID string,
) (songrequestssettingsentity.Settings, error) {
	return s.settingsRepository.GetByChannelID(ctx, channelID)
}

func (s *Service) UpsertSettings(
	ctx context.Context,
	settings songrequestssettingsentity.Settings,
) (songrequestssettingsentity.Settings, songrequestssettingsentity.Settings, error) {
	oldSettings, err := s.settingsRepository.GetByChannelID(ctx, settings.ChannelID.String())
	if err != nil {
		if !errors.Is(err, songrequestssettingsrepository.ErrNotFound) {
			return songrequestssettingsentity.Nil, songrequestssettingsentity.Nil, fmt.Errorf(
				"failed to get current song requests settings: %w",
				err,
			)
		}

		oldSettings = songrequestssettingsentity.Nil
	}

	newSettings, err := s.settingsRepository.Upsert(ctx, settings)
	if err != nil {
		return songrequestssettingsentity.Nil, songrequestssettingsentity.Nil, fmt.Errorf(
			"failed to update song requests settings: %w",
			err,
		)
	}

	return oldSettings, newSettings, nil
}

func (s *Service) GetVolume(ctx context.Context, channelID string) int {
	settings, err := s.settingsRepository.GetByChannelID(ctx, channelID)
	if err != nil {
		return DefaultVolume
	}

	return settings.Volume
}

func (s *Service) SetVolume(ctx context.Context, channelID string, volume int) error {
	return s.settingsRepository.SetVolume(ctx, channelID, volume)
}

func (s *Service) Skip(ctx context.Context, channelID string) error {
	state, err := s.playbackState.GetState(ctx, channelID)
	if err != nil {
		return fmt.Errorf("failed to get playback state: %w", err)
	}

	if state != nil {
		if err := s.songsRepository.SoftDeleteByVideoID(ctx, channelID, state.VideoID); err != nil &&
			!errors.Is(err, requestedsongsrepository.ErrNotFound) {
			return fmt.Errorf("failed to remove current song: %w", err)
		}

		_ = s.bus.Api.SongRequestRemoveFromQueue.Publish(
			ctx,
			api.SongRequestRemoveFromQueue{ChannelID: channelID, VideoID: state.VideoID},
		)
	}

	if err := s.playbackState.ClearState(ctx, channelID); err != nil {
		return err
	}

	s.playbackState.PublishClearedState(channelID)

	return nil
}

func (s *Service) DeleteFromQueue(ctx context.Context, channelID, videoID string) error {
	if err := s.songsRepository.SoftDeleteByVideoID(ctx, channelID, videoID); err != nil {
		if errors.Is(err, requestedsongsrepository.ErrNotFound) {
			return fmt.Errorf("song not found")
		}

		return err
	}

	return s.bus.Api.SongRequestRemoveFromQueue.Publish(
		ctx,
		api.SongRequestRemoveFromQueue{ChannelID: channelID, VideoID: videoID},
	)
}

func (s *Service) ClearQueue(ctx context.Context, channelID string) error {
	if err := s.songsRepository.SoftDeleteAll(ctx, channelID); err != nil {
		return err
	}

	if err := s.playbackState.ClearState(ctx, channelID); err != nil {
		return err
	}

	s.playbackState.PublishClearedState(channelID)

	return nil
}

func (s *Service) ReorderQueue(ctx context.Context, channelID string, videoIDs []string) error {
	if err := s.songsRepository.UpdateQueuePositions(ctx, channelID, videoIDs); err != nil {
		return err
	}

	return s.bus.Api.SongRequestAddToQueue.Publish(
		ctx,
		api.SongRequestAddToQueue{ChannelID: channelID},
	)
}

func (s *Service) GetPublicQueue(ctx context.Context, channelID string) (
	[]entity.SongRequestPublic,
	error,
) {
	queue, err := s.songsRepository.GetQueue(ctx, channelID)
	if err != nil {
		return nil, fmt.Errorf("failed to get queue: %w", err)
	}

	songs := make([]entity.SongRequestPublic, 0, len(queue))

	for _, song := range queue {
		songs = append(
			songs, entity.SongRequestPublic{
				Title:           song.Title,
				UserID:          song.OrderedByID.String(),
				CreatedAt:       song.CreatedAt,
				SongLink:        song.Link(),
				DurationSeconds: int(song.Duration),
			},
		)
	}

	return songs, nil
}
