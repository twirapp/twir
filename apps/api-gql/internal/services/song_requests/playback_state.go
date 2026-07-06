package song_requests

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	gojson "github.com/goccy/go-json"
	"github.com/redis/go-redis/v9"
	"github.com/twirapp/twir/libs/bus-core/api"
	"github.com/twirapp/twir/libs/wsrouter"
	"go.uber.org/fx"
)

const playbackKeyPrefix = "songrequests:playback:"

type PlaybackStateOpts struct {
	fx.In

	Redis    *redis.Client
	Logger   *slog.Logger
	WsRouter wsrouter.WsRouter
}

type PlaybackStateService struct {
	redis    *redis.Client
	logger   *slog.Logger
	wsRouter wsrouter.WsRouter
}

func NewPlaybackStateService(opts PlaybackStateOpts) *PlaybackStateService {
	return &PlaybackStateService{
		redis:    opts.Redis,
		logger:   opts.Logger,
		wsRouter: opts.WsRouter,
	}
}

type PlaybackState struct {
	VideoID   string  `json:"videoId"`
	Title     string  `json:"title"`
	Position  float64 `json:"position"`
	IsPlaying bool    `json:"isPlaying"`
	Volume    int     `json:"volume"`
	StartedAt int64   `json:"startedAt"`
	UpdatedAt int64   `json:"updatedAt"`
}

func playbackKey(channelID string) string {
	return playbackKeyPrefix + channelID
}

func (s *PlaybackStateService) save(ctx context.Context, channelID string, state PlaybackState) error {
	state.UpdatedAt = time.Now().UnixMilli()

	data, err := gojson.Marshal(state)
	if err != nil {
		return fmt.Errorf("failed to marshal playback state: %w", err)
	}

	return s.redis.Set(ctx, playbackKey(channelID), data, 0).Err()
}

func (s *PlaybackStateService) SetPlaying(
	ctx context.Context,
	channelID string,
	videoID string,
	title string,
) error {
	now := time.Now().UnixMilli()
	state := PlaybackState{
		VideoID:   videoID,
		Title:     title,
		Position:  0,
		IsPlaying: true,
		Volume:    100,
		StartedAt: now,
		UpdatedAt: now,
	}

	return s.save(ctx, channelID, state)
}

func (s *PlaybackStateService) SetPaused(ctx context.Context, channelID string) error {
	state, err := s.GetState(ctx, channelID)
	if err != nil {
		return err
	}

	if state == nil {
		return nil
	}

	state.Position = s.computePosition(state)
	state.IsPlaying = false
	state.StartedAt = 0

	return s.save(ctx, channelID, *state)
}

func (s *PlaybackStateService) SetVolume(
	ctx context.Context,
	channelID string,
	volume int,
) error {
	state, err := s.GetState(ctx, channelID)
	if err != nil {
		return err
	}

	if state == nil {
		return nil
	}

	state.Volume = volume

	return s.save(ctx, channelID, *state)
}

func (s *PlaybackStateService) UpdatePosition(
	ctx context.Context,
	channelID string,
	position float64,
) error {
	state, err := s.GetState(ctx, channelID)
	if err != nil {
		return err
	}

	if state == nil {
		return nil
	}

	state.Position = position
	state.StartedAt = time.Now().UnixMilli() - int64(position*1000)

	return s.save(ctx, channelID, *state)
}

func (s *PlaybackStateService) computePosition(state *PlaybackState) float64 {
	if !state.IsPlaying || state.StartedAt == 0 {
		return state.Position
	}

	now := time.Now().UnixMilli()
	elapsed := float64(now-state.StartedAt) / 1000.0
	return state.Position + elapsed
}

func (s *PlaybackStateService) GetState(
	ctx context.Context,
	channelID string,
) (*PlaybackState, error) {
	data, err := s.redis.Get(ctx, playbackKey(channelID)).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get playback state: %w", err)
	}

	var state PlaybackState
	if err := gojson.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("failed to unmarshal playback state: %w", err)
	}

	state.Position = s.computePosition(&state)

	return &state, nil
}

func (s *PlaybackStateService) ClearState(
	ctx context.Context,
	channelID string,
) error {
	return s.redis.Del(ctx, playbackKey(channelID)).Err()
}

func PlaybackStateWsKey(channelID string) string {
	return "api.songRequestPlayback." + channelID
}

func (s *PlaybackStateService) PublishState(ctx context.Context, channelID string) {
	state, err := s.GetState(ctx, channelID)
	if err != nil || state == nil {
		return
	}

	msg := api.SongRequestPlaybackState{
		ChannelID: channelID,
		VideoID:   state.VideoID,
		Title:     state.Title,
		Position:  state.Position,
		IsPlaying: state.IsPlaying,
		Volume:    state.Volume,
		UpdatedAt: state.UpdatedAt,
	}

	if err := s.wsRouter.Publish(PlaybackStateWsKey(channelID), msg); err != nil {
		s.logger.Error("failed to publish playback state", slog.String("channelID", channelID), slog.Any("error", err))
	}
}

func (s *PlaybackStateService) PublishClearedState(channelID string) {
	msg := api.SongRequestPlaybackState{
		ChannelID: channelID,
		VideoID:   "",
		Title:     "",
		Position:  0,
		IsPlaying: false,
		Volume:    0,
		UpdatedAt: time.Now().UnixMilli(),
	}

	if err := s.wsRouter.Publish(PlaybackStateWsKey(channelID), msg); err != nil {
		s.logger.Error("failed to publish cleared state", slog.String("channelID", channelID), slog.Any("error", err))
	}
}

func (s *PlaybackStateService) publishState(
	channelID string,
	state *PlaybackState,
) {
	msg := api.SongRequestPlaybackState{
		ChannelID: channelID,
		VideoID:   state.VideoID,
		Title:     state.Title,
		Position:  state.Position,
		IsPlaying: state.IsPlaying,
		Volume:    state.Volume,
		UpdatedAt: state.UpdatedAt,
	}

	if err := s.wsRouter.Publish(PlaybackStateWsKey(channelID), msg); err != nil {
		s.logger.Error("failed to publish playback state", slog.String("channelID", channelID), slog.Any("error", err))
	}
}

func (s *PlaybackStateService) StartTicker(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				s.tick(ctx)
			}
		}
	}()
}

func (s *PlaybackStateService) tick(ctx context.Context) {
	var cursor uint64
	for {
		keys, nextCursor, err := s.redis.Scan(ctx, cursor, playbackKeyPrefix+"*", 100).Result()
		if err != nil {
			s.logger.Error("failed to scan playback keys", slog.Any("error", err))
			return
		}

		for _, key := range keys {
			channelID := strings.TrimPrefix(key, playbackKeyPrefix)
			if channelID == key {
				continue
			}

			state, err := s.GetState(ctx, channelID)
			if err != nil {
				s.logger.Error("failed to get playback state for tick", slog.String("channelID", channelID), slog.Any("error", err))
				continue
			}

			if state == nil || !state.IsPlaying {
				continue
			}

			s.publishState(channelID, state)
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}
}
