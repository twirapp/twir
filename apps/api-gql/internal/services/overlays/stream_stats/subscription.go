package stream_stats

import (
	"context"
	"fmt"
	"time"

	"github.com/goccy/go-json"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	"github.com/twirapp/twir/libs/logger"
)

func createSettingsSubscriptionKey(channelID string) string {
	return fmt.Sprintf("overlays:stream_stats:settings:%s", channelID)
}

func (s *Service) SettingsSubscriptionSignalerByApiKey(ctx context.Context, apiKey string) (<-chan entity.StreamStatsOverlay, error) {
	channelID, err := s.resolveChannelIDByAPIKey(ctx, apiKey)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve channel by api key: %w", err)
	}

	wsRouterSub, err := s.wsRouter.Subscribe([]string{createSettingsSubscriptionKey(channelID)})
	if err != nil {
		return nil, err
	}

	initialSettings, err := s.GetOrCreate(ctx, channelID)
	if err != nil {
		wsRouterSub.Unsubscribe()
		return nil, fmt.Errorf("failed to get stream stats overlay: %w", err)
	}

	outputChan := make(chan entity.StreamStatsOverlay, 1)
	outputChan <- initialSettings

	go func() {
		defer func() {
			wsRouterSub.Unsubscribe()
			close(outputChan)
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case data, ok := <-wsRouterSub.GetChannel():
				if !ok {
					return
				}

				var settings entity.StreamStatsOverlay
				if err := json.Unmarshal(data, &settings); err != nil {
					s.logger.Error("failed to unmarshal stream stats overlay update", logger.Error(err))
					continue
				}

				select {
				case outputChan <- settings:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	return outputChan, nil
}

func (s *Service) CountersSubscriptionSignalerByApiKey(ctx context.Context, apiKey string) (<-chan entity.StreamStatsOverlayCounters, error) {
	channelID, err := s.resolveChannelIDByAPIKey(ctx, apiKey)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve channel by api key: %w", err)
	}

	initialCounters, err := s.buildCounters(ctx, channelID)
	if err != nil {
		return nil, fmt.Errorf("failed to build stream stats counters: %w", err)
	}

	outputChan := make(chan entity.StreamStatsOverlayCounters, 1)
	outputChan <- initialCounters

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer func() {
			ticker.Stop()
			close(outputChan)
		}()

		lastKey := countersKey(initialCounters)
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				counters, buildErr := s.buildCounters(ctx, channelID)
				if buildErr != nil {
					continue
				}

				key := countersKey(counters)
				if key == lastKey {
					continue
				}
				lastKey = key

				select {
				case outputChan <- counters:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	return outputChan, nil
}

func countersKey(counters entity.StreamStatsOverlayCounters) string {
	startedAt := ""
	if counters.StartedAt != nil {
		startedAt = counters.StartedAt.Format(time.RFC3339Nano)
	}

	subscribers := ""
	if counters.Subscribers != nil {
		subscribers = fmt.Sprint(*counters.Subscribers)
	}

	followers := ""
	if counters.Followers != nil {
		followers = fmt.Sprint(*counters.Followers)
	}

	return fmt.Sprintf(
		"%t:%d:%v:%d:%s:%s:%s",
		counters.Live,
		counters.Viewers,
		counters.PlatformViewers,
		counters.Messages,
		startedAt,
		subscribers,
		followers,
	)
}
