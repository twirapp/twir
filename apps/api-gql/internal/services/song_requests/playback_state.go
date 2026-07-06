package song_requests

import (
	"context"
	"fmt"
	"log/slog"
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

// PlaybackState stored in Redis.
// Position = base position (frozen at pause, 0 at fresh play).
// StartedAt = unix ms when playback started (0 when paused).
// Current position = Position + (now - StartedAt)/1000 when playing.
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

// computePosition returns the current playback position.
// If playing: base position + elapsed time since StartedAt.
// If paused: base position (frozen).
func computePosition(state *PlaybackState) float64 {
	if !state.IsPlaying || state.StartedAt == 0 {
		return state.Position
	}
	elapsed := float64(time.Now().UnixMilli()-state.StartedAt) / 1000.0
	return state.Position + elapsed
}

// readState reads raw state from Redis.
func (s *PlaybackStateService) readState(ctx context.Context, channelID string) (*PlaybackState, error) {
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

	return &state, nil
}

// GetState returns state with computed position for clients.
func (s *PlaybackStateService) GetState(ctx context.Context, channelID string) (*PlaybackState, error) {
	state, err := s.readState(ctx, channelID)
	if err != nil || state == nil {
		return state, err
	}

	computed := *state
	computed.Position = computePosition(state)
	return &computed, nil
}

// save writes state to Redis.
func (s *PlaybackStateService) save(ctx context.Context, channelID string, state PlaybackState) error {
	state.UpdatedAt = time.Now().UnixMilli()
	data, err := gojson.Marshal(state)
	if err != nil {
		return fmt.Errorf("failed to marshal playback state: %w", err)
	}
	return s.redis.Set(ctx, playbackKey(channelID), data, 0).Err()
}

// SetPlaying starts or resumes playback.
// position=0 for new play, or paused position for resume.
func (s *PlaybackStateService) SetPlaying(
	ctx context.Context,
	channelID string,
	videoID string,
	title string,
	position float64,
) error {
	// Preserve volume from existing state
	volume := 100
	existing, _ := s.readState(ctx, channelID)
	if existing != nil {
		volume = existing.Volume
	}

	now := time.Now().UnixMilli()
	state := PlaybackState{
		VideoID:   videoID,
		Title:     title,
		Position:  position, // base position
		IsPlaying: true,
		Volume:    volume,
		StartedAt: now, // playback started now
		UpdatedAt: now,
	}

	return s.save(ctx, channelID, state)
}

// SetPaused freezes the current position.
func (s *PlaybackStateService) SetPaused(ctx context.Context, channelID string) error {
	state, err := s.readState(ctx, channelID)
	if err != nil || state == nil {
		return err
	}

	state.Position = computePosition(state) // freeze at current position
	state.IsPlaying = false
	state.StartedAt = 0

	return s.save(ctx, channelID, *state)
}

// SetVolume updates volume without affecting position.
func (s *PlaybackStateService) SetVolume(ctx context.Context, channelID string, volume int) error {
	state, err := s.readState(ctx, channelID)
	if err != nil || state == nil {
		return err
	}

	state.Volume = volume
	return s.save(ctx, channelID, *state)
}

// UpdatePosition sets a new position (from seek).
func (s *PlaybackStateService) UpdatePosition(ctx context.Context, channelID string, position float64) error {
	state, err := s.readState(ctx, channelID)
	if err != nil || state == nil {
		return err
	}

	state.Position = position
	if state.IsPlaying {
		state.StartedAt = time.Now().UnixMilli()
	}

	return s.save(ctx, channelID, *state)
}

// ClearState removes the playback state from Redis.
func (s *PlaybackStateService) ClearState(ctx context.Context, channelID string) error {
	return s.redis.Del(ctx, playbackKey(channelID)).Err()
}

// --- Publishing to wsRouter ---

func PlaybackStateWsKey(channelID string) string {
	return "api.songRequestPlayback." + channelID
}

func (s *PlaybackStateService) publishToWsRouter(channelID string, state *PlaybackState) {
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

// PublishState reads current state and publishes to wsRouter (used by mutations).
func (s *PlaybackStateService) PublishState(ctx context.Context, channelID string) {
	state, err := s.GetState(ctx, channelID)
	if err != nil || state == nil {
		return
	}
	s.publishToWsRouter(channelID, state)
}

// PublishClearedState publishes an empty state (used by skip/clear).
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
