package stream_stats

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/twirapp/kv"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	"github.com/twirapp/twir/libs/cache/twitch"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	"github.com/twirapp/twir/libs/logger"
	"github.com/twirapp/twir/libs/redis_keys"
	"github.com/twirapp/twir/libs/repositories/channels"
	"github.com/twirapp/twir/libs/repositories/overlays_stream_stats"
	"github.com/twirapp/twir/libs/repositories/overlays_stream_stats/model"
	"github.com/twirapp/twir/libs/repositories/streams"
	"github.com/twirapp/twir/libs/repositories/users"
	channelservice "github.com/twirapp/twir/libs/services/channels"
	"github.com/twirapp/twir/libs/wsrouter"
)

func New(
	repository overlays_stream_stats.Repository,
	wsRouter wsrouter.WsRouter,
	kvClient kv.KV,
	cachedTwitchClient *twitch.CachedTwitchClient,
	streamsRepository streams.Repository,
	channelsRepository channels.Repository,
	usersRepository users.Repository,
	channelService *channelservice.ChannelService,
	loggerClient *slog.Logger,
) *Service {
	return &Service{
		repository:         repository,
		wsRouter:           wsRouter,
		kv:                 kvClient,
		cachedTwitchClient: cachedTwitchClient,
		streamsRepository:  streamsRepository,
		channelsRepository: channelsRepository,
		usersRepository:    usersRepository,
		channelService:     channelService,
		logger:             loggerClient,
	}
}

type Service struct {
	repository         overlays_stream_stats.Repository
	wsRouter           wsrouter.WsRouter
	kv                 kv.KV
	cachedTwitchClient *twitch.CachedTwitchClient
	streamsRepository  streams.Repository
	channelsRepository channels.Repository
	usersRepository    users.Repository
	channelService     *channelservice.ChannelService
	logger             *slog.Logger
}

func (s *Service) GetOrCreate(ctx context.Context, channelID string) (entity.StreamStatsOverlay, error) {
	overlay, err := s.repository.GetByChannelID(ctx, channelID)
	if err != nil {
		if errors.Is(err, overlays_stream_stats.ErrNotFound) {
			created, createErr := s.repository.Create(ctx, createDefaultOverlayInput(channelID))
			if createErr != nil {
				return entity.StreamStatsOverlay{}, createErr
			}

			return mapModelToEntity(created), nil
		}

		return entity.StreamStatsOverlay{}, err
	}

	return mapModelToEntity(overlay), nil
}

type UpdateInput struct {
	ChannelID string
	Settings  entity.StreamStatsOverlay
}

func (s *Service) Update(ctx context.Context, input UpdateInput) (entity.StreamStatsOverlay, error) {
	if _, err := s.GetOrCreate(ctx, input.ChannelID); err != nil {
		return entity.StreamStatsOverlay{}, err
	}

	updated, err := s.repository.Update(ctx, input.ChannelID, overlays_stream_stats.UpdateInput{
		Design:             string(input.Settings.Design),
		ViewersEnabled:     input.Settings.ViewersEnabled,
		ViewersMode:        string(input.Settings.ViewersMode),
		MessagesEnabled:    input.Settings.MessagesEnabled,
		UptimeEnabled:      input.Settings.UptimeEnabled,
		SubscribersEnabled: input.Settings.SubscribersEnabled,
		FollowersEnabled:   input.Settings.FollowersEnabled,
		CustomHTMLEnabled:  input.Settings.CustomHTMLEnabled,
		CustomHTML:         input.Settings.CustomHTML,
		CustomCSS:          input.Settings.CustomCSS,
	})
	if err != nil {
		return entity.StreamStatsOverlay{}, err
	}

	updatedEntity := mapModelToEntity(updated)
	if err := s.wsRouter.Publish(createSettingsSubscriptionKey(input.ChannelID), updatedEntity); err != nil {
		return entity.StreamStatsOverlay{}, err
	}

	return updatedEntity, nil
}

func (s *Service) resolveChannelIDByAPIKey(ctx context.Context, apiKey string) (string, error) {
	if s.channelsRepository != nil {
		channel, err := s.channelsRepository.GetByApiKey(ctx, apiKey)
		if err != nil && !errors.Is(err, channels.ErrNotFound) {
			return "", err
		}
		if !channel.IsNil() {
			return channel.ID.String(), nil
		}
	}

	user, err := s.usersRepository.GetByApiKey(ctx, apiKey)
	if err != nil {
		return "", err
	}
	if user.IsNil() {
		return "", errors.New("user not found for provided api key")
	}

	channel, err := s.channelService.GetChannelByBindingUserID(ctx, user.Platform, user.ID)
	if err != nil {
		return "", err
	}

	return channel.ID.String(), nil
}

