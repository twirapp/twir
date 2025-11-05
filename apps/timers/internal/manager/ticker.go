package manager

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/goccy/go-json"
	"github.com/redis/go-redis/v9"
	"github.com/twirapp/twir/libs/redis_keys"
	streamsmodel "github.com/twirapp/twir/libs/repositories/streams/model"
	timersmodel "github.com/twirapp/twir/libs/repositories/timers/model"
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
			slog.Any("err", err),
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
			slog.Any("err", err),
			slog.String("channelId", t.dbRow.ChannelID),
			slog.String("timerId", id.String()),
		)
		return
	}

	if stream == nil {
		return
	}

	streamParsedMessages, err := c.getStreamChatLines(ctx, stream.ID)
	if err != nil {
		c.logger.Error(
			"[tick] cannot get stream parsed messages",
			slog.Any("err", err),
			slog.String("channelId", t.dbRow.ChannelID),
			slog.String("timerId", id.String()),
		)
		return
	}

	var (
		now               = time.Now()
		shouldSend        bool
		timeInterval      = time.Duration(t.dbRow.TimeInterval) * time.Minute
		messageInterval   = t.dbRow.MessageInterval
		messagesSinceLast = streamParsedMessages - t.lastTriggerMessageNumber
		secondsSinceLast  = now.Sub(t.lastTriggerTimestamp).Seconds() + 1 // https://go.dev/pkg/time/?m=old#hdr-Timer_Resolution
	)

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

	t.lastTriggerMessageNumber = streamParsedMessages
	t.lastTriggerTimestamp = now

	var response timersmodel.Response
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
			slog.Any("err", err),
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
