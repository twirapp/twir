package manager

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/goccy/go-json"
	"github.com/redis/go-redis/v9"
	timersentity "github.com/twirapp/twir/libs/entities/timers"
	"github.com/twirapp/twir/libs/logger"
	"github.com/twirapp/twir/libs/redis_keys"
	streamsmodel "github.com/twirapp/twir/libs/repositories/streams/model"
)

func (c *Manager) tryTick(id TimerID) {
	t, ok := c.timers[id]
	if !ok {
		return
	}

	ctx := context.Background()

	channel, err := c.channelCachedRepo.Get(ctx, t.dbRow.ChannelID)
	if err != nil {
		c.logger.Error(
			"[tick] cannot get channel",
			logger.Error(err),
			slog.String("channelId", t.dbRow.ChannelID),
			slog.String("timerId", id.String()),
		)
		return
	}

	if !channel.IsBotMod || !channel.IsEnabled {
		return
	}

	stream, err := c.getChannelStream(ctx, t.dbRow.ChannelID)
	if err != nil {
		c.logger.Error(
			"[tick] cannot get channel stream",
			logger.Error(err),
			slog.String("channelId", t.dbRow.ChannelID),
			slog.String("timerId", id.String()),
		)
		return
	}

	isOffline := stream == nil

	if isOffline {
		if !t.dbRow.OfflineEnabled {
			return
		}
		t.offlineMessageNumber++
	} else {
		if !t.dbRow.OnlineEnabled {
			return
		}
		t.offlineMessageNumber = 0
		t.lastTriggerOfflineNumber = 0
	}

	var (
		currentMessageNumber    int
		lastTriggerMessageCount int
	)

	if isOffline {
		currentMessageNumber = t.offlineMessageNumber
		lastTriggerMessageCount = t.lastTriggerOfflineNumber
	} else {
		streamParsedMessages, err := c.getStreamChatLines(ctx, stream.ID)
		if err != nil {
			c.logger.Error(
				"[tick] cannot get stream parsed messages",
				logger.Error(err),
				slog.String("channelId", t.dbRow.ChannelID),
				slog.String("timerId", id.String()),
			)
			return
		}

		currentMessageNumber = streamParsedMessages
		lastTriggerMessageCount = t.lastTriggerMessageNumber
	}

	var (
		now              = time.Now()
		shouldSend       bool
		timeInterval     = time.Duration(t.dbRow.TimeInterval) * time.Minute
		messageInterval  = t.dbRow.MessageInterval
		secondsSinceLast = now.Sub(t.lastTriggerTimestamp).
					Seconds() +
			1 // https://go.dev/pkg/time/?m=old#hdr-Timer_Resolution
	)

	if currentMessageNumber < lastTriggerMessageCount {
		if isOffline {
			t.lastTriggerOfflineNumber = currentMessageNumber
		} else {
			t.lastTriggerMessageNumber = currentMessageNumber
		}
		lastTriggerMessageCount = currentMessageNumber
	}

	messagesSinceLast := currentMessageNumber - lastTriggerMessageCount

	switch {
	case timeInterval == 0 && messageInterval > 0:
		if messagesSinceLast >= messageInterval {
			shouldSend = true
		}
		break
	case timeInterval > 0 && messageInterval == 0:
		if secondsSinceLast >= timeInterval.Seconds() {
			shouldSend = true
		}
		break
	case timeInterval > 0 && messageInterval > 0:
		var (
			timeTriggered bool
			msgTriggered  bool
		)
		if secondsSinceLast >= timeInterval.Seconds() {
			timeTriggered = true
		}
		if messagesSinceLast >= messageInterval {
			msgTriggered = true
		}
		shouldSend = timeTriggered && msgTriggered
	}

	if !shouldSend {
		return
	}

	if isOffline {
		t.lastTriggerOfflineNumber = currentMessageNumber
	} else {
		t.lastTriggerMessageNumber = currentMessageNumber
	}
	t.lastTriggerTimestamp = now

	var response timersentity.Response
	for index, r := range t.dbRow.Responses {
		if index == t.currentResponseIndex {
			response = r
			break
		}
	}

	err = c.sendMessage(
		ctx,
		channel.ID,
		response.Text,
		response.IsAnnounce,
		response.AnnounceColor,
		response.Count,
	)
	if err != nil {
		c.logger.Error(
			"[tick] cannot send timer message",
			logger.Error(err),
			slog.String("channelId", t.dbRow.ChannelID),
			slog.String("timerId", id.String()),
		)
		return
	}

	nextIndex := t.currentResponseIndex + 1

	if nextIndex >= len(t.dbRow.Responses) {
		nextIndex = 0
	}

	t.currentResponseIndex = nextIndex
}

func (c *Manager) getChannelStream(ctx context.Context, channelID string) (
	*streamsmodel.Stream,
	error,
) {
	cacheKey := redis_keys.StreamByChannelID(channelID)
	cachedBytes, err := c.redis.Get(ctx, cacheKey).Bytes()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, fmt.Errorf("failed to get stream cache: %w", err)
	}

	if len(cachedBytes) > 0 {
		var stream streamsmodel.Stream
		if err := json.Unmarshal(cachedBytes, &stream); err != nil {
			return nil, err
		}

		return &stream, nil
	}

	return nil, nil
}

func (c *Manager) getStreamChatLines(ctx context.Context, streamID string) (int, error) {
	streamParsedMessages, err := c.redis.Get(
		ctx,
		redis_keys.StreamParsedMessages(streamID),
	).Int()
	if err != nil && !errors.Is(err, redis.Nil) {
		return 0, fmt.Errorf("failed to get stream parsed messages: %w", err)
	}

	return streamParsedMessages, nil
}