func (s *Service) buildCounters(ctx context.Context, channelID string) (entity.StreamStatsOverlayCounters, error) {
	parsedChannelID, err := uuid.Parse(channelID)
	if err != nil {
		return entity.StreamStatsOverlayCounters{}, fmt.Errorf("parse channel ID: %w", err)
	}

	streamList, err := s.streamsRepository.GetListByChannelID(ctx, parsedChannelID)
	if err != nil {
		return entity.StreamStatsOverlayCounters{}, fmt.Errorf("get streams by channel ID: %w", err)
	}

	counters := entity.StreamStatsOverlayCounters{
		Live:            len(streamList) > 0,
		PlatformViewers: make([]entity.StreamStatsOverlayPlatformViewers, 0, len(streamList)),
	}
	for _, stream := range streamList {
		counters.Viewers += stream.ViewerCount
		counters.PlatformViewers = append(counters.PlatformViewers, entity.StreamStatsOverlayPlatformViewers{
			Platform: stream.Platform,
			Viewers:  stream.ViewerCount,
		})

		parsedMessages, messageErr := s.kv.Get(ctx, redis_keys.StreamParsedMessages(stream.ID)).Int()
		if messageErr == nil {
			counters.Messages += int(parsedMessages)
		}

		if counters.StartedAt == nil || stream.StartedAt.Before(*counters.StartedAt) {
			startedAt := stream.StartedAt
			counters.StartedAt = &startedAt
		}
	}

	channel, err := s.channelService.GetChannelByID(ctx, parsedChannelID)
	if err != nil {
		return entity.StreamStatsOverlayCounters{}, fmt.Errorf("get channel by ID: %w", err)
	}
	if channel.IsNil() {
		return counters, nil
	}

	twitchBinding, hasTwitchBinding := channel.Binding(platformentity.PlatformTwitch)
	if !hasTwitchBinding || twitchBinding.PlatformChannelID == "" {
		return counters, nil
	}

	subscribers, err := s.cachedTwitchClient.GetChannelSubscribersCountByChannelId(
		ctx,
		twitchBinding.UserID,
		twitchBinding.PlatformChannelID,
	)
	if err != nil {
		s.logger.Error("failed to get channel subscribers count", logger.Error(err), slog.String("channel_id", channelID))
	} else {
		counters.Subscribers = &subscribers
	}

	followers, err := s.cachedTwitchClient.GetChannelFollowersCountByChannelId(
		ctx,
		twitchBinding.UserID,
		twitchBinding.PlatformChannelID,
	)
	if err != nil {
		s.logger.Error("failed to get channel followers count", logger.Error(err), slog.String("channel_id", channelID))
	} else {
		counters.Followers = &followers
	}

	return counters, nil
}

func mapModelToEntity(m model.StreamStatsOverlay) entity.StreamStatsOverlay {
	return entity.StreamStatsOverlay{
		ID:                 m.ID,
		ChannelID:          m.ChannelID,
		Design:             entity.StreamStatsOverlayDesign(m.Design),
		ViewersEnabled:     m.ViewersEnabled,
		ViewersMode:        entity.StreamStatsOverlayViewersMode(m.ViewersMode),
		MessagesEnabled:    m.MessagesEnabled,
		UptimeEnabled:      m.UptimeEnabled,
		SubscribersEnabled: m.SubscribersEnabled,
		FollowersEnabled:   m.FollowersEnabled,
		CustomHTMLEnabled:  m.CustomHTMLEnabled,
		CustomHTML:         m.CustomHTML,
		CustomCSS:          m.CustomCSS,
		CreatedAt:          m.CreatedAt,
		UpdatedAt:          m.UpdatedAt,
	}
}

func createDefaultOverlayInput(channelID string) overlays_stream_stats.CreateInput {
	return overlays_stream_stats.CreateInput{
		ChannelID:          channelID,
		Design:             string(entity.StreamStatsOverlayDesignBar),
		ViewersEnabled:     true,
		ViewersMode:        string(entity.StreamStatsOverlayViewersModeCumulative),
		MessagesEnabled:    true,
		UptimeEnabled:      true,
		SubscribersEnabled: true,
		FollowersEnabled:   true,
		CustomHTMLEnabled:  false,
		CustomHTML:         "",
		CustomCSS:          "",
	}
}
